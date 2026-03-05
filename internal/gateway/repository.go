package gateway

import (
	"github.com/nats-io/nats.go"
)

type NATSCommandResponseRepository struct {
	natsMsg *nats.Msg
}

func (r *NATSCommandResponseRepository) Reply(response Response) error {
	// Logic to publish the response to a NATS subject
	err := r.natsMsg.Respond([]byte(response.Message))
	return err
}
