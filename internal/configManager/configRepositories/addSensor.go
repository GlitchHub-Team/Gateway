package configrepositories

import (
	"fmt"

	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
)

func (r *SQLiteConfigRepository) AddSensor(cmdData *commanddata.AddSensor, status sensor.SensorStatus) error {
	query := `
		INSERT INTO sensors (id, gatewayId, profile, status, interval)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.dbConnection.ExecContext(
		r.ctx,
		query,
		cmdData.SensorId.String(),
		cmdData.GatewayId.String(),
		cmdData.Profile.String(),
		string(status),
		cmdData.Interval.Milliseconds(),
	)
	if err != nil {
		return fmt.Errorf("fallito ad aggiungere il sensore %s al gateway %s: %w", cmdData.SensorId, cmdData.GatewayId, err)
	}

	return nil
}
