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
	Stopped        GatewayStatus = "stopped"
)

type Gateway struct {
	Id       uuid.UUID
	TenantId *uuid.UUID
	Sensors  map[uuid.UUID]*sensor.Sensor
	Status   GatewayStatus
	Interval time.Duration
}

type ConfigPort interface {
	GatewaysFetcherPort
	GatewayCommissionerPort
	GatewayCreatorPort
	GatewayDecommissionerPort
	GatewayDeleterPort
	GatewayInterrupterPort
	GatewayResetterPort
	GatewayResumerPort
	SensorInterrupterPort
	SensorResumerPort
	SensorAdderPort
	SensorDeleterPort
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

type GatewayResetterPort interface {
	ResetGateway(cmdData *commanddata.ResetGateway, defaultInterval time.Duration) error
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
