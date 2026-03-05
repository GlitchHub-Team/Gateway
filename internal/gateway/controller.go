package gateway

import (
	"github.com/nats-io/nats.go"
)

type NatsCommandController struct {
	natsConn *nats.Conn
}

func (c *NatsCommandController) Listen() {
	_, err := c.natsConn.Subscribe("gateway.commands", func(msg *nats.Msg) {
		// Creerà una repository per rispondere al comando
		// repo := &NATSCommandResponseRepository{natsMsg: msg}
	})
	if err != nil {
		panic(err)
	}
}
