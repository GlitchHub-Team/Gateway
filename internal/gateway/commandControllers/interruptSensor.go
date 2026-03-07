package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gateway/commandData"
	gatewayusecases "Gateway/internal/gateway/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type InterruptSensorSubject string

type NATSInterruptSensorController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.InterruptSensorUseCase
	logger   *zap.Logger
}

type NATSInterruptSensorDTO struct {
	GatewayId string `json:"gatewayId"`
	SensorId  string `json:"sensorId"`
}

func (c *NATSInterruptSensorController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		cmd, err := c.parseInterruptSensorCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.InterruptSensor(cmd)
		err = responseHandler(res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSInterruptSensorController) parseInterruptSensorCommand(msg *nats.Msg) (*commanddata.InterruptSensor, error) {
	var req NATSInterruptSensorDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	sensorId, err := uuid.Parse(req.SensorId)
	if err != nil {
		return nil, err
	}

	return &commanddata.InterruptSensor{
		GatewayId: gatewayId,
		SensorId:  sensorId,
	}, nil
}

func NewNATSInterruptSensorController(natsConn *nats.Conn, subject InterruptSensorSubject, useCase gatewayusecases.InterruptSensorUseCase, logger *zap.Logger) *NATSInterruptSensorController {
	return &NATSInterruptSensorController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
