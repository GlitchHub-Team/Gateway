package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
)

type AddSensorCmd struct {
	cmdData       *commanddata.AddSensor
	sensorAdder   configmanager.SensorAdderPort
	sensorStarter sensor.SensorStarter
}

func (c *AddSensorCmd) Execute() error {
	if err := c.sensorAdder.AddSensor(c.cmdData); err != nil {
		return err
	}

	go c.sensorStarter.Start()

	return nil
}

func NewAddSensorCmd(cmdData *commanddata.AddSensor, sensorAdder configmanager.SensorAdderPort, sensorStarter sensor.SensorStarter) *AddSensorCmd {
	return &AddSensorCmd{
		cmdData:       cmdData,
		sensorAdder:   sensorAdder,
		sensorStarter: sensorStarter,
	}
}

func (c *AddSensorCmd) String() string {
	return "AddSensorCmd"
}
