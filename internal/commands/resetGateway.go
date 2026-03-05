package commands

import (
	configmanager "Gateway/internal/configManager"
)

type ResetGatewayCmd struct {
	gatewayId     string
	configService *configmanager.ConfigManagerService
}

func (c *ResetGatewayCmd) Execute() error {
	// Logic to reset a gateway using the configService
	return nil
}

func NewResetGatewayCmd(gatewayId string, configService *configmanager.ConfigManagerService) *ResetGatewayCmd {
	return &ResetGatewayCmd{
		gatewayId:     gatewayId,
		configService: configService,
	}
}
