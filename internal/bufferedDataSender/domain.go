package buffereddatasender

import (
	sensor "Gateway/internal/sensor"
)

type DataSender interface {
	Start()
	Stop()
}

type SendSensorDataPort interface {
	Send(data *sensor.SensorData) error
}

type BufferedDataPort interface {
	GetOrderedBufferedData() ([]*sensor.SensorData, error)
}
