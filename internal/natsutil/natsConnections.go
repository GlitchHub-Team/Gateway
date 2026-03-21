package natsutil

import (
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NewNATSConnection(address NatsAddress, port NatsPort, credsPath NatsCredsPath, caPemPath NatsCAPemPath) *nats.Conn {
	options := make([]nats.Option, 0, 2)
	if string(credsPath) != "" {
		options = append(options, CredsFileAuth(string(credsPath)))
	}
	options = append(options, CAPemAuth(string(caPemPath)))

	nc, err := nats.Connect("nats://"+string(address)+":"+strconv.Itoa(int(port)), options...)
	if err != nil {
		log.Fatalf("Error while connecting to NATS server: %v", err)
	}

	return nc
}

func NewJetStreamContext(nc *nats.Conn) (jetstream.JetStream, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	return js, nil
}
