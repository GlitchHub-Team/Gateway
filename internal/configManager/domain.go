package configmanager

import (
	gateway "Gateway/internal/gateway"
	sensor "Gateway/internal/sensor"

	"github.com/google/uuid"
)

type ConfigRepository interface {
	GetAllGatewaysByTenantId(tenantId uuid.UUID) ([]gateway.Gateway, error)
	GetGatewayById(gatewayId uuid.UUID) (*gateway.Gateway, error)
	GetSensorById(gatewayId uuid.UUID, sensorId uuid.UUID) (*sensor.Sensor, error)
	ChangeSensorFrequency(dto gateway.ChangeSensorFrequency) error
	CommissionGateway(dto gateway.CommissionGateway) error
	CreateGateway(dto gateway.CreateGateway) error
	DecommissionGateway(dto gateway.DecommissionGateway) error
	DeleteGateway(dto gateway.DeleteGateway) error
	InterruptGateway(dto gateway.InterruptGateway) error
	RebootGateway(dto gateway.RebootGateway) error
	ResetGateway(dto gateway.ResetGateway) error
	ResumeGateway(dto gateway.ResumeGateway) error
	InterruptSensor(dto gateway.InterruptSensor) error
	ResumeSensor(dto gateway.ResumeSensor) error
	AddSensor(dto gateway.AddSensor) error
	DeleteSensor(dto gateway.DeleteSensor) error
}

// Interfaces for defining methods in ConfigManagerService
type SensorFrequencySetter interface {
	ChangeSensorFrequency(dto gateway.ChangeSensorFrequency) error
}

type GatewayCommissioner interface {
	CommissionGateway(dto gateway.CommissionGateway) error
}

type GatewayCreator interface {
	CreateGateway(dto gateway.CreateGateway) error
}

type GatewayDecommissioner interface {
	DecommissionGateway(dto gateway.DecommissionGateway) error
}

type GatewayDeleter interface {
	DeleteGateway(dto gateway.DeleteGateway) error
}

type GatewayInterrupter interface {
	InterruptGateway(dto gateway.InterruptGateway) error
}

type GatewayRebooter interface {
	RebootGateway(dto gateway.RebootGateway) error
}

type GatewayResetter interface {
	ResetGateway(dto gateway.ResetGateway) error
}

type GatewayResumer interface {
	ResumeGateway(dto gateway.ResumeGateway) error
}

type SensorInterrupter interface {
	InterruptSensor(dto gateway.InterruptSensor) error
}

type SensorResumer interface {
	ResumeSensor(dto gateway.ResumeSensor) error
}

type SensorAdder interface {
	AddSensor(dto gateway.AddSensor) error
}

type SensorDeleter interface {
	DeleteSensor(dto gateway.DeleteSensor) error
}
