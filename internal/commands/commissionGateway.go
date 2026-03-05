package commands

import (
	configmanager "Gateway/internal/configManager"
)

type CommissionGatewayCmd struct {
	gatewayId     string
	tenantId      string
	configService *configmanager.ConfigManagerService
}

func (c *CommissionGatewayCmd) Execute() error {
	// Logic to create a new gateway using the configService
	return nil
}

func NewCommissionGatewayCmd(gatewayId string, tenantId string, configService *configmanager.ConfigManagerService) *CommissionGatewayCmd {
	return &CommissionGatewayCmd{
		gatewayId:     gatewayId,
		tenantId:      tenantId,
		configService: configService,
	}
}
