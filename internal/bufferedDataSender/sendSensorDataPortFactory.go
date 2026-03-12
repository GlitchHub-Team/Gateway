package buffereddatasender

import (
	"fmt"
	"strconv"

	"github.com/nats-io/nats.go"
)

func NewNATSDataPublisherFactory(address NatsAddress, port NatsPort, baseToken BaseToken, baseSeed BaseSeed) *NATSDataPublisherFactory {
	return &NATSDataPublisherFactory{
		address: address,
		port:    port,
		token:   baseToken,
		seed:    baseSeed,
	}
}

func (f *NATSDataPublisherFactory) Create() (SendSensorDataPort, error) {
	//TODO token e seed
	nc, err := nats.Connect("nats://" + string(f.address) + ":" + strconv.Itoa(int(f.port)))
	if err != nil {
		return nil, fmt.Errorf("errore nella connessione a NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("errore nell'ottenimento del contesto JetStream: %w", err)
	}

	return NewNATSDataPublisherRepository(js), nil
}

func (f *NATSDataPublisherFactory) Reload(token string, seed string) (SendSensorDataPort, error) {
	//TODO token e seed
	nc, err := nats.Connect("nats://" + string(f.address) + ":" + strconv.Itoa(int(f.port)))
	if err != nil {
		return nil, fmt.Errorf("errore nella connessione a NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("errore nell'ottenimento del contesto JetStream: %w", err)
	}

	return NewNATSDataPublisherRepository(js), nil
}
