package configmanager

import (
	"time"

	credentialsgenerator "Gateway/internal/credentialsGenerator"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensor "Gateway/internal/sensor"

	"github.com/google/uuid"
)

type Gateway struct {
	Id               uuid.UUID
	TenantId         *uuid.UUID
	Sensors          map[uuid.UUID]*sensor.Sensor
	Status           domain.GatewayStatus
	Interval         time.Duration
	PublicIdentifier string  // Public Key
	SecretKey        string  // Private Key
	Token            *string // JWT when gateway is commissioned
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
	CommissionGateway(cmdData *commanddata.CommissionGateway, status domain.GatewayStatus) error
}

type GatewayCreatorPort interface {
	CreateGateway(cmdData *commanddata.CreateGateway, credentials *credentialsgenerator.Credentials, status domain.GatewayStatus) error
}

type GatewayDecommissionerPort interface {
	DecommissionGateway(cmdData *commanddata.DecommissionGateway, status domain.GatewayStatus) error
}

type GatewayDeleterPort interface {
	DeleteGateway(cmdData *commanddata.DeleteGateway) error
}

type GatewayInterrupterPort interface {
	InterruptGateway(cmdData *commanddata.InterruptGateway, status domain.GatewayStatus) error
}

type GatewayResetterPort interface {
	ResetGateway(cmdData *commanddata.ResetGateway, defaultInterval time.Duration) error
}

type GatewayResumerPort interface {
	ResumeGateway(cmdData *commanddata.ResumeGateway, status domain.GatewayStatus) error
}

type SensorInterrupterPort interface {
	InterruptSensor(cmdData *commanddata.InterruptSensor, status sensor.SensorStatus) error
}

type SensorResumerPort interface {
	ResumeSensor(cmdData *commanddata.ResumeSensor, status sensor.SensorStatus) error
}

type SensorAdderPort interface {
	AddSensor(cmdData *commanddata.AddSensor, status sensor.SensorStatus) error
}

type SensorDeleterPort interface {
	DeleteSensor(cmdData *commanddata.DeleteSensor) error
}
