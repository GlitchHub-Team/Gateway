package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type ResumeGatewayCmd struct {
	cmdData     *commanddata.ResumeGateway
	resumerPort configmanager.GatewayResumerPort
	sender      buffereddatasender.DataSenderResumer
}

func (c *ResumeGatewayCmd) Execute() error {
	if err := c.resumerPort.ResumeGateway(c.cmdData); err != nil {
		return err
	}

	c.sender.Resume()
	return nil
}

func NewResumeGatewayCmd(cmdData *commanddata.ResumeGateway, sender buffereddatasender.DataSenderResumer, resumerPort configmanager.GatewayResumerPort) *ResumeGatewayCmd {
	return &ResumeGatewayCmd{
		cmdData:     cmdData,
		resumerPort: resumerPort,
		sender:      sender,
	}
}

func (c *ResumeGatewayCmd) String() string {
	return "ResumeGatewayCmd"
}
