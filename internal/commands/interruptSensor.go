package commands

import (
	configmanager "Gateway/internal/configManager"
)

type InterruptSensorCmd struct {
	gatewayId     string
	sensorId      string
	configService *configmanager.ConfigManagerService
}

func (c *InterruptSensorCmd) Execute() error {
	// Logic to interrupt a sensor using the configService
	return nil
}

func NewInterruptSensorCmd(gatewayId string, sensorId string, configService *configmanager.ConfigManagerService) *InterruptSensorCmd {
	return &InterruptSensorCmd{
		gatewayId:     gatewayId,
		sensorId:      sensorId,
		configService: configService,
	}
}
