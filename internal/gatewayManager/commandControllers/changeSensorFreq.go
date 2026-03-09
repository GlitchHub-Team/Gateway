package commandcontrollers

import (
	"encoding/json"

	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"
	sensorprofiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type ChangeSensorFrequencySubject string

type NATSChangeSensorFrequencyController struct {
	natsConn *nats.Conn
	subject  string
	useCase  gatewayusecases.ChangeSensorFrequencyUseCase
	rand     sensorprofiles.Rand
	logger   *zap.Logger
}

type NATSChangeSensorFrequencyDTO struct {
	GatewayId string `json:"gatewayId"`
	Profile   string `json:"profile"`
	Frequency int    `json:"frequency"`
}

func (c *NATSChangeSensorFrequencyController) Listen() {
	_, err := c.natsConn.Subscribe(c.subject, func(msg *nats.Msg) {
		c.logger.Info("Ricevuto comando su subject: ", zap.String("subject", c.subject))
		cmd, err := c.parseChangeSensorFrequencyCommand(msg)
		if err != nil {
			err := wrongCommandErrorHandler(err, msg, c.logger)
			if err != nil {
				c.logger.Error("Errore durante la comunicazione del messaggio di formato errato", zap.String("subject", c.subject), zap.Error(err))
			}
			return
		}
		res := c.useCase.ChangeSensorFrequency(cmd)
		err = responseHandler(&res, msg)
		if err != nil {
			c.logger.Error("Errore durante la comunicazione della risposta", zap.String("subject", c.subject), zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("Errore nel subscribe: ", zap.String("subject", c.subject), zap.Error(err))
	}
}

func (c *NATSChangeSensorFrequencyController) parseChangeSensorFrequencyCommand(msg *nats.Msg) (*commanddata.ChangeSensorFrequency, error) {
	var req NATSChangeSensorFrequencyDTO

	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		return nil, err
	}

	gatewayId, err := uuid.Parse(req.GatewayId)
	if err != nil {
		return nil, err
	}

	profile := sensorprofiles.ParseSensorProfile(req.Profile, c.rand)
	if profile == nil {
		return nil, err
	}

	frequency := commanddata.SensorFrequency(req.Frequency)

	return &commanddata.ChangeSensorFrequency{
		GatewayId: gatewayId,
		Profile:   profile,
		Frequency: frequency,
	}, nil
}

func NewNATSChangeSensorFrequencyController(natsConn *nats.Conn, subject ChangeSensorFrequencySubject, useCase gatewayusecases.ChangeSensorFrequencyUseCase, rand sensorprofiles.Rand, logger *zap.Logger) *NATSChangeSensorFrequencyController {
	return &NATSChangeSensorFrequencyController{
		natsConn: natsConn,
		subject:  string(subject),
		useCase:  useCase,
		rand:     rand,
		logger:   logger,
	}
}
