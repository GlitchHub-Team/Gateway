package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type DecommissionGatewaySubject string

type NATSDecommissionGatewayController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.DecommissionGatewayUseCase
	logger   *zap.Logger
}

type NATSDecommissionGatewayDTO struct {
	GatewayId string `json:"gatewayId"`
}

func (c *NATSDecommissionGatewayController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		c.logger.Info("Ricevuto comando su subject: ", zap.String("subject", c.subject))
		cmd, err := c.parseDecommissionGatewayCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.DecommissionGateway(cmd)
		err = responseHandler(&res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSDecommissionGatewayController) parseDecommissionGatewayCommand(msg *nats.Msg) (*commanddata.DecommissionGateway, error) {
	var req NATSDecommissionGatewayDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	return &commanddata.DecommissionGateway{
		GatewayId: gatewayId,
	}, nil
}

func NewNATSDecommissionGatewayController(natsConn *nats.Conn, subject DecommissionGatewaySubject, useCase gatewayusecases.DecommissionGatewayUseCase, logger *zap.Logger) *NATSDecommissionGatewayController {
	return &NATSDecommissionGatewayController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
