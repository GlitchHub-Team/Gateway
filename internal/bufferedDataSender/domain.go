package buffereddatasender

import (
	"time"

	"github.com/google/uuid"
)

type DataSender interface {
	Start()
	Stop()
	Interrupt()
	Resume()
	Reset() error
	Hello() error
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
