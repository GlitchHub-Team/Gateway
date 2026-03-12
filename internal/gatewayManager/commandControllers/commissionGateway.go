package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type CommissionGatewaySubject string

type NATSCommissionGatewayController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.CommissionGatewayUseCase
	logger   *zap.Logger
}

type NATSCommissionGatewayDTO struct {
	GatewayId         string `json:"gatewayId"`
	TenantId          string `json:"tenantId"`
	CommissionedToken string `json:"commissionedToken"`
}

func (c *NATSCommissionGatewayController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		c.logger.Info("Ricevuto comando su subject: ", zap.String("subject", c.subject))
		cmd, err := c.parseCommissionGatewayCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.CommissionGateway(cmd)
		err = responseHandler(&res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSCommissionGatewayController) parseCommissionGatewayCommand(msg *nats.Msg) (*commanddata.CommissionGateway, error) {
	var req NATSCommissionGatewayDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	tenantId, err := uuid.Parse(req.TenantId)
	if err != nil {
		return nil, err
	}

	return &commanddata.CommissionGateway{
		GatewayId:         gatewayId,
		TenantId:          tenantId,
		CommissionedToken: req.CommissionedToken,
	}, nil
}

func NewNATSCommissionGatewayController(natsConn *nats.Conn, subject CommissionGatewaySubject, useCase gatewayusecases.CommissionGatewayUseCase, logger *zap.Logger) *NATSCommissionGatewayController {
	return &NATSCommissionGatewayController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
