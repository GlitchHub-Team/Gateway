package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gateway/commandData"
)

type RebootGatewayCmd struct {
	cmdData       commanddata.RebootGateway
	configService *configmanager.GatewayRebooter
}

func (c *RebootGatewayCmd) Execute() error {
	// Logic to reboot a gateway using the configService
	return nil
}

func NewRebootGatewayCmd(cmdData commanddata.RebootGateway, configService *configmanager.GatewayRebooter) *RebootGatewayCmd {
	return &RebootGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}
