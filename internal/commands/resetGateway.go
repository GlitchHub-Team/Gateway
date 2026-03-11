package commands

import (
	"fmt"

	configmanager "Gateway/internal/configManager"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type ResetGatewayCmd struct {
	cmdData        *commanddata.ResetGateway
	configService  configmanager.GatewayResetterPort
	gatewayWorkers *gatewaymanager.GatewayWorkers
}

func (c *ResetGatewayCmd) Execute() error {
	if err := c.configService.ResetGateway(c.cmdData); err != nil {
		return err
	}

	c.gatewayWorkers.Mu.RLock()
	worker, exists := c.gatewayWorkers.Workers[c.cmdData.GatewayId]
	c.gatewayWorkers.Mu.RUnlock()

	if !exists {
		return fmt.Errorf("nessun gateway trovato per il reset, id %s", c.cmdData.GatewayId)
	}

	c.gatewayWorkers.Mu.Lock()
	if err := worker.Reset(); err != nil {
		c.gatewayWorkers.Mu.Unlock()
		return err
	}
	c.gatewayWorkers.Mu.Unlock()

	return nil
}

func NewResetGatewayCmd(cmdData *commanddata.ResetGateway, configService configmanager.GatewayResetterPort, gateways *gatewaymanager.GatewayWorkers) *ResetGatewayCmd {
	return &ResetGatewayCmd{
		cmdData:        cmdData,
		configService:  configService,
		gatewayWorkers: gateways,
	}
}

func (c *ResetGatewayCmd) String() string {
	return "ResetGatewayCmd"
}
