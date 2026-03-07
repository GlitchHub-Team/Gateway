package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gateway/commandData"
	gatewayusecases "Gateway/internal/gateway/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type DeleteGatewaySubject string

type NATSDeleteGatewayController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.DeleteGatewayUseCase
	logger   *zap.Logger
}

type NATSDeleteGatewayDTO struct {
	GatewayId string `json:"gatewayId"`
}

func (c *NATSDeleteGatewayController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		cmd, err := c.parseDeleteGatewayCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.DeleteGateway(cmd)
		err = responseHandler(res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSDeleteGatewayController) parseDeleteGatewayCommand(msg *nats.Msg) (*commanddata.DeleteGateway, error) {
	var req NATSDeleteGatewayDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	return &commanddata.DeleteGateway{
		GatewayId: gatewayId,
	}, nil
}

func NewNATSDeleteGatewayController(natsConn *nats.Conn, subject DeleteGatewaySubject, useCase gatewayusecases.DeleteGatewayUseCase, logger *zap.Logger) *NATSDeleteGatewayController {
	return &NATSDeleteGatewayController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
