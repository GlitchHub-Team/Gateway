package buffereddatasender

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type NATSDataPublisherRepository struct {
	js nats.JetStreamContext
}

func NewNATSDataPublisherRepository(js nats.JetStreamContext) *NATSDataPublisherRepository {
	return &NATSDataPublisherRepository{
		js: js,
	}
}

type SensorDataDTO struct {
	SensorId  uuid.UUID       `json:"sensorId"`
	GatewayId uuid.UUID       `json:"gatewayId"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

func (r *NATSDataPublisherRepository) Send(d *sensorData) error {
	subject := fmt.Sprintf("sensor.data.%s.%s", d.GatewayId, d.SensorId)

	dto := SensorDataDTO{
		SensorId:  d.SensorId,
		GatewayId: d.GatewayId,
		Timestamp: d.Timestamp,
		Data:      json.RawMessage(d.Data),
	}

	marshaledSensorData, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("errore nel marshaling del dato in JSON: %w, gatewayId: %s, sensorId: %s, timestamp: %s", err, d.GatewayId.String(), d.SensorId.String(), d.Timestamp.String())
	}

	_, err = r.js.Publish(subject, marshaledSensorData)
	if err != nil {
		return err
	}

	return nil
}
