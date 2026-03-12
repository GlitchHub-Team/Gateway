package configrepositories

import (
	"fmt"

	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (r *SQLiteConfigRepository) DeleteGateway(cmdData *commanddata.DeleteGateway) error {
	deleteGatewayQuery := `DELETE FROM gateways WHERE id = ?`
	_, err := r.dbConnection.ExecContext(r.ctx, deleteGatewayQuery, cmdData.GatewayId.String())
	if err != nil {
		return fmt.Errorf("fallito a eliminare il gateway %s: %w", cmdData.GatewayId, err)
	}

	return nil
}
