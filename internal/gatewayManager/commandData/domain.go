package commanddata

import (
	"time"

	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

type SensorFrequency int

type ChangeSensorFrequency struct {
	GatewayId uuid.UUID
	Profile   profiles.SensorProfile
	Frequency SensorFrequency
}

type CommissionGateway struct {
	GatewayId         uuid.UUID
	TenantId          uuid.UUID
	CommissionedToken string
}

type CreateGateway struct {
	GatewayId uuid.UUID
	Interval  time.Duration
}

type DecommissionGateway struct {
	GatewayId uuid.UUID
}

type DeleteGateway struct {
	GatewayId uuid.UUID
}

type InterruptGateway struct {
	GatewayId uuid.UUID
}

type RebootGateway struct {
	GatewayId uuid.UUID
}

type ResetGateway struct {
	GatewayId uuid.UUID
}

type ResumeGateway struct {
	GatewayId uuid.UUID
}

type InterruptSensor struct {
	GatewayId uuid.UUID
	SensorId  uuid.UUID
}

type ResumeSensor struct {
	GatewayId uuid.UUID
	SensorId  uuid.UUID
}

type AddSensor struct {
	GatewayId uuid.UUID
	SensorId  uuid.UUID
	Profile   profiles.SensorProfile
	Interval  time.Duration
}

type DeleteSensor struct {
	GatewayId uuid.UUID
	SensorId  uuid.UUID
}
