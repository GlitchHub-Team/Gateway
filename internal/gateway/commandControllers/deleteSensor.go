package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gateway/commandData"
	gatewayusecases "Gateway/internal/gateway/gatewayUseCases"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type DeleteSensorSubject string

type NATSDeleteSensorController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.DeleteSensorUseCase
	logger   *zap.Logger
}

type NATSDeleteSensorDTO struct {
	GatewayId string `json:"gatewayId"`
	SensorId  string `json:"sensorId"`
}

func (c *NATSDeleteSensorController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		cmd, err := c.parseDeleteSensorCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.DeleteSensor(cmd)
		err = responseHandler(res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSDeleteSensorController) parseDeleteSensorCommand(msg *nats.Msg) (*commanddata.DeleteSensor, error) {
	var req NATSDeleteSensorDTO

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

	return &commanddata.DeleteSensor{
		GatewayId: gatewayId,
		SensorId:  sensorId,
	}, nil
}

func NewNATSDeleteSensorController(natsConn *nats.Conn, subject DeleteSensorSubject, useCase gatewayusecases.DeleteSensorUseCase, logger *zap.Logger) *NATSDeleteSensorController {
	return &NATSDeleteSensorController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		logger:   logger,
	}
}
