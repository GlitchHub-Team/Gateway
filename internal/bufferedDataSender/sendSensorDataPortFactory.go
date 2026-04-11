package buffereddatasender

import (
	"context"
	"fmt"
	"strings"

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
	normalizedToken := normalizeCommissionedToken(token)
	if normalizedToken == "" {
		return nil, fmt.Errorf("commissioned token vuoto")
	}

	options := make([]nats.Option, 0, 2)
	options = append(options, natsutil.JWTAuth(normalizedToken, seed))
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

func normalizeCommissionedToken(token string) string {
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		return ""
	}

	const beginJWTBlock = "-----BEGIN NATS USER JWT-----"
	const endJWTBlock = "------END NATS USER JWT------"

	beginIndex := strings.Index(trimmed, beginJWTBlock)
	if beginIndex == -1 {
		return trimmed
	}

	jwtSection := trimmed[beginIndex+len(beginJWTBlock):]
	endIndex := strings.Index(jwtSection, endJWTBlock)
	if endIndex >= 0 {
		jwtSection = jwtSection[:endIndex]
	}

	return strings.TrimSpace(jwtSection)
}
