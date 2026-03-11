package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
)

type InterruptSensorCmd struct {
	cmdData               *commanddata.InterruptSensor
	simulatedSensor       sensor.SensorInterrupter
	sensorInterrupterPort configmanager.SensorInterrupterPort
}

func (c *InterruptSensorCmd) Execute() error {
	if err := c.sensorInterrupterPort.InterruptSensor(c.cmdData); err != nil {
		return err
	}

	c.simulatedSensor.Interrupt()
	return nil
}

func NewInterruptSensorCmd(cmdData *commanddata.InterruptSensor, simulatedSensor sensor.SensorInterrupter, interrupterPort configmanager.SensorInterrupterPort) *InterruptSensorCmd {
	return &InterruptSensorCmd{
		cmdData:               cmdData,
		simulatedSensor:       simulatedSensor,
		sensorInterrupterPort: interrupterPort,
	}
}

func (c *InterruptSensorCmd) String() string {
	return "InterruptSensorCmd"
}
