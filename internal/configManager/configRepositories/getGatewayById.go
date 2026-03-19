package configrepositories

import (
	"fmt"
	"time"

	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"

	"github.com/google/uuid"
)

func (r *SQLiteConfigRepository) GetGatewayById(gatewayId uuid.UUID) (*configmanager.Gateway, error) {
	query := `
		SELECT id, tenantId, status, interval, publicIdentifier, secretKey, token
		FROM gateways 
		WHERE id = ?
	`

	row := r.dbConnection.QueryRowContext(r.ctx, query, gatewayId.String())

	var id uuid.UUID
	var tenantId *uuid.UUID
	var statusStr string
	var interval int
	var publicIdentifier, secretKey string
	var token *string

	if err := row.Scan(&id, &tenantId, &statusStr, &interval, &publicIdentifier, &secretKey, &token); err != nil {
		return nil, fmt.Errorf("fallito a recuperare il gateway %s: %w", gatewayId, err)
	}

	sensors, err := r.loadSensors(gatewayId)
	if err != nil {
		return nil, fmt.Errorf("fallito a recuperare i sensori del gateway %s: %w", gatewayId, err)
	}

	gateway := &configmanager.Gateway{
		Id:               id,
		TenantId:         tenantId,
		Status:           domain.GatewayStatus(statusStr),
		Sensors:          sensors,
		Interval:         time.Duration(interval) * time.Millisecond,
		PublicIdentifier: publicIdentifier,
		SecretKey:        secretKey,
		Token:            token,
	}

	return gateway, nil
}
