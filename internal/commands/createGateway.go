package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	credentialsgenerator "Gateway/internal/credentialsGenerator"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type CreateGatewayCmd struct {
	cmdData     *commanddata.CreateGateway
	configPort  configmanager.GatewayCreatorPort
	sender      buffereddatasender.DataSenderStarter
	greeter     buffereddatasender.DataSenderGreeter
	credentials *credentialsgenerator.Credentials
	status      domain.GatewayStatus
}

func (c *CreateGatewayCmd) Execute() error {
	if err := c.configPort.CreateGateway(c.cmdData, c.credentials, c.status); err != nil {
		return err
	}

	if err := c.greeter.Hello(); err != nil {
		return err
	}

	go c.sender.Start()

	return nil
}

func NewCreateGatewayCmd(cmdData *commanddata.CreateGateway, configPort configmanager.GatewayCreatorPort, sender buffereddatasender.DataSenderStarter, greeter buffereddatasender.DataSenderGreeter, credentials *credentialsgenerator.Credentials, status domain.GatewayStatus) *CreateGatewayCmd {
	return &CreateGatewayCmd{
		cmdData:     cmdData,
		configPort:  configPort,
		sender:      sender,
		greeter:     greeter,
		credentials: credentials,
		status:      status,
	}
}

func (c *CreateGatewayCmd) String() string {
	return "CreateGatewayCmd"
}
