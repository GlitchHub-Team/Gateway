package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gateway/commandData"
)

type DecommissionGatewayCmd struct {
	cmdData       commanddata.DecommissionGateway
	configService *configmanager.GatewayDecommissioner
}

func (c *DecommissionGatewayCmd) Execute() error {
	// Logic to decommission a gateway using the configService
	return nil
}

func NewDecommissionGatewayCmd(cmdData commanddata.DecommissionGateway, configService *configmanager.GatewayDecommissioner) *DecommissionGatewayCmd {
	return &DecommissionGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}
