package natsserver

import (
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type (
	NatsAddress string
	NatsPort    int
)

func NewNATSConnection(address NatsAddress, port NatsPort) *nats.Conn {
	opts := &server.Options{
		Host: string(address),
		Port: int(port),
	}
	s, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("Impossible to create a new NATS server: %v", err)
	}
	go s.Start()
	if !s.ReadyForConnections(5 * time.Second) {
		log.Fatalf("NATS server hasn't started on time")
	}
	log.Printf("NATS server is running on %s:%d", address, port)

	nc, err := nats.Connect("nats://" + string(address) + ":" + strconv.Itoa(int(port)))
	if err != nil {
		log.Fatalf("Error while connecting to NATS server: %v", err)
	}
	return nc
}
