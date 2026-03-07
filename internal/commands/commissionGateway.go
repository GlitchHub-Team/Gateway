package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gateway/commandData"
)

type CommissionGatewayCmd struct {
	cmdData       commanddata.CommissionGateway
	configService *configmanager.GatewayCommissioner
}

func (c *CommissionGatewayCmd) Execute() error {
	// Logic to create a new gateway using the configService
	return nil
}

func NewCommissionGatewayCmd(cmdData commanddata.CommissionGateway, configService *configmanager.GatewayCommissioner) *CommissionGatewayCmd {
	return &CommissionGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}
