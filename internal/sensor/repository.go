package sensor

import (
	"database/sql"

	profiles "Gateway/internal/sensor/sensorProfiles"
)

type BufferDbConnection struct {
	*sql.DB
}

type SQLiteSaveSensorDataRepository struct {
	dbConnection BufferDbConnection
}

func NewSQLiteSaveSensorDataRepository(conn BufferDbConnection) *SQLiteSaveSensorDataRepository {
	return &SQLiteSaveSensorDataRepository{
		dbConnection: conn,
	}
}

func (r *SQLiteSaveSensorDataRepository) Save(data *profiles.GeneratedSensorData) error {
	// Logic to save the generated sensor data to SQLite database
	return nil
}
