package commands

import (
	configmanager "Gateway/internal/configManager"
)

type RebootGatewayCmd struct {
	gatewayId     string
	configService *configmanager.ConfigManagerService
}

func (c *RebootGatewayCmd) Execute() error {
	// Logic to reboot a gateway using the configService
	return nil
}

func NewRebootGatewayCmd(gatewayId string, configService *configmanager.ConfigManagerService) *RebootGatewayCmd {
	return &RebootGatewayCmd{
		gatewayId:     gatewayId,
		configService: configService,
	}
}
