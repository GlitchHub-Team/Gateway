package commands

import (
	configmanager "Gateway/internal/configManager"
)

type DeleteSensorCmd struct {
	gatewayId     string
	sensorId      string
	configService *configmanager.ConfigManagerService
}

func (c *DeleteSensorCmd) Execute() error {
	// Logic to delete a sensor using the configService
	return nil
}

func NewDeleteSensorCmd(gatewayId string, sensorId string, configService *configmanager.ConfigManagerService) *DeleteSensorCmd {
	return &DeleteSensorCmd{
		gatewayId:     gatewayId,
		sensorId:      sensorId,
		configService: configService,
	}
}
