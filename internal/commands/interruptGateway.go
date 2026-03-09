package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type InterruptGatewayCmd struct {
	cmdData       commanddata.InterruptGateway
	configService *configmanager.GatewayInterrupter
}

func (c *InterruptGatewayCmd) Execute() error {
	// Logic to interrupt a gateway using the configService
	return nil
}

func NewInterruptGatewayCmd(cmdData commanddata.InterruptGateway, configService *configmanager.GatewayInterrupter) *InterruptGatewayCmd {
	return &InterruptGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}

func (c *InterruptGatewayCmd) String() string {
	return "InterruptGatewayCmd"
}
