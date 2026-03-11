package commands

import (
	"fmt"

	configmanager "Gateway/internal/configManager"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type DeleteGatewayCmd struct {
	cmdData        *commanddata.DeleteGateway
	configService  configmanager.GatewayDeleterPort
	gatewayWorkers *gatewaymanager.GatewayWorkers
}

func (c *DeleteGatewayCmd) Execute() error {
	if err := c.configService.DeleteGateway(c.cmdData); err != nil {
		return err
	}

	c.gatewayWorkers.Mu.RLock()
	worker, exists := c.gatewayWorkers.Workers[c.cmdData.GatewayId]
	c.gatewayWorkers.Mu.RUnlock()

	if !exists {
		return fmt.Errorf("gateway con Id %s non trovato nello stato del gateway manager", c.cmdData.GatewayId)
	}

	worker.Stop()
	c.gatewayWorkers.Mu.Lock()
	delete(c.gatewayWorkers.Workers, c.cmdData.GatewayId)
	c.gatewayWorkers.Mu.Unlock()

	return nil
}

func NewDeleteGatewayCmd(cmdData *commanddata.DeleteGateway, configService configmanager.GatewayDeleterPort, gatewayWorkers *gatewaymanager.GatewayWorkers) *DeleteGatewayCmd {
	return &DeleteGatewayCmd{
		cmdData:        cmdData,
		configService:  configService,
		gatewayWorkers: gatewayWorkers,
	}
}

func (c *DeleteGatewayCmd) String() string {
	return "DeleteGatewayCmd"
}
