package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gateway/commandData"
)

type ChangeSensorFrequencyCmd struct {
	cmdData       commanddata.ChangeSensorFrequency
	configService *configmanager.SensorFrequencySetter
}

func (c *ChangeSensorFrequencyCmd) Execute() error {
	// Logic to change sensor frequency using the configService
	return nil
}

func NewChangeSensorFrequencyCmd(cmdData commanddata.ChangeSensorFrequency, configService *configmanager.SensorFrequencySetter) *ChangeSensorFrequencyCmd {
	return &ChangeSensorFrequencyCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}
