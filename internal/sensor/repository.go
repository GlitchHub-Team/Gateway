package sensor

import (
	"context"
	"database/sql"
	"fmt"

	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

type BufferDbConnection struct {
	*sql.DB
}

type SQLiteSaveSensorDataRepository struct {
	ctx          context.Context
	dbConnection BufferDbConnection
}

func NewSQLiteSaveSensorDataRepository(ctx context.Context, conn BufferDbConnection) *SQLiteSaveSensorDataRepository {
	return &SQLiteSaveSensorDataRepository{
		ctx:          ctx,
		dbConnection: conn,
	}
}

func (r *SQLiteSaveSensorDataRepository) Save(data *profiles.GeneratedSensorData, gatewayId uuid.UUID) error {
	query := `INSERT INTO buffer (gatewayId, sensorId, timestamp, value) VALUES (?, ?, ?, jsonb(?))`
	serializedData, err := data.Data.Serialize()
	if err != nil {
		return fmt.Errorf("errore nella serializzazione dei dati: %w, gatewayId: %s, sensorId: %s", err, gatewayId, data.SensorId)
	}
	_, err = r.dbConnection.ExecContext(r.ctx, query, gatewayId, data.SensorId, data.Timestamp, serializedData)
	if err != nil {
		return fmt.Errorf("errore nel salvataggio del dato: %w, gatewayId: %s, sensorId: %s", err, gatewayId, data.SensorId)
	}
	return nil
}
