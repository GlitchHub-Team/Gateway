package buffereddatasender

import (
	"fmt"

	"Gateway/internal/natsutil"

	"github.com/nats-io/nats.go"
)

type (
	NatsAddress string
	NatsPort    int
	NatsToken   string
	NatsSeed    string
)

type NATSDataPublisherFactory struct {
	js      nats.JetStreamContext
	address NatsAddress
	port    NatsPort
}

func NewNATSDataPublisherFactory(js nats.JetStreamContext, address NatsAddress, port NatsPort) *NATSDataPublisherFactory {
	return &NATSDataPublisherFactory{
		js:      js,
		address: address,
		port:    port,
	}
}

func (f *NATSDataPublisherFactory) Create() SendSensorDataPort {
	return NewNATSDataPublisherRepository(f.js)
}

func (f *NATSDataPublisherFactory) Reload(token string, seed string) (SendSensorDataPort, error) {
	opt := natsutil.JWTAuth(token, seed)

	url := fmt.Sprintf("nats://%s:%d", f.address, f.port)
	nc, err := nats.Connect(url, opt)
	if err != nil {
		return nil, fmt.Errorf("errore nella connessione a NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("errore nell'ottenimento del contesto JetStream: %w", err)
	}

	return NewNATSDataPublisherRepository(js), nil
}
