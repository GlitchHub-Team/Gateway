package commands

import (
	"fmt"

	configmanager "Gateway/internal/configManager"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type ResumeSensorCmd struct {
	cmdData       *commanddata.ResumeSensor
	configService configmanager.SensorResumerPort
	sensorWorkers *gatewaymanager.SensorWorkers
}

func (c *ResumeSensorCmd) Execute() error {
	if err := c.configService.ResumeSensor(c.cmdData); err != nil {
		return err
	}

	c.sensorWorkers.Mu.RLock()
	sensors, exists := c.sensorWorkers.Workers[c.cmdData.GatewayId]
	c.sensorWorkers.Mu.RUnlock()

	if !exists {
		return fmt.Errorf("nessun gateway trovato per la ripresa del sensore, id gateway %s", c.cmdData.GatewayId)
	}

	c.sensorWorkers.Mu.RLock()
	sensorWorker, sensorExists := sensors[c.cmdData.SensorId]
	c.sensorWorkers.Mu.RUnlock()

	if !sensorExists {
		return fmt.Errorf("nessun sensore trovato per la ripresa, id sensore %s", c.cmdData.SensorId)
	}

	c.sensorWorkers.Mu.Lock()
	sensorWorker.Resume()
	c.sensorWorkers.Mu.Unlock()

	return nil
}

func NewResumeSensorCmd(cmdData *commanddata.ResumeSensor, configService configmanager.SensorResumerPort, sensorWorkers *gatewaymanager.SensorWorkers) *ResumeSensorCmd {
	return &ResumeSensorCmd{
		cmdData:       cmdData,
		configService: configService,
		sensorWorkers: sensorWorkers,
	}
}

func (c *ResumeSensorCmd) String() string {
	return "ResumeSensorCmd"
}
