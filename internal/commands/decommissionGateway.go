package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type DecommissionGatewayCmd struct {
	cmdData       *commanddata.DecommissionGateway
	configService configmanager.GatewayDecommissionerPort
	errChannel    chan error
}

func (c *DecommissionGatewayCmd) Execute() error {
	// Logic to decommission a gateway using the configService
	return nil
}

func NewDecommissionGatewayCmd(cmdData *commanddata.DecommissionGateway, configService configmanager.GatewayDecommissionerPort, errChannel chan error) *DecommissionGatewayCmd {
	return &DecommissionGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
		errChannel:    errChannel,
	}
}

func (c *DecommissionGatewayCmd) String() string {
	return "DecommissionGatewayCmd"
}
