package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type RebootGatewaySubject string

type NATSRebootGatewayController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.RebootGatewayUseCase
	logger   *zap.Logger
}

type NATSRebootGatewayDTO struct {
	GatewayId string `json:"gatewayId"`
}

func (c *NATSRebootGatewayController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		c.logger.Info("Ricevuto comando su subject: ", zap.String("subject", c.subject))
		cmd, err := c.parseRebootGatewayCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.RebootGateway(cmd)
		err = responseHandler(&res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSRebootGatewayController) parseRebootGatewayCommand(msg *nats.Msg) (*commanddata.RebootGateway, error) {
	var req NATSRebootGatewayDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	return &commanddata.RebootGateway{
		GatewayId: gatewayId,
	}, nil
}

func NewNATSRebootGatewayController(natsConn *nats.Conn, subject RebootGatewaySubject, useCase gatewayusecases.RebootGatewayUseCase, logger *zap.Logger) *NATSRebootGatewayController {
	return &NATSRebootGatewayController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
