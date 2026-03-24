package configrepositories

import (
	"fmt"
	"strings"
	"time"

	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"
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
	gatewayIDs := make([]uuid.UUID, 0)
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

		gateways[gatewayId] = &configmanager.Gateway{
			Id:               gatewayId,
			TenantId:         tenantId,
			Status:           domain.GatewayStatus(statusStr),
			Sensors:          make(map[uuid.UUID]*sensor.Sensor),
			Interval:         time.Duration(interval) * time.Millisecond,
			PublicIdentifier: publicIdentifier,
			SecretKey:        secretKey,
			Token:            token,
		}
		gatewayIDs = append(gatewayIDs, gatewayId)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("fallito a chiudere le righe nel caricamento dei gateway: %w", err)
	}

	sensorsByGateway, err := r.loadSensors(gatewayIDs)
	if err != nil {
		return nil, fmt.Errorf("fallito a recuperare i sensori dei gateway: %w", err)
	}

	for id, gateway := range gateways {
		if sensors, exists := sensorsByGateway[id]; exists {
			gateway.Sensors = sensors
		}
	}

	return gateways, nil
}

func (r *SQLiteConfigRepository) loadSensors(gatewayIDs []uuid.UUID) (map[uuid.UUID]map[uuid.UUID]*sensor.Sensor, error) {
	sensorsByGateway := make(map[uuid.UUID]map[uuid.UUID]*sensor.Sensor, len(gatewayIDs))
	if len(gatewayIDs) == 0 {
		return sensorsByGateway, nil
	}

	args := make([]any, 0, len(gatewayIDs))
	placeholders := make([]string, 0, len(gatewayIDs))
	for _, gatewayID := range gatewayIDs {
		args = append(args, gatewayID.String())
		placeholders = append(placeholders, "?")
		sensorsByGateway[gatewayID] = make(map[uuid.UUID]*sensor.Sensor)
	}

	query := fmt.Sprintf(
		`SELECT id, gatewayId, profile, status, interval FROM sensors WHERE gatewayId IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := r.dbConnection.QueryContext(r.ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var sensorId uuid.UUID
		var gatewayId uuid.UUID
		var profileStr, statusStr string
		var interval int
		if err := rows.Scan(&sensorId, &gatewayId, &profileStr, &statusStr, &interval); err != nil {
			return nil, err
		}

		profile := profiles.ParseSensorProfile(profileStr, profiles.NewRand())
		if profile == nil {
			return nil, fmt.Errorf("fallito a fare il parsing del profilo: %s", profileStr)
		}

		sensorsByGateway[gatewayId][sensorId] = &sensor.Sensor{
			Id:        sensorId,
			GatewayId: gatewayId,
			Profile:   profile,
			Status:    sensor.SensorStatus(statusStr),
			Interval:  time.Duration(interval) * time.Millisecond,
		}
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("fallito a chiudere le righe nel caricamento dei sensori: %w", err)
	}

	return sensorsByGateway, nil
}
