package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type ResumeGatewayCmd struct {
	cmdData       commanddata.ResumeGateway
	configService *configmanager.GatewayResumer
}

func (c *ResumeGatewayCmd) Execute() error {
	// Logic to resume a gateway using the configService
	return nil
}

func NewResumeGatewayCmd(cmdData commanddata.ResumeGateway, configService *configmanager.GatewayResumer) *ResumeGatewayCmd {
	return &ResumeGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}

func (c *ResumeGatewayCmd) String() string {
	return "ResumeGatewayCmd"
}
