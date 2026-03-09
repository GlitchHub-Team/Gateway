package configrepositories

import (
	"database/sql"

	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"

	"github.com/google/uuid"
)

type ConfigDbConnection struct {
	*sql.DB
}

type SQLiteConfigRepository struct {
	dbConnection ConfigDbConnection
}

func NewSQLiteConfigRepository(conn ConfigDbConnection) *SQLiteConfigRepository {
	return &SQLiteConfigRepository{
		dbConnection: conn,
	}
}

func (r *SQLiteConfigRepository) GetAllGatewaysByTenantId(tenantId uuid.UUID) (map[uuid.UUID]configmanager.Gateway, error) {
	// TODO: Implement database query
	return nil, nil
}

func (r *SQLiteConfigRepository) GetGatewayById(gatewayId uuid.UUID) (*configmanager.Gateway, error) {
	// TODO: Implement database query
	return nil, nil
}

func (r *SQLiteConfigRepository) GetSensorById(gatewayId uuid.UUID, sensorId uuid.UUID) (*sensor.Sensor, error) {
	// TODO: Implement database query
	return nil, nil
}

func (r *SQLiteConfigRepository) ChangeSensorFrequency(cmdData *commanddata.ChangeSensorFrequency) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) CommissionGateway(cmdData *commanddata.CommissionGateway) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) CreateGateway(cmdData *commanddata.CreateGateway) error {
	// TODO: Implement database insert
	return nil
}

func (r *SQLiteConfigRepository) DecommissionGateway(cmdData *commanddata.DecommissionGateway) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) DeleteGateway(cmdData *commanddata.DeleteGateway) error {
	// TODO: Implement database delete
	return nil
}

func (r *SQLiteConfigRepository) InterruptGateway(cmdData *commanddata.InterruptGateway) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) RebootGateway(cmdData *commanddata.RebootGateway) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) ResetGateway(cmdData *commanddata.ResetGateway) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) ResumeGateway(cmdData *commanddata.ResumeGateway) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) InterruptSensor(cmdData *commanddata.InterruptSensor) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) ResumeSensor(cmdData *commanddata.ResumeSensor) error {
	// TODO: Implement database update
	return nil
}

func (r *SQLiteConfigRepository) AddSensor(cmdData *commanddata.AddSensor) error {
	// TODO: Implement database insert
	return nil
}

func (r *SQLiteConfigRepository) DeleteSensor(cmdData *commanddata.DeleteSensor) error {
	// TODO: Implement database delete
	return nil
}
