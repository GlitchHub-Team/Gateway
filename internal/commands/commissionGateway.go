package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type CommissionGatewayCmd struct {
	cmdData       *commanddata.CommissionGateway
	configService configmanager.GatewayCommissionerPort
	commissioner  buffereddatasender.DataSenderCommissioner
}

func (c *CommissionGatewayCmd) Execute() error {
	if err := c.configService.CommissionGateway(c.cmdData); err != nil {
		return err
	}

	if err := c.commissioner.Commission(c.cmdData.TenantId, c.cmdData.CommissionedToken); err != nil {
		return err
	}

	return nil
}

func NewCommissionGatewayCmd(cmdData *commanddata.CommissionGateway, configService configmanager.GatewayCommissionerPort, commissioner buffereddatasender.DataSenderCommissioner) *CommissionGatewayCmd {
	return &CommissionGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
		commissioner:  commissioner,
	}
}

func (c *CommissionGatewayCmd) String() string {
	return "CommissionGatewayCmd"
}
