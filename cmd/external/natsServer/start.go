package natsserver

import (
	"log"
	"strconv"
	"time"

	"Gateway/internal/natsutil"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func NewMockNATSConnection(address natsutil.NatsAddress, port natsutil.NatsPort, token natsutil.NatsToken, seed natsutil.NatsSeed) *nats.Conn {
	opts := &server.Options{
		Host:      string(address),
		Port:      int(port),
		JetStream: true,
	}
	s, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("Impossible to create a new NATS server: %v", err)
	}
	go s.Start()
	if !s.ReadyForConnections(10 * time.Second) {
		log.Fatalf("NATS server hasn't started on time")
	}
	log.Printf("NATS server is running on %s:%d", address, port)

	opt := natsutil.JWTAuth(string(token), string(seed))

	nc, err := nats.Connect("nats://"+string(address)+":"+strconv.Itoa(int(port)), opt)
	if err != nil {
		log.Fatalf("Error while connecting to NATS server: %v", err)
	}
	return nc
}
