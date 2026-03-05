package commands

import (
	configmanager "Gateway/internal/configManager"
)

type DeleteGatewayCmd struct {
	gatewayName   string
	configService *configmanager.ConfigManagerService
}

func (c *DeleteGatewayCmd) Execute() error {
	// Logic to delete a gateway using the configService
	return nil
}

func NewDeleteGatewayCmd(gatewayName string, configService *configmanager.ConfigManagerService) *DeleteGatewayCmd {
	return &DeleteGatewayCmd{
		gatewayName:   gatewayName,
		configService: configService,
	}
}
