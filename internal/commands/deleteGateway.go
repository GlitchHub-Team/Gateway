package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type DeleteGatewayCmd struct {
	cmdData       *commanddata.DeleteGateway
	configService configmanager.GatewayDeleterPort
	senderStopper buffereddatasender.DataSenderStopper
}

func (c *DeleteGatewayCmd) Execute() error {
	if err := c.configService.DeleteGateway(c.cmdData); err != nil {
		return err
	}

	if err := c.senderStopper.Stop(); err != nil {
		return err
	}

	return nil
}

func NewDeleteGatewayCmd(cmdData *commanddata.DeleteGateway, configService configmanager.GatewayDeleterPort, senderStopper buffereddatasender.DataSenderStopper) *DeleteGatewayCmd {
	return &DeleteGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
		senderStopper: senderStopper,
	}
}

func (c *DeleteGatewayCmd) String() string {
	return "DeleteGatewayCmd"
}
