package buffereddatasender

import (
	sensor "Gateway/internal/sensor"
)

type DataSender interface {
	Start() error
	Stop() error
}

type SendSensorDataRepository interface {
	Send(data *sensor.SensorData) error
}
