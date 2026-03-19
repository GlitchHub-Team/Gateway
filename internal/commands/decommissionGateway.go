package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type DecommissionGatewayCmd struct {
	cmdData        *commanddata.DecommissionGateway
	configService  configmanager.GatewayDecommissionerPort
	decommissioner buffereddatasender.DataSenderDecommissioner
	greeter        buffereddatasender.DataSenderGreeter
	status         domain.GatewayStatus
}

func (c *DecommissionGatewayCmd) Execute() error {
	if err := c.configService.DecommissionGateway(c.cmdData, c.status); err != nil {
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

func NewDecommissionGatewayCmd(cmdData *commanddata.DecommissionGateway, configService configmanager.GatewayDecommissionerPort, decommissioner buffereddatasender.DataSenderDecommissioner, greeter buffereddatasender.DataSenderGreeter, status domain.GatewayStatus) *DecommissionGatewayCmd {
	return &DecommissionGatewayCmd{
		cmdData:        cmdData,
		configService:  configService,
		decommissioner: decommissioner,
		greeter:        greeter,
		status:         status,
	}
}

func (c *DecommissionGatewayCmd) String() string {
	return "DecommissionGatewayCmd"
}
