package gateway

import (
	sensor "Gateway/internal/sensor"
	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

type Response struct {
	Status  bool
	Message string
}

type GatewayStatus string

const (
	Active   GatewayStatus = "active"
	Inactive GatewayStatus = "inactive"
)

type Gateway struct {
	Id                     uuid.UUID
	TenantId               uuid.UUID
	Sensors                map[uuid.UUID]sensor.Sensor
	Status                 GatewayStatus
	SensorProfileFrequency map[profiles.SensorProfile]int
}

type CommandResponseRepository interface {
	Reply(r Response) error
}

// DTOs for the methods in ConfigManagerService and in GatewayManagerService
type ChangeSensorFrequency struct {
	GatewayId string
	Profile   profiles.SensorProfile
	Frequency int
}

type CommissionGateway struct {
	GatewayId string
	TenantId  string
}

type CreateGateway struct {
	GatewayId string
}

type DecommissionGateway struct {
	GatewayId string
}

type DeleteGateway struct {
	GatewayId string
}

type InterruptGateway struct {
	GatewayId string
}

type RebootGateway struct {
	GatewayId string
}

type ResetGateway struct {
	GatewayId string
}

type ResumeGateway struct {
	GatewayId string
}

type InterruptSensor struct {
	GatewayId string
	SensorId  string
}

type ResumeSensor struct {
	GatewayId string
	SensorId  string
}

type AddSensor struct {
	GatewayId string
	SensorId  string
	Profile   profiles.SensorProfile
	Frequency int
}

type DeleteSensor struct {
	GatewayId string
	SensorId  string
}
