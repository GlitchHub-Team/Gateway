package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type DeleteGatewayCmd struct {
	cmdData       commanddata.DeleteGateway
	configService *configmanager.GatewayDeleter
}

func (c *DeleteGatewayCmd) Execute() error {
	// Logic to delete a gateway using the configService
	return nil
}

func NewDeleteGatewayCmd(cmdData commanddata.DeleteGateway, configService *configmanager.GatewayDeleter) *DeleteGatewayCmd {
	return &DeleteGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}

func (c *DeleteGatewayCmd) String() string {
	return "DeleteGatewayCmd"
}
