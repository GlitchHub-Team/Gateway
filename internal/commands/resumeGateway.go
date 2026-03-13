package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type ResumeGatewayCmd struct {
	cmdData     *commanddata.ResumeGateway
	resumerPort configmanager.GatewayResumerPort
	sender      buffereddatasender.DataSenderResumer
	status      domain.GatewayStatus
}

func (c *ResumeGatewayCmd) Execute() error {
	if err := c.resumerPort.ResumeGateway(c.cmdData, c.status); err != nil {
		return err
	}

	c.sender.Resume()
	return nil
}

func NewResumeGatewayCmd(cmdData *commanddata.ResumeGateway, sender buffereddatasender.DataSenderResumer, resumerPort configmanager.GatewayResumerPort, status domain.GatewayStatus) *ResumeGatewayCmd {
	return &ResumeGatewayCmd{
		cmdData:     cmdData,
		resumerPort: resumerPort,
		sender:      sender,
		status:      status,
	}
}

func (c *ResumeGatewayCmd) String() string {
	return "ResumeGatewayCmd"
}
