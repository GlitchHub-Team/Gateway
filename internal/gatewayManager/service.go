package gatewaymanager

import (
	"sync"

	buffereddatasender "Gateway/internal/bufferedDataSender"
	"Gateway/internal/domain"
	"Gateway/internal/sensor"

	"github.com/google/uuid"
)

type GatewayId = uuid.UUID

type SensorId = uuid.UUID

type GatewayWorker struct {
	Sender      buffereddatasender.DataSender
	ErrChannel  chan error
	CmdChannel  chan domain.BaseCommand
}

type GatewayWorkers struct {
	Workers map[GatewayId]GatewayWorker
	Mu      *sync.RWMutex
}

type SensorWorker struct {
	SimulatedSensor sensor.SimulatedSensor
	ErrChannel      chan error
	CmdChannel      chan domain.BaseCommand
}

type SensorWorkers struct {
	Workers map[GatewayId]map[SensorId]SensorWorker
	Mu      *sync.RWMutex
}

func NewGatewayWorkers() GatewayWorkers {
	return GatewayWorkers{
		Workers: make(map[GatewayId]GatewayWorker),
		Mu:      &sync.RWMutex{},
	}
}

func NewSensorWorkers() SensorWorkers {
	return SensorWorkers{
		Workers: make(map[GatewayId]map[SensorId]SensorWorker),
		Mu:      &sync.RWMutex{},
	}
}
