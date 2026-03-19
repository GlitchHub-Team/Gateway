package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
)

type DeleteSensorCmd struct {
	cmdData       *commanddata.DeleteSensor
	configService configmanager.SensorDeleterPort
	sensorStopper sensor.SensorStopper
}

func (c *DeleteSensorCmd) Execute() error {
	if err := c.configService.DeleteSensor(c.cmdData); err != nil {
		return err
	}

	c.sensorStopper.Stop()

	return nil
}

func NewDeleteSensorCmd(cmdData *commanddata.DeleteSensor, configService configmanager.SensorDeleterPort, sensorStopper sensor.SensorStopper) *DeleteSensorCmd {
	return &DeleteSensorCmd{
		cmdData:       cmdData,
		configService: configService,
		sensorStopper: sensorStopper,
	}
}

func (c *DeleteSensorCmd) String() string {
	return "DeleteSensorCmd"
}
