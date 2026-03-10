package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type ResetGatewayCmd struct {
	cmdData       commanddata.ResetGateway
	configService *configmanager.GatewayResetterPort
	errChannel    chan error
}

func (c *ResetGatewayCmd) Execute() error {
	// Logic to reset a gateway using the configService
	return nil
}

func NewResetGatewayCmd(cmdData commanddata.ResetGateway, configService *configmanager.GatewayResetterPort, errChannel chan error) *ResetGatewayCmd {
	return &ResetGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
		errChannel:    errChannel,
	}
}

func (c *ResetGatewayCmd) String() string {
	return "ResetGatewayCmd"
}
