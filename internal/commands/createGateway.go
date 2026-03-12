package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	credentialsgenerator "Gateway/internal/credentialsGenerator"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type CreateGatewayCmd struct {
	cmdData     *commanddata.CreateGateway
	configPort  configmanager.GatewayCreatorPort
	sender      buffereddatasender.DataSenderStarter
	greeter     buffereddatasender.DataSenderGreeter
	credentials *credentialsgenerator.Credentials
}

func (c *CreateGatewayCmd) Execute() error {
	if err := c.configPort.CreateGateway(c.cmdData, c.credentials); err != nil {
		return err
	}

	if err := c.greeter.Hello(); err != nil {
		return err
	}

	go c.sender.Start()

	return nil
}

func NewCreateGatewayCmd(cmdData *commanddata.CreateGateway, configPort configmanager.GatewayCreatorPort, sender buffereddatasender.DataSenderStarter, greeter buffereddatasender.DataSenderGreeter, credentials *credentialsgenerator.Credentials) *CreateGatewayCmd {
	return &CreateGatewayCmd{
		cmdData:     cmdData,
		configPort:  configPort,
		sender:      sender,
		greeter:     greeter,
		credentials: credentials,
	}
}

func (c *CreateGatewayCmd) String() string {
	return "CreateGatewayCmd"
}
