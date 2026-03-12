package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type DecommissionGatewayCmd struct {
	cmdData        *commanddata.DecommissionGateway
	configService  configmanager.GatewayDecommissionerPort
	decommissioner buffereddatasender.DataSenderDecommissioner
	greeter        buffereddatasender.DataSenderGreeter
}

func (c *DecommissionGatewayCmd) Execute() error {
	if err := c.configService.DecommissionGateway(c.cmdData); err != nil {
		return err
	}

	if err := c.decommissioner.Decommission(); err != nil {
		return err
	}

	if err := c.greeter.Hello(); err != nil {
		return err
	}

	return nil
}

func NewDecommissionGatewayCmd(cmdData *commanddata.DecommissionGateway, configService configmanager.GatewayDecommissionerPort, decommissioner buffereddatasender.DataSenderDecommissioner, greeter buffereddatasender.DataSenderGreeter) *DecommissionGatewayCmd {
	return &DecommissionGatewayCmd{
		cmdData:        cmdData,
		configService:  configService,
		decommissioner: decommissioner,
		greeter:        greeter,
	}
}

func (c *DecommissionGatewayCmd) String() string {
	return "DecommissionGatewayCmd"
}
