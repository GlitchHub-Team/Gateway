package sensor

import (
	profiles "Gateway/internal/sensor/sensorProfiles"
	"github.com/google/uuid"
)

type SensorService struct {
	sensorId  uuid.UUID
	profile   profiles.SensorProfile
	frequency int
	status    SensorStatus
}

func NewSensorService(id uuid.UUID, profile profiles.SensorProfile, frequency int, status SensorStatus) *SensorService {
	return &SensorService{
		sensorId:  id,
		profile:   profile,
		frequency: frequency,
		status:    status,
	}
}

func (s *SensorService) Start() {
	// Logic to start the sensor data generation based on the profile and frequency
}

func (s *SensorService) Stop() {
	// Logic to stop the sensor data generation
}
