package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type InterruptGatewayCmd struct {
	cmdData                *commanddata.InterruptGateway
	sender                 buffereddatasender.DataSenderInterrupter
	gatewayInterrupterPort configmanager.GatewayInterrupterPort
}

func (c *InterruptGatewayCmd) Execute() error {
	if err := c.gatewayInterrupterPort.InterruptGateway(c.cmdData); err != nil {
		return err
	}

	c.sender.Interrupt()
	return nil
}

func NewInterruptGatewayCmd(cmdData *commanddata.InterruptGateway, sender buffereddatasender.DataSenderInterrupter, gatewayInterrupterPort configmanager.GatewayInterrupterPort) *InterruptGatewayCmd {
	return &InterruptGatewayCmd{
		cmdData:                cmdData,
		sender:                 sender,
		gatewayInterrupterPort: gatewayInterrupterPort,
	}
}

func (c *InterruptGatewayCmd) String() string {
	return "InterruptGatewayCmd"
}
