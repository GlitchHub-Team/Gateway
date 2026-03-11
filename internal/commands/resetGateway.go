package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"time"
)

type ResetGatewayCmd struct {
	cmdData             *commanddata.ResetGateway
	sender              buffereddatasender.DataSenderResetter
	gatewayResetterPort configmanager.GatewayResetterPort
}

const (
	defaultInterval = 5 * time.Second
)

func (c *ResetGatewayCmd) Execute() error {
	if err := c.gatewayResetterPort.ResetGateway(c.cmdData, defaultInterval); err != nil {
		return err
	}

	return c.sender.Reset(defaultInterval)
}

func NewResetGatewayCmd(cmdData *commanddata.ResetGateway, sender buffereddatasender.DataSenderResetter, gatewayResetterPort configmanager.GatewayResetterPort) *ResetGatewayCmd {
	return &ResetGatewayCmd{
		cmdData:             cmdData,
		sender:              sender,
		gatewayResetterPort: gatewayResetterPort,
	}
}

func (c *ResetGatewayCmd) String() string {
	return "ResetGatewayCmd"
}
