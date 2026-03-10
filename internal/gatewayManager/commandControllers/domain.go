package commandcontrollers

import (
	gatewayservices "Gateway/internal/gatewayManager/gatewayServices"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NATSCommandController interface {
	Listen()
}

func wrongCommandErrorHandler(err error, msg *nats.Msg, logger *zap.Logger) error {
	res, err := json.Marshal(gatewayservices.Response{Success: false, Message: fmt.Sprintf("Formato del comando incorretto: %v", err)})
	if err != nil {
		logger.Panic("Errore durante la serializzazione del messaggio di formato errato", zap.Error(err))
	}
	err = msg.Respond(res)
	if err != nil {
		return err
	}
	return nil
}

func responseHandler(res *gatewayservices.Response, msg *nats.Msg) error {
	resBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}
	err = msg.Respond(resBytes)
	if err != nil {
		return err
	}
	return nil
}
