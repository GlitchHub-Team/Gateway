package commands

import (
	"fmt"

	configmanager "Gateway/internal/configManager"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type DeleteSensorCmd struct {
	cmdData       *commanddata.DeleteSensor
	configService configmanager.SensorDeleterPort
	sensorWorkers *gatewaymanager.SensorWorkers
}

func (c *DeleteSensorCmd) Execute() error {
	if err := c.configService.DeleteSensor(c.cmdData); err != nil {
		return err
	}

	c.sensorWorkers.Mu.RLock()
	sensors, exists := c.sensorWorkers.Workers[c.cmdData.GatewayId]
	c.sensorWorkers.Mu.RUnlock()

	if !exists {
		return fmt.Errorf("gateway con Id %s non trovato nello stato del gateway manager", c.cmdData.GatewayId)
	}

	c.sensorWorkers.Mu.RLock()
	sensorWorker, sensorExists := sensors[c.cmdData.SensorId]
	c.sensorWorkers.Mu.RUnlock()

	if !sensorExists {
		return fmt.Errorf("sensor con Id %s non trovato nello stato del gateway manager", c.cmdData.SensorId)
	}

	sensorWorker.Stop()

	c.sensorWorkers.Mu.Lock()
	delete(sensors, c.cmdData.SensorId)
	c.sensorWorkers.Mu.Unlock()

	return nil
}

func NewDeleteSensorCmd(cmdData *commanddata.DeleteSensor, configService configmanager.SensorDeleterPort, sensorWorkers *gatewaymanager.SensorWorkers) *DeleteSensorCmd {
	return &DeleteSensorCmd{
		cmdData:       cmdData,
		configService: configService,
		sensorWorkers: sensorWorkers,
	}
}

func (c *DeleteSensorCmd) String() string {
	return "DeleteSensorCmd"
}
