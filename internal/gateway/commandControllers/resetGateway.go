package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gateway/commandData"
	gatewayusecases "Gateway/internal/gateway/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type ResetGatewaySubject string

type NATSResetGatewayController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.ResetGatewayUseCase
	logger   *zap.Logger
}

type NATSResetGatewayDTO struct {
	GatewayId string `json:"gatewayId"`
}

func (c *NATSResetGatewayController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		cmd, err := c.parseResetGatewayCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.ResetGateway(cmd)
		err = responseHandler(res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSResetGatewayController) parseResetGatewayCommand(msg *nats.Msg) (*commanddata.ResetGateway, error) {
	var req NATSResetGatewayDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	return &commanddata.ResetGateway{
		GatewayId: gatewayId,
	}, nil
}

func NewNATSResetGatewayController(natsConn *nats.Conn, subject ResetGatewaySubject, useCase gatewayusecases.ResetGatewayUseCase, logger *zap.Logger) *NATSResetGatewayController {
	return &NATSResetGatewayController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
