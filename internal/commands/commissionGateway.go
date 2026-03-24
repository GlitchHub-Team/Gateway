package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type CommissionGatewayCmd struct {
	cmdData       *commanddata.CommissionGateway
	configService configmanager.GatewayCommissionerPort
	commissioner  buffereddatasender.DataSenderCommissioner
	status        domain.GatewayStatus
}

func (c *CommissionGatewayCmd) Execute() error {
	if err := c.commissioner.Commission(c.cmdData.TenantId, c.cmdData.CommissionedToken); err != nil {
		return err
	}

	if err := c.configService.CommissionGateway(c.cmdData, c.status); err != nil {
		return err
	}

	return nil
}

func NewCommissionGatewayCmd(cmdData *commanddata.CommissionGateway, configService configmanager.GatewayCommissionerPort, commissioner buffereddatasender.DataSenderCommissioner, status domain.GatewayStatus) *CommissionGatewayCmd {
	return &CommissionGatewayCmd{
		cmdData:       cmdData,
		configService: configService,
		commissioner:  commissioner,
		status:        status,
	}
}

func (c *CommissionGatewayCmd) String() string {
	return "CommissionGatewayCmd"
}
