package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gateway/commandData"
)

type AddSensorCmd struct {
	cmdData       commanddata.AddSensor
	configService *configmanager.SensorAdder
}

func (c *AddSensorCmd) Execute() error {
	// Logic to add a new sensor using the configService
	return nil
}

func NewAddSensorCmd(cmdData commanddata.AddSensor, configService *configmanager.SensorAdder) *AddSensorCmd {
	return &AddSensorCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}
