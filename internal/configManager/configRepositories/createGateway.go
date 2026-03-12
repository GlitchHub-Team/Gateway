package configrepositories

import (
	"fmt"

	credentialsgenerator "Gateway/internal/credentialsGenerator"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (r *SQLiteConfigRepository) CreateGateway(cmdData *commanddata.CreateGateway, credentials *credentialsgenerator.Credentials, status domain.GatewayStatus) error {
	query := `
		INSERT INTO gateways (id, tenantId, status, interval, publicIdentifier, secretKey, token)
		VALUES (?, NULL, ?, ?, ?, ?, NULL)
	`

	intervalMs := int64(cmdData.Interval.Milliseconds())

	_, err := r.dbConnection.ExecContext(
		r.ctx,
		query,
		cmdData.GatewayId.String(),
		string(status),
		intervalMs,
		credentials.PublicIdentifier,
		credentials.SecretKey,
	)
	if err != nil {
		return fmt.Errorf("fallito a creare il gateway %s: %w", cmdData.GatewayId, err)
	}

	return nil
}
