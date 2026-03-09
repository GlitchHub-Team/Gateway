package configmanager

import (
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensor "Gateway/internal/sensor"
	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

type GatewayStatus string

type ProfileSensorFrequency int

const (
	Active   GatewayStatus = "active"
	Inactive GatewayStatus = "inactive"
)

type Gateway struct {
	Id                       uuid.UUID
	TenantId                 uuid.UUID
	Sensors                  map[uuid.UUID]*sensor.Sensor
	Status                   GatewayStatus
	SensorProfileFrequencies map[profiles.SensorProfile]ProfileSensorFrequency
}

type ConfigPort interface {
	GetAllGateways() (map[uuid.UUID]*Gateway, error)
	GetAllGatewaysByTenantId(tenantId uuid.UUID) (map[uuid.UUID]Gateway, error)
	GetGatewayById(gatewayId uuid.UUID) (*Gateway, error)
	GetSensorById(gatewayId uuid.UUID, sensorId uuid.UUID) (*sensor.Sensor, error)
	ChangeSensorFrequency(cmdData *commanddata.ChangeSensorFrequency) error
	CommissionGateway(cmdData *commanddata.CommissionGateway) error
	CreateGateway(cmdData *commanddata.CreateGateway) error
	DecommissionGateway(cmdData *commanddata.DecommissionGateway) error
	DeleteGateway(cmdData *commanddata.DeleteGateway) error
	InterruptGateway(cmdData *commanddata.InterruptGateway) error
	RebootGateway(cmdData *commanddata.RebootGateway) error
	ResetGateway(cmdData *commanddata.ResetGateway) error
	ResumeGateway(cmdData *commanddata.ResumeGateway) error
	InterruptSensor(cmdData *commanddata.InterruptSensor) error
	ResumeSensor(cmdData *commanddata.ResumeSensor) error
	AddSensor(cmdData *commanddata.AddSensor) error
	DeleteSensor(cmdData *commanddata.DeleteSensor) error
}

type GatewaysFetcher interface {
	GetAllGateways() (map[uuid.UUID]*Gateway, error)
}

// Interfaces for defining methods in ConfigManagerService
type SensorFrequencySetter interface {
	ChangeSensorFrequency(cmdData *commanddata.ChangeSensorFrequency) error
}

type GatewayCommissioner interface {
	CommissionGateway(cmdData *commanddata.CommissionGateway) error
}

type GatewayCreator interface {
	CreateGateway(cmdData *commanddata.CreateGateway) error
}

type GatewayDecommissioner interface {
	DecommissionGateway(cmdData *commanddata.DecommissionGateway) error
}

type GatewayDeleter interface {
	DeleteGateway(cmdData *commanddata.DeleteGateway) error
}

type GatewayInterrupter interface {
	InterruptGateway(cmdData *commanddata.InterruptGateway) error
}

type GatewayRebooter interface {
	RebootGateway(cmdData *commanddata.RebootGateway) error
}

type GatewayResetter interface {
	ResetGateway(cmdData *commanddata.ResetGateway) error
}

type GatewayResumer interface {
	ResumeGateway(cmdData *commanddata.ResumeGateway) error
}

type SensorInterrupter interface {
	InterruptSensor(cmdData *commanddata.InterruptSensor) error
}

type SensorResumer interface {
	ResumeSensor(cmdData *commanddata.ResumeSensor) error
}

type SensorAdder interface {
	AddSensor(cmdData *commanddata.AddSensor) error
}

type SensorDeleter interface {
	DeleteSensor(cmdData *commanddata.DeleteSensor) error
}
