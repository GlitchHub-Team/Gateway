package buffereddatasender

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSDataPublisherRepository struct {
	nc  *nats.Conn
	js  jetstream.JetStream
	ctx context.Context
}

func NewNATSDataPublisherRepository(nc *nats.Conn, js jetstream.JetStream, ctx context.Context) *NATSDataPublisherRepository {
	return &NATSDataPublisherRepository{
		nc:  nc,
		js:  js,
		ctx: ctx,
	}
}

type SensorDataDTO struct {
	SensorId  uuid.UUID       `json:"sensorId"`
	GatewayId uuid.UUID       `json:"gatewayId"`
	TenantId  uuid.UUID       `json:"tenantId"`
	Timestamp time.Time       `json:"timestamp"`
	Profile   string          `json:"profile"`
	Data      json.RawMessage `json:"data"`
}

type HelloMessageDTO struct {
	GatewayId        uuid.UUID `json:"gatewayId"`
	PublicIdentifier string    `json:"publicIdentifier"`
}

func (r *NATSDataPublisherRepository) Send(d *sensorData, tenantId uuid.UUID) error {
	subject := fmt.Sprintf("sensor.%s.%s", d.GatewayId, d.SensorId)

	dto := SensorDataDTO{
		SensorId:  d.SensorId,
		GatewayId: d.GatewayId,
		TenantId:  tenantId,
		Timestamp: d.Timestamp,
		Profile:   d.Profile,
		Data:      json.RawMessage(d.Data),
	}

	marshaledSensorData, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("errore nel marshaling del dato in JSON: %w, gatewayId: %s, sensorId: %s, timestamp: %s", err, d.GatewayId.String(), d.SensorId.String(), d.Timestamp.String())
	}

	err = r.nc.Publish(subject, marshaledSensorData)
	if err != nil {
		return err
	}

	return nil
}

func (r *NATSDataPublisherRepository) Hello(gatewayId uuid.UUID, publicIdentifier string) error {
	subject := fmt.Sprintf("gateway.hello.%s", gatewayId)

	dto := HelloMessageDTO{
		GatewayId:        gatewayId,
		PublicIdentifier: publicIdentifier,
	}

	marshaledHelloMessage, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("errore nel marshaling del messaggio di hello: %w, gatewayId: %s", err, gatewayId.String())
	}

	_, err = r.js.Publish(r.ctx, subject, marshaledHelloMessage)
	if err != nil {
		return fmt.Errorf("errore nell'invio del messaggio di hello: %w, gatewayId: %s", err, gatewayId.String())
	}

	return nil
}
