package configrepositories

import (
	"fmt"
	"time"

	configmanager "Gateway/internal/configManager"
	"Gateway/internal/sensor"
	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

func (r *SQLiteConfigRepository) GetAllGateways() (map[uuid.UUID]*configmanager.Gateway, error) {
	query := `
		SELECT id, tenantId, status, interval, publicIdentifier, secretKey, token
		FROM gateways
	`
	rows, err := r.dbConnection.QueryContext(r.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fallito a fare la query dei gateway: %w", err)
	}

	gateways := make(map[uuid.UUID]*configmanager.Gateway)
	for rows.Next() {
		var gatewayId uuid.UUID
		var tenantId *uuid.UUID
		var statusStr string
		var interval int
		var publicIdentifier, secretKey string
		var token *string
		if err := rows.Scan(&gatewayId, &tenantId, &statusStr, &interval, &publicIdentifier, &secretKey, &token); err != nil {
			return nil, fmt.Errorf("fallito a scansionare una riga gateway: %w", err)
		}

		sensors, err := r.loadSensors(gatewayId)
		if err != nil {
			return nil, fmt.Errorf("fallito a recuperare i sensori del gateway %s: %w", gatewayId, err)
		}

		gateways[gatewayId] = &configmanager.Gateway{
			Id:               gatewayId,
			TenantId:         tenantId,
			Status:           configmanager.GatewayStatus(statusStr),
			Sensors:          sensors,
			Interval:         time.Duration(interval) * time.Millisecond,
			PublicIdentifier: publicIdentifier,
			SecretKey:        secretKey,
			Token:            token,
		}
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("fallito a chiudere le righe nel caricamento dei gateway: %w", err)
	}

	return gateways, nil
}

func (r *SQLiteConfigRepository) loadSensors(gatewayId uuid.UUID) (map[uuid.UUID]*sensor.Sensor, error) {
	query := `SELECT id, profile, status, interval FROM sensors WHERE gatewayId = ?`
	rows, err := r.dbConnection.QueryContext(r.ctx, query, gatewayId.String())
	if err != nil {
		return nil, err
	}

	sensors := make(map[uuid.UUID]*sensor.Sensor)
	for rows.Next() {
		var sensorId uuid.UUID
		var profileStr, statusStr string
		var interval int
		if err := rows.Scan(&sensorId, &profileStr, &statusStr, &interval); err != nil {
			return nil, err
		}

		profile := profiles.ParseSensorProfile(profileStr, profiles.NewRand())
		if profile == nil {
			return nil, fmt.Errorf("fallito a fare il parsing del profilo: %s", profileStr)
		}

		sensors[sensorId] = &sensor.Sensor{
			Id:        sensorId,
			GatewayId: gatewayId,
			Profile:   profile,
			Status:    sensor.SensorStatus(statusStr),
			Interval:  time.Duration(interval) * time.Millisecond,
		}
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("fallito a chiudere le righe nel caricamento dei sensori per il gateway %s: %w", gatewayId, err)
	}

	return sensors, nil
}

//Requisito valutato troppo oneroso da implementare, tuttavia è possibile reimplementare il campo presente nel gateway
/*
func (r *SQLiteConfigRepository) loadFrequencies(gatewayId uuid.UUID) (map[profiles.SensorProfile]configmanager.ProfileSensorFrequency, error) {
	query := `SELECT sensorType, frequency FROM sensor_type_frequencies WHERE gatewayId = ?`
	rows, err := r.dbConnection.QueryContext(r.ctx, query, gatewayId.String())
	if err != nil {
		return nil, err
	}

	frequencies := make(map[profiles.SensorProfile]configmanager.ProfileSensorFrequency)
	for rows.Next() {
		var sensorTypeStr string
		var frequency int
		if err := rows.Scan(&sensorTypeStr, &frequency); err != nil {
			continue
		}

		profile := profiles.ParseSensorProfile(sensorTypeStr, profiles.NewRand())
		if profile == nil {
			return nil, fmt.Errorf("fallito a fare il parsing del profilo: %s", sensorTypeStr)
		}

		frequencies[profile] = configmanager.ProfileSensorFrequency(time.Duration(frequency) * time.Millisecond)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("fallito a chiudere le righe nel caricamento delle frequenze per il gateway %s: %w", gatewayId, err)
	}

	return frequencies, nil
}
*/
