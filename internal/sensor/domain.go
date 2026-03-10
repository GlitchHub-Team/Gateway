package sensor

import (
	"github.com/google/uuid"

	profiles "Gateway/internal/sensor/sensorProfiles"
)

type SimulatedSensor interface {
	Start()
	Stop()
}

type SaveSensorDataPort interface {
	Save(data *profiles.GeneratedSensorData, gatewayId uuid.UUID) error
}

type SensorStatus string

const (
	Active   SensorStatus = "active"
	Inactive SensorStatus = "inactive"
)

type SensorFrequency int

type Sensor struct {
	Id        uuid.UUID
	GatewayId uuid.UUID
	Profile   profiles.SensorProfile
	Frequency SensorFrequency
	Status    SensorStatus
}
