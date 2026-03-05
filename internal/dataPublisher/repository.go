package datapublisher

import (
	sensor "Gateway/internal/sensor"
	"github.com/nats-io/nats.go"
)

type NATSDataPublisherRepository struct {
	nc *nats.Conn
}

func NewNATSDataPublisherRepository(nc *nats.Conn) *NATSDataPublisherRepository {
	return &NATSDataPublisherRepository{
		nc: nc,
	}
}

func (r *NATSDataPublisherRepository) Send(data *sensor.SensorData) error {
	return nil
}
