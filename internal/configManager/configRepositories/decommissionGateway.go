package configrepositories

import (
	"fmt"

	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (r *SQLiteConfigRepository) DecommissionGateway(cmdData *commanddata.DecommissionGateway, status domain.GatewayStatus) error {
	query := `
		UPDATE gateways 
		SET tenantId = NULL, token = NULL, status = ?
		WHERE id = ?
	`

	_, err := r.dbConnection.ExecContext(
		r.ctx,
		query,
		string(status),
		cmdData.GatewayId.String(),
	)
	if err != nil {
		return fmt.Errorf("fallito a decommissionare il gateway %s: %w", cmdData.GatewayId, err)
	}

	return nil
}
