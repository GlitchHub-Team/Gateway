package commands

import (
	configmanager "Gateway/internal/configManager"
)

type ResumeSensorCmd struct {
	gatewayId     string
	sensorId      string
	configService *configmanager.ConfigManagerService
}

func (c *ResumeSensorCmd) Execute() error {
	// Logic to resume a sensor using the configService
	return nil
}

func NewResumeSensorCmd(gatewayId string, sensorId string, configService *configmanager.ConfigManagerService) *ResumeSensorCmd {
	return &ResumeSensorCmd{
		gatewayId:     gatewayId,
		sensorId:      sensorId,
		configService: configService,
	}
}
