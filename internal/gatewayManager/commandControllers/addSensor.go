package commandcontrollers

import (
	"encoding/json"
	"time"

	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"
	sensorprofiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type AddSensorSubject string

type NATSAddSensorController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.AddSensorUseCase
	rand     sensorprofiles.Rand
	logger   *zap.Logger
}

type NATSAddSensorDTO struct {
	GatewayId string `json:"gatewayId"`
	SensorId  string `json:"sensorId"`
	Profile   string `json:"profile"`
	Interval  int    `json:"interval"`
}

func (c *NATSAddSensorController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		c.logger.Info("Ricevuto comando su subject: ", zap.String("subject", c.subject))
		cmd, err := c.parseAddSensorCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.AddSensor(cmd)
		err = responseHandler(&res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSAddSensorController) parseAddSensorCommand(msg *nats.Msg) (*commanddata.AddSensor, error) {
	var req NATSAddSensorDTO

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

	profile := sensorprofiles.ParseSensorProfile(req.Profile, c.rand)
	if profile == nil {
		return nil, err
	}

	interval := time.Duration(req.Interval) * time.Millisecond
	if interval <= 0 {
		return nil, err
	}

	return &commanddata.AddSensor{
		GatewayId: gatewayId,
		SensorId:  sensorId,
		Profile:   profile,
		Interval:  interval,
	}, nil
}

func NewNATSAddSensorController(natsConn *nats.Conn, subject AddSensorSubject, useCase gatewayusecases.AddSensorUseCase, rand sensorprofiles.Rand, logger *zap.Logger) *NATSAddSensorController {
	return &NATSAddSensorController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		rand:     rand,
		logger:   logger,
	}
}
