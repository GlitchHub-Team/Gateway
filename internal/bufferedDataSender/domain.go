package buffereddatasender

import (
	"time"

	"github.com/google/uuid"
)

type DataSender interface {
	DataSenderStarter
	DataSenderStopper
	DataSenderInterrupter
	DataSenderResumer
	DataSenderResetter
}

type DataSenderStarter interface {
	Start()
	Hello() error
}

type DataSenderStopper interface {
	Stop()
}

type DataSenderInterrupter interface {
	Interrupt()
}

type DataSenderResumer interface {
	Resume()
}

type DataSenderResetter interface {
	Reset(defaultInterval time.Duration) error
}

type sensorData struct {
	SensorId  uuid.UUID
	GatewayId uuid.UUID
	Timestamp time.Time
	Data      []byte
}

type SendSensorDataPort interface {
	Send(data *sensorData) error
	Hello(gatewayId uuid.UUID) error
}

type BufferedDataPort interface {
	GetOrderedBufferedData(gatewayId uuid.UUID) ([]*sensorData, error)
	CleanBufferedData(data []*sensorData) error
	CleanWholeBuffer(gatewayId uuid.UUID) error
}
