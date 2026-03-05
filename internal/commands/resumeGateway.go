package commands

import (
	configmanager "Gateway/internal/configManager"
)

type ResumeGatewayCmd struct {
	gatewayId     string
	configService *configmanager.ConfigManagerService
}

func (c *ResumeGatewayCmd) Execute() error {
	// Logic to resume a gateway using the configService
	return nil
}

func NewResumeGatewayCmd(gatewayId string, configService *configmanager.ConfigManagerService) *ResumeGatewayCmd {
	return &ResumeGatewayCmd{
		gatewayId:     gatewayId,
		configService: configService,
	}
}
