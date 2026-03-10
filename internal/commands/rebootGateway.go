package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type RebootGatewayCmd struct {
	cmdData       commanddata.RebootGateway
	configService *configmanager.GatewayRebooterPort
	errChannel    chan error
}

func (c *RebootGatewayCmd) Execute() error {
	// Logic to reboot a gateway using the configService
	return nil
}

func NewRebootGatewayCmd(cmdData commanddata.RebootGateway, configService *configmanager.GatewayRebooterPort, errChannel chan error) *RebootGatewayCmd {
	return &RebootGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
		errChannel:    errChannel,
	}
}

func (c *RebootGatewayCmd) String() string {
	return "RebootGatewayCmd"
}
