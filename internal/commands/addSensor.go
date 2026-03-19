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
	status        sensor.SensorStatus
}

func (c *AddSensorCmd) Execute() error {
	if err := c.sensorAdder.AddSensor(c.cmdData, c.status); err != nil {
		return err
	}

	go c.sensorStarter.Start()

	return nil
}

func NewAddSensorCmd(cmdData *commanddata.AddSensor, sensorAdder configmanager.SensorAdderPort, sensorStarter sensor.SensorStarter, status sensor.SensorStatus) *AddSensorCmd {
	return &AddSensorCmd{
		cmdData:       cmdData,
		sensorAdder:   sensorAdder,
		sensorStarter: sensorStarter,
		status:        status,
	}
}

func (c *AddSensorCmd) String() string {
	return "AddSensorCmd"
}
