package configrepositories

import (
	"fmt"

	configmanager "Gateway/internal/configManager"
	"Gateway/internal/sensor"
	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

func (r *SQLiteConfigRepository) GetAllGateways() (map[uuid.UUID]*configmanager.Gateway, error) {
	query := `
		SELECT *
		FROM gateways
	`
	rows, err := r.dbConnection.QueryContext(r.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fallito a fare la query dei gateway: %w", err)
	}

	gateways := make(map[uuid.UUID]*configmanager.Gateway)
	for rows.Next() {
		var gatewayId, tenantId uuid.UUID
		var statusStr string
		if err := rows.Scan(&gatewayId, &tenantId, &statusStr); err != nil {
			return nil, fmt.Errorf("fallito a scansionare una riga gateway: %w", err)
		}

		sensors, err := r.loadSensors(gatewayId)
		if err != nil {
			return nil, fmt.Errorf("fallito a recuperare i sensori del gateway %s: %w", gatewayId, err)
		}

		frequencies, err := r.loadFrequencies(gatewayId)
		if err != nil {
			return nil, fmt.Errorf("fallito a recuperare le frequenze del gateway %s: %w", gatewayId, err)
		}

		gateways[gatewayId] = &configmanager.Gateway{
			Id:                       gatewayId,
			TenantId:                 tenantId,
			Status:                   configmanager.GatewayStatus(statusStr),
			Sensors:                  sensors,
			SensorProfileFrequencies: frequencies,
		}
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("fallito a chiudere le righe nel caricamento dei gateway: %w", err)
	}

	return gateways, rows.Err()
}

func (r *SQLiteConfigRepository) loadSensors(gatewayId uuid.UUID) (map[uuid.UUID]*sensor.Sensor, error) {
	query := `SELECT id, profile, status, frequency FROM sensors WHERE gatewayId = ?`
	rows, err := r.dbConnection.QueryContext(r.ctx, query, gatewayId.String())
	if err != nil {
		return nil, err
	}

	sensors := make(map[uuid.UUID]*sensor.Sensor)
	for rows.Next() {
		var sensorId uuid.UUID
		var profileStr, statusStr string
		var frequency int
		if err := rows.Scan(&sensorId, &profileStr, &statusStr, &frequency); err != nil {
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
			Frequency: sensor.SensorFrequency(frequency),
		}
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("fallito a chiudere le righe nel caricamento dei sensori per il gateway %s: %w", gatewayId, err)
	}

	return sensors, nil
}

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

		frequencies[profile] = configmanager.ProfileSensorFrequency(frequency)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("fallito a chiudere le righe nel caricamento delle frequenze per il gateway %s: %w", gatewayId, err)
	}

	return frequencies, nil
}
