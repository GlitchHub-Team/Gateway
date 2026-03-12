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
	status                sensor.SensorStatus
}

func (c *InterruptSensorCmd) Execute() error {
	if err := c.sensorInterrupterPort.InterruptSensor(c.cmdData, c.status); err != nil {
		return err
	}

	c.simulatedSensor.Interrupt()
	return nil
}

func NewInterruptSensorCmd(cmdData *commanddata.InterruptSensor, simulatedSensor sensor.SensorInterrupter, interrupterPort configmanager.SensorInterrupterPort, status sensor.SensorStatus) *InterruptSensorCmd {
	return &InterruptSensorCmd{
		cmdData:               cmdData,
		simulatedSensor:       simulatedSensor,
		sensorInterrupterPort: interrupterPort,
		status:                status,
	}
}

func (c *InterruptSensorCmd) String() string {
	return "InterruptSensorCmd"
}
