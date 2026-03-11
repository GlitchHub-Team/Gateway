package gatewaymanager

import (
	"sync"

	buffereddatasender "Gateway/internal/bufferedDataSender"
	"Gateway/internal/sensor"

	"github.com/google/uuid"
)

type GatewayId = uuid.UUID

type SensorId = uuid.UUID

type GatewayWorkers struct {
	Workers map[GatewayId]buffereddatasender.DataSender
	Mu      *sync.RWMutex
}

type SensorWorkers struct {
	Workers map[GatewayId]map[SensorId]sensor.SimulatedSensor
	Mu      *sync.RWMutex
}

func NewGatewayWorkers() GatewayWorkers {
	return GatewayWorkers{
		Workers: make(map[GatewayId]buffereddatasender.DataSender),
		Mu:      &sync.RWMutex{},
	}
}

func NewSensorWorkers() SensorWorkers {
	return SensorWorkers{
		Workers: make(map[GatewayId]map[SensorId]sensor.SimulatedSensor),
		Mu:      &sync.RWMutex{},
	}
}
