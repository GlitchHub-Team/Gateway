package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gateway/commandData"
)

type DeleteSensorCmd struct {
	cmdData       commanddata.DeleteSensor
	configService *configmanager.SensorDeleter
}

func (c *DeleteSensorCmd) Execute() error {
	// Logic to delete a sensor using the configService
	return nil
}

func NewDeleteSensorCmd(cmdData commanddata.DeleteSensor, configService *configmanager.SensorDeleter) *DeleteSensorCmd {
	return &DeleteSensorCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}
