package sensorprofiles

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type SensorProfile interface {
	Generate() *GeneratedSensorData
	String() string
}

type SerializableData interface {
	Serialize() ([]byte, error)
}

type GeneratedSensorData struct {
	SensorId  uuid.UUID
	Timestamp time.Time
	Profile   string
	Data      SerializableData
}

type Rand interface {
	Intn(n int) int
	Float64() float64
}

func NewRand() Rand {
	source := rand.NewSource(time.Now().UnixNano())
	return rand.New(source)
}

func ParseSensorProfile(sensorId uuid.UUID, profileType string, rand Rand) SensorProfile {
	switch profileType {
	case "ecg_custom":
		return NewEcgProfile(sensorId, rand)
	case "environmental_sensing":
		return NewEnvironmentalSensingProfile(sensorId, rand)
	case "health_thermometer":
		return NewHealthThermometerProfile(sensorId, rand)
	case "heart_rate":
		return NewHeartRateProfile(sensorId, rand)
	case "pulse_oximeter":
		return NewPulseOximeterProfile(sensorId, rand)
	default:
		return nil
	}
}
