package commands

import (
	configmanager "Gateway/internal/configManager"
)

type DecommissionGatewayCmd struct {
	gatewayId     string
	configService *configmanager.ConfigManagerService
}

func (c *DecommissionGatewayCmd) Execute() error {
	// Logic to decommission a gateway using the configService
	return nil
}

func NewDecommissionGatewayCmd(gatewayId string, configService *configmanager.ConfigManagerService) *DecommissionGatewayCmd {
	return &DecommissionGatewayCmd{
		gatewayId:     gatewayId,
		configService: configService,
	}
}
