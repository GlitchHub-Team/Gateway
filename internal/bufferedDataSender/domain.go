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
	DataSenderGreeter
	DataSenderDecommissioner
	DataSenderCommissioner
}

type DataSenderStarter interface {
	Start()
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

type DataSenderGreeter interface {
	Hello() error
}

type DataSenderDecommissioner interface {
	Decommission() error
}

type DataSenderCommissioner interface {
	Commission(tenantId uuid.UUID, commissionedToken string) error
}

type sensorData struct {
	SensorId  uuid.UUID
	GatewayId uuid.UUID
	Timestamp time.Time
	Data      []byte
}

type SendSensorDataPort interface {
	Send(data *sensorData, tenantId uuid.UUID) error
	Hello(gatewayId uuid.UUID, publicIdentifier string) error
}

type BufferedDataPort interface {
	GetOrderedBufferedData(gatewayId uuid.UUID) ([]*sensorData, error)
	CleanBufferedData(data []*sensorData) error
	CleanWholeBuffer(gatewayId uuid.UUID) error
}

type SendSensorDataPortFactory interface {
	Create() SendSensorDataPort
	Reload(token string, seed string) (SendSensorDataPort, error)
}
