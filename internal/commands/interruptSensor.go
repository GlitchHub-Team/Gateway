package commands

import (
	"fmt"

	configmanager "Gateway/internal/configManager"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type InterruptSensorCmd struct {
	cmdData       *commanddata.InterruptSensor
	configService configmanager.SensorInterrupterPort
	sensorWorkers *gatewaymanager.SensorWorkers
}

func (c *InterruptSensorCmd) Execute() error {
	if err := c.configService.InterruptSensor(c.cmdData); err != nil {
		return err
	}

	c.sensorWorkers.Mu.RLock()
	sensors, exists := c.sensorWorkers.Workers[c.cmdData.GatewayId]
	c.sensorWorkers.Mu.RUnlock()

	if !exists {
		return fmt.Errorf("nessun gateway trovato per l'interruzione del sensore, id gateway %s", c.cmdData.GatewayId)
	}

	c.sensorWorkers.Mu.RLock()
	sensorWorker, sensorExists := sensors[c.cmdData.SensorId]
	c.sensorWorkers.Mu.RUnlock()

	if !sensorExists {
		return fmt.Errorf("nessun sensore trovato per l'interruzione, id sensore %s", c.cmdData.SensorId)
	}

	c.sensorWorkers.Mu.Lock()
	sensorWorker.Interrupt()
	c.sensorWorkers.Mu.Unlock()

	return nil
}

func NewInterruptSensorCmd(cmdData *commanddata.InterruptSensor, configService configmanager.SensorInterrupterPort, sensorWorkers *gatewaymanager.SensorWorkers) *InterruptSensorCmd {
	return &InterruptSensorCmd{
		cmdData:       cmdData,
		configService: configService,
		sensorWorkers: sensorWorkers,
	}
}

func (c *InterruptSensorCmd) String() string {
	return "InterruptSensorCmd"
}
