package commands

import (
	"context"

	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AddSensorCmd struct {
	cmdData       *commanddata.AddSensor
	sensorAdder   configmanager.SensorAdderPort
	sensorWorkers *gatewaymanager.SensorWorkers
	bufferPort    sensor.SaveSensorDataPort
	ctx           context.Context
	logger        *zap.Logger
}

func (c *AddSensorCmd) Execute() error {
	if err := c.sensorAdder.AddSensor(c.cmdData); err != nil {
		return err
	}

	sensorEntity := &sensor.Sensor{
		Id:        c.cmdData.SensorId,
		GatewayId: c.cmdData.GatewayId,
		Profile:   c.cmdData.Profile,
		Interval:  c.cmdData.Interval,
		Status:    sensor.Active,
	}

	simulatedSensor := sensor.NewSensorService(
		sensorEntity,
		c.bufferPort,
		make(chan domain.BaseCommand),
		make(chan struct{}),
		c.ctx,
		c.logger,
	)

	c.sensorWorkers.Mu.Lock()
	if c.sensorWorkers.Workers[c.cmdData.GatewayId] == nil {
		c.sensorWorkers.Workers[c.cmdData.GatewayId] = make(map[uuid.UUID]sensor.SimulatedSensor)
	}
	c.sensorWorkers.Workers[c.cmdData.GatewayId][c.cmdData.SensorId] = simulatedSensor
	c.sensorWorkers.Mu.Unlock()

	go simulatedSensor.Start()

	return nil
}

func NewAddSensorCmd(cmdData *commanddata.AddSensor, sensorAdder configmanager.SensorAdderPort, sensorWorkers *gatewaymanager.SensorWorkers, bufferPort sensor.SaveSensorDataPort, ctx context.Context, logger *zap.Logger) *AddSensorCmd {
	return &AddSensorCmd{
		cmdData:       cmdData,
		sensorAdder:   sensorAdder,
		sensorWorkers: sensorWorkers,
		bufferPort:    bufferPort,
		ctx:           ctx,
		logger:        logger,
	}
}

func (c *AddSensorCmd) String() string {
	return "AddSensorCmd"
}
