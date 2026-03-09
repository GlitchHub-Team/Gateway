package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type CreateGatewaySubject string

type NATSCreateGatewayController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.CreateGatewayUseCase
	logger   *zap.Logger
}

type NATSCreateGatewayDTO struct {
	GatewayId string `json:"gatewayId"`
}

func (c *NATSCreateGatewayController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		c.logger.Info("Ricevuto comando su subject: ", zap.String("subject", c.subject))
		cmd, err := c.parseCreateGatewayCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.CreateGateway(cmd)
		err = responseHandler(&res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSCreateGatewayController) parseCreateGatewayCommand(msg *nats.Msg) (*commanddata.CreateGateway, error) {
	var req NATSCreateGatewayDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	return &commanddata.CreateGateway{
		GatewayId: gatewayId,
	}, nil
}

func NewNATSCreateGatewayController(natsConn *nats.Conn, subject CreateGatewaySubject, useCase gatewayusecases.CreateGatewayUseCase, logger *zap.Logger) *NATSCreateGatewayController {
	return &NATSCreateGatewayController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
