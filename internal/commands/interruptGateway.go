package commands

import (
	"fmt"

	configmanager "Gateway/internal/configManager"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type InterruptGatewayCmd struct {
	cmdData        *commanddata.InterruptGateway
	configService  configmanager.GatewayInterrupterPort
	gatewayWorkers *gatewaymanager.GatewayWorkers
}

func (c *InterruptGatewayCmd) Execute() error {
	if err := c.configService.InterruptGateway(c.cmdData); err != nil {
		return err
	}

	c.gatewayWorkers.Mu.RLock()
	worker, exists := c.gatewayWorkers.Workers[c.cmdData.GatewayId]
	c.gatewayWorkers.Mu.RUnlock()

	if !exists {
		return fmt.Errorf("nessun gateway trovato per l'interruzione, id %s", c.cmdData.GatewayId)
	}

	c.gatewayWorkers.Mu.Lock()
	worker.Interrupt()
	c.gatewayWorkers.Mu.Unlock()

	return nil
}

func NewInterruptGatewayCmd(cmdData *commanddata.InterruptGateway, configService configmanager.GatewayInterrupterPort, gatewayWorkers *gatewaymanager.GatewayWorkers) *InterruptGatewayCmd {
	return &InterruptGatewayCmd{
		cmdData:        cmdData,
		configService:  configService,
		gatewayWorkers: gatewayWorkers,
	}
}

func (c *InterruptGatewayCmd) String() string {
	return "InterruptGatewayCmd"
}
