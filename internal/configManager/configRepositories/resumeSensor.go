package configrepositories

import (
	"fmt"

	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
)

func (r *SQLiteConfigRepository) ResumeSensor(cmdData *commanddata.ResumeSensor, status sensor.SensorStatus) error {
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
		return fmt.Errorf("fallito a riprendere il sensore %s del gateway %s: %w", cmdData.SensorId, cmdData.GatewayId, err)
	}

	return nil
}
