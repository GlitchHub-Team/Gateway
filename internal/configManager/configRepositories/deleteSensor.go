package configrepositories

import (
	"fmt"

	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (r *SQLiteConfigRepository) DeleteSensor(cmdData *commanddata.DeleteSensor) error {
	query := `DELETE FROM sensors WHERE id = ? AND gatewayId = ?`

	_, err := r.dbConnection.ExecContext(
		r.ctx,
		query,
		cmdData.SensorId.String(),
		cmdData.GatewayId.String(),
	)
	if err != nil {
		return fmt.Errorf("fallito a eliminare il sensore %s dal gateway %s: %w", cmdData.SensorId, cmdData.GatewayId, err)
	}

	return nil
}
