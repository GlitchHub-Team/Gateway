package buffereddatasender

import (
	"time"

	"github.com/google/uuid"
)

type DataSender interface {
	Start()
	Stop()
}

type sensorData struct {
	SensorId  uuid.UUID
	GatewayId uuid.UUID
	Timestamp time.Time
	Data      []byte
}

type SendSensorDataPort interface {
	Send(data *sensorData) error
}

type BufferedDataPort interface {
	GetOrderedBufferedData(gatewayId uuid.UUID) ([]*sensorData, error)
	CleanBufferedData(data []*sensorData) error
}
