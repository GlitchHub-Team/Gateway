package sensor

import (
	"time"

	"github.com/google/uuid"

	profiles "Gateway/internal/sensor/sensorProfiles"
)

type SimulatedSensor interface {
	Start() error
	Stop() error
}

type GeneratedSensorData struct {
	SensorId  string
	Timestamp time.Time
	Data      []byte
}

type SaveSensorDataRepository interface {
	Save(data GeneratedSensorData) error
}

type SensorStatus string

const (
	Active   SensorStatus = "active"
	Inactive SensorStatus = "inactive"
)

type SensorData struct {
	SensorId  uuid.UUID
	GatewayId uuid.UUID
	TenantId  uuid.UUID
	Timestamp time.Time
	Data      []byte
}

type Sensor struct {
	Id        uuid.UUID
	GatewayId uuid.UUID
	Profile   profiles.SensorProfile
	Status    SensorStatus
}
