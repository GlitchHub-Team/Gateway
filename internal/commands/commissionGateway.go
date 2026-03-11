package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type CommissionGatewayCmd struct {
	cmdData       *commanddata.CommissionGateway
	configService configmanager.GatewayCommissionerPort
	errChannel    chan error
}

func (c *CommissionGatewayCmd) Execute() error {
	// Logic to create a new gateway using the configService
	return nil
}

func NewCommissionGatewayCmd(cmdData *commanddata.CommissionGateway, configService configmanager.GatewayCommissionerPort, errChannel chan error) *CommissionGatewayCmd {
	return &CommissionGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
		errChannel:    errChannel,
	}
}

func (c *CommissionGatewayCmd) String() string {
	return "CommissionGatewayCmd"
}
