package buffereddatasender

import (
	"context"
	"fmt"

	"Gateway/internal/natsutil"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	natsConnect         = nats.Connect
	newJetStreamContext = natsutil.NewJetStreamContext
)

type NATSDataPublisherFactory struct {
	nc        *nats.Conn
	js        jetstream.JetStream
	caPemPath natsutil.NatsCAPemPath
	address   natsutil.NatsAddress
	port      natsutil.NatsPort
	ctx       context.Context
}

func NewNATSDataPublisherFactory(js jetstream.JetStream, nc *nats.Conn, address natsutil.NatsAddress, port natsutil.NatsPort, ctx context.Context, caPemPath natsutil.NatsCAPemPath) *NATSDataPublisherFactory {
	return &NATSDataPublisherFactory{
		nc:        nc,
		js:        js,
		caPemPath: caPemPath,
		address:   address,
		port:      port,
		ctx:       ctx,
	}
}

func (f *NATSDataPublisherFactory) Create() SendSensorDataPort {
	return NewNATSDataPublisherRepository(f.nc, f.js, f.ctx)
}

func (f *NATSDataPublisherFactory) Reload(token string, seed string) (SendSensorDataPort, error) {
	options := make([]nats.Option, 0, 2)
	options = append(options, natsutil.JWTAuth(token, seed))
	options = append(options, natsutil.CAPemAuth(string(f.caPemPath)))

	url := fmt.Sprintf("nats://%s:%d", f.address, f.port)
	nc, err := natsConnect(url, options...)
	if err != nil {
		return nil, fmt.Errorf("errore nella connessione a NATS: %w", err)
	}

	js, err := newJetStreamContext(nc)
	if err != nil {
		return nil, fmt.Errorf("errore nella creazione del contesto JetStream: %w", err)
	}

	return NewNATSDataPublisherRepository(nc, js, f.ctx), nil
}
