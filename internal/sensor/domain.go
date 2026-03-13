package sensor

import (
	"time"

	"github.com/google/uuid"

	profiles "Gateway/internal/sensor/sensorProfiles"
)

type SimulatedSensor interface {
	SensorStarter
	SensorStopper
	SensorInterrupter
	SensorResumer
}

type SensorStarter interface {
	Start()
}

type SensorStopper interface {
	Stop()
}

type SensorInterrupter interface {
	Interrupt()
}

type SensorResumer interface {
	Resume()
}

type SaveSensorDataPort interface {
	Save(data *profiles.GeneratedSensorData, gatewayId uuid.UUID) error
}

type SensorStatus string

const (
	Active   SensorStatus = "active"
	Inactive SensorStatus = "inactive"
	Stopped  SensorStatus = "stopped"
)

type Sensor struct {
	Id        uuid.UUID
	GatewayId uuid.UUID
	Profile   profiles.SensorProfile
	Interval  time.Duration
	Status    SensorStatus
}
