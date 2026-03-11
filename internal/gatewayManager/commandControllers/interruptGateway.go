package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type InterruptGatewaySubject string

type NATSInterruptGatewayController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.InterruptGatewayUseCase
	logger   *zap.Logger
}

type NATSInterruptGatewayDTO struct {
	GatewayId string `json:"gatewayId"`
}

func (c *NATSInterruptGatewayController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		c.logger.Info("Ricevuto comando su subject: ", zap.String("subject", c.subject))
		cmd, err := c.parseInterruptGatewayCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.InterruptGateway(cmd)
		err = responseHandler(&res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSInterruptGatewayController) parseInterruptGatewayCommand(msg *nats.Msg) (*commanddata.InterruptGateway, error) {
	var req NATSInterruptGatewayDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	return &commanddata.InterruptGateway{
		GatewayId: gatewayId,
	}, nil
}

func NewNATSInterruptGatewayController(natsConn *nats.Conn, subject InterruptGatewaySubject, useCase gatewayusecases.InterruptGatewayUseCase, logger *zap.Logger) *NATSInterruptGatewayController {
	return &NATSInterruptGatewayController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
