package commands

import (
	configmanager "Gateway/internal/configManager"
)

type CreateGatewayCmd struct {
	gatewayName   string
	configService *configmanager.ConfigManagerService
}

func (c *CreateGatewayCmd) Execute() error {
	// Logic to create a new gateway using the configService
	return nil
}

func NewCreateGatewayCmd(gatewayName string, configService *configmanager.ConfigManagerService) *CreateGatewayCmd {
	return &CreateGatewayCmd{
		gatewayName:   gatewayName,
		configService: configService,
	}
}
