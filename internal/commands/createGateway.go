package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type CreateGatewayCmd struct {
	cmdData    *commanddata.CreateGateway
	configPort configmanager.GatewayCreatorPort
	sender     buffereddatasender.DataSenderStarter
}

func (c *CreateGatewayCmd) Execute() error {
	if err := c.configPort.CreateGateway(c.cmdData); err != nil {
		return err
	}

	if err := c.sender.Hello(); err != nil {
		return err
	}

	go c.sender.Start()

	return nil
}

func NewCreateGatewayCmd(cmdData *commanddata.CreateGateway, configPort configmanager.GatewayCreatorPort, sender buffereddatasender.DataSenderStarter) *CreateGatewayCmd {
	return &CreateGatewayCmd{
		cmdData:    cmdData,
		configPort: configPort,
		sender:     sender,
	}
}

func (c *CreateGatewayCmd) String() string {
	return "CreateGatewayCmd"
}
