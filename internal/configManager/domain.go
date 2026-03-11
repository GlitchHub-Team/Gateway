package configmanager

import (
	"time"

	commanddata "Gateway/internal/gatewayManager/commandData"
	sensor "Gateway/internal/sensor"

	"github.com/google/uuid"
)

type GatewayStatus string

const (
	Active         GatewayStatus = "active"
	Inactive       GatewayStatus = "inactive"
	Decommissioned GatewayStatus = "decommissioned"
)

type Gateway struct {
	Id       uuid.UUID
	TenantId *uuid.UUID
	Sensors  map[uuid.UUID]*sensor.Sensor
	Status   GatewayStatus
	Interval time.Duration
}

type ConfigPort interface {
	GetAllGateways() (map[uuid.UUID]*Gateway, error)
	GetAllGatewaysByTenantId(tenantId uuid.UUID) (map[uuid.UUID]Gateway, error)
	GetGatewayById(gatewayId uuid.UUID) (*Gateway, error)
	GetSensorById(gatewayId uuid.UUID, sensorId uuid.UUID) (*sensor.Sensor, error)
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

type GatewaysFetcherPort interface {
	GetAllGateways() (map[uuid.UUID]*Gateway, error)
}

type GatewayCommissionerPort interface {
	CommissionGateway(cmdData *commanddata.CommissionGateway) error
}

type GatewayCreatorPort interface {
	CreateGateway(cmdData *commanddata.CreateGateway) error
}

type GatewayDecommissionerPort interface {
	DecommissionGateway(cmdData *commanddata.DecommissionGateway) error
}

type GatewayDeleterPort interface {
	DeleteGateway(cmdData *commanddata.DeleteGateway) error
}

type GatewayInterrupterPort interface {
	InterruptGateway(cmdData *commanddata.InterruptGateway) error
}

type GatewayRebooterPort interface {
	RebootGateway(cmdData *commanddata.RebootGateway) error
}

type GatewayResetterPort interface {
	ResetGateway(cmdData *commanddata.ResetGateway) error
}

type GatewayResumerPort interface {
	ResumeGateway(cmdData *commanddata.ResumeGateway) error
}

type SensorInterrupterPort interface {
	InterruptSensor(cmdData *commanddata.InterruptSensor) error
}

type SensorResumerPort interface {
	ResumeSensor(cmdData *commanddata.ResumeSensor) error
}

type SensorAdderPort interface {
	AddSensor(cmdData *commanddata.AddSensor) error
}

type SensorDeleterPort interface {
	DeleteSensor(cmdData *commanddata.DeleteSensor) error
}
