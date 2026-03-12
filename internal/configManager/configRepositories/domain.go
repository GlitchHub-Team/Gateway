package configrepositories

import (
	"context"
	"database/sql"
	"time"

	configmanager "Gateway/internal/configManager"
	credentialsgenerator "Gateway/internal/credentialsGenerator"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"

	"github.com/google/uuid"
)

type ConfigDbConnection struct {
	*sql.DB
}

type SQLiteConfigRepository struct {
	ctx          context.Context
	dbConnection ConfigDbConnection
}

func NewSQLiteConfigRepository(ctx context.Context, conn ConfigDbConnection) *SQLiteConfigRepository {
	return &SQLiteConfigRepository{
		ctx:          ctx,
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

func (r *SQLiteConfigRepository) CreateGateway(cmdData *commanddata.CreateGateway, credentials *credentialsgenerator.Credentials) error {
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

func (r *SQLiteConfigRepository) ResetGateway(cmdData *commanddata.ResetGateway, defaultInterval time.Duration) error {
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
