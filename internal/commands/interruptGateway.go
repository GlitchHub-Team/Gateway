package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type InterruptGatewayCmd struct {
	cmdData                *commanddata.InterruptGateway
	sender                 buffereddatasender.DataSenderInterrupter
	gatewayInterrupterPort configmanager.GatewayInterrupterPort
	status                 domain.GatewayStatus
}

func (c *InterruptGatewayCmd) Execute() error {
	if err := c.gatewayInterrupterPort.InterruptGateway(c.cmdData, c.status); err != nil {
		return err
	}

	c.sender.Interrupt()
	return nil
}

func NewInterruptGatewayCmd(cmdData *commanddata.InterruptGateway, sender buffereddatasender.DataSenderInterrupter, gatewayInterrupterPort configmanager.GatewayInterrupterPort, status domain.GatewayStatus) *InterruptGatewayCmd {
	return &InterruptGatewayCmd{
		cmdData:                cmdData,
		sender:                 sender,
		gatewayInterrupterPort: gatewayInterrupterPort,
		status:                 status,
	}
}

func (c *InterruptGatewayCmd) String() string {
	return "InterruptGatewayCmd"
}
