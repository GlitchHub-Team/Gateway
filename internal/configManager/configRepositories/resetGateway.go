package configrepositories

import (
	"fmt"
	"time"

	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (r *SQLiteConfigRepository) ResetGateway(cmdData *commanddata.ResetGateway, defaultInterval time.Duration) error {
	query := `
		UPDATE gateways 
		SET interval = ?
		WHERE id = ?
	`

	defaultIntervalMs := int64(defaultInterval.Milliseconds())

	_, err := r.dbConnection.ExecContext(
		r.ctx,
		query,
		defaultIntervalMs,
		cmdData.GatewayId.String(),
	)
	if err != nil {
		return fmt.Errorf("fallito a resettare il gateway %s: %w", cmdData.GatewayId, err)
	}

	return nil
}
