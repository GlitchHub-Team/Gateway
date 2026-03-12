package configrepositories

import (
	"fmt"

	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (r *SQLiteConfigRepository) CommissionGateway(cmdData *commanddata.CommissionGateway, status domain.GatewayStatus) error {
	query := `
		UPDATE gateways 
		SET tenantId = ?, token = ?, status = ?
		WHERE id = ?
	`

	_, err := r.dbConnection.ExecContext(
		r.ctx,
		query,
		cmdData.TenantId.String(),
		cmdData.CommissionedToken,
		string(status),
		cmdData.GatewayId.String(),
	)
	if err != nil {
		return fmt.Errorf("fallito a commissionare il gateway %s: %w", cmdData.GatewayId, err)
	}

	return nil
}
