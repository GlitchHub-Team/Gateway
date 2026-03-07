package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gateway/commandData"
)

type CreateGatewayCmd struct {
	cmdData       commanddata.CreateGateway
	configService *configmanager.GatewayCreator
}

func (c *CreateGatewayCmd) Execute() error {
	// Logic to create a new gateway using the configService
	return nil
}

func NewCreateGatewayCmd(cmdData commanddata.CreateGateway, configService *configmanager.GatewayCreator) *CreateGatewayCmd {
	return &CreateGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}
