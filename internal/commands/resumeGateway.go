package commands

import (
	"fmt"

	configmanager "Gateway/internal/configManager"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type ResumeGatewayCmd struct {
	cmdData        *commanddata.ResumeGateway
	configService  configmanager.GatewayResumerPort
	gatewayWorkers *gatewaymanager.GatewayWorkers
}

func (c *ResumeGatewayCmd) Execute() error {
	if err := c.configService.ResumeGateway(c.cmdData); err != nil {
		return err
	}

	c.gatewayWorkers.Mu.RLock()
	worker, exists := c.gatewayWorkers.Workers[c.cmdData.GatewayId]
	c.gatewayWorkers.Mu.RUnlock()

	if !exists {
		return fmt.Errorf("nessun gateway trovato per la ripresa, id %s", c.cmdData.GatewayId)
	}

	c.gatewayWorkers.Mu.Lock()
	worker.Resume()
	c.gatewayWorkers.Mu.Unlock()

	return nil
}

func NewResumeGatewayCmd(cmdData *commanddata.ResumeGateway, configService configmanager.GatewayResumerPort, gatewayWorkers *gatewaymanager.GatewayWorkers) *ResumeGatewayCmd {
	return &ResumeGatewayCmd{
		cmdData:        cmdData,
		configService:  configService,
		gatewayWorkers: gatewayWorkers,
	}
}

func (c *ResumeGatewayCmd) String() string {
	return "ResumeGatewayCmd"
}
