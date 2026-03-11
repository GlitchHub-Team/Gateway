package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type ResumeGatewaySubject string

type NATSResumeGatewayController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.ResumeGatewayUseCase
	logger   *zap.Logger
}

type NATSResumeGatewayDTO struct {
	GatewayId string `json:"gatewayId"`
}

func (c *NATSResumeGatewayController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		c.logger.Info("Ricevuto comando su subject: ", zap.String("subject", c.subject))
		cmd, err := c.parseResumeGatewayCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.ResumeGateway(cmd)
		err = responseHandler(&res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSResumeGatewayController) parseResumeGatewayCommand(msg *nats.Msg) (*commanddata.ResumeGateway, error) {
	var req NATSResumeGatewayDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	return &commanddata.ResumeGateway{
		GatewayId: gatewayId,
	}, nil
}

func NewNATSResumeGatewayController(natsConn *nats.Conn, subject ResumeGatewaySubject, useCase gatewayusecases.ResumeGatewayUseCase, logger *zap.Logger) *NATSResumeGatewayController {
	return &NATSResumeGatewayController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
