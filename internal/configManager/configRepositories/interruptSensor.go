package configrepositories

import (
	"fmt"

	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
)

func (r *SQLiteConfigRepository) InterruptSensor(cmdData *commanddata.InterruptSensor, status sensor.SensorStatus) error {
	query := `
		UPDATE sensors 
		SET status = ?
		WHERE id = ? AND gatewayId = ?
	`

	_, err := r.dbConnection.ExecContext(
		r.ctx,
		query,
		string(status),
		cmdData.SensorId.String(),
		cmdData.GatewayId.String(),
	)
	if err != nil {
		return fmt.Errorf("fallito a interrompere il sensore %s del gateway %s: %w", cmdData.SensorId, cmdData.GatewayId, err)
	}

	return nil
}
