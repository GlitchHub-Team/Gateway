package commands

import (
	configmanager "Gateway/internal/configManager"
)

type InterruptGatewayCmd struct {
	gatewayId     string
	configService *configmanager.ConfigManagerService
}

func (c *InterruptGatewayCmd) Execute() error {
	// Logic to interrupt a gateway using the configService
	return nil
}

func NewInterruptGatewayCmd(gatewayId string, configService *configmanager.ConfigManagerService) *InterruptGatewayCmd {
	return &InterruptGatewayCmd{
		gatewayId:     gatewayId,
		configService: configService,
	}
}
