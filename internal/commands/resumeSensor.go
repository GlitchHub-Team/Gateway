package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gateway/commandData"
)

type ResumeSensorCmd struct {
	cmdData       commanddata.ResumeSensor
	configService *configmanager.SensorResumer
}

func (c *ResumeSensorCmd) Execute() error {
	// Logic to resume a sensor using the configService
	return nil
}

func NewResumeSensorCmd(cmdData commanddata.ResumeSensor, configService *configmanager.SensorResumer) *ResumeSensorCmd {
	return &ResumeSensorCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}
