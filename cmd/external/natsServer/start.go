package natsserver

import (
	"log"
	"strconv"
	"time"

	buffereddatasender "Gateway/internal/bufferedDataSender"
	"Gateway/internal/natsutil"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func NewNATSConnection(address buffereddatasender.NatsAddress, port buffereddatasender.NatsPort, token buffereddatasender.NatsToken, seed buffereddatasender.NatsSeed) *nats.Conn {
	opts := &server.Options{
		Host: string(address),
		Port: int(port),
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

func NewJetStreamContext(nc *nats.Conn) nats.JetStreamContext {
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error while creating JetStream context: %v", err)
	}
	return js
}
