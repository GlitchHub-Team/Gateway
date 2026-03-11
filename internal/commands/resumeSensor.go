package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
)

type ResumeSensorCmd struct {
	cmdData         *commanddata.ResumeSensor
	resumerPort     configmanager.SensorResumerPort
	simulatedSensor sensor.SensorResumer
}

func (c *ResumeSensorCmd) Execute() error {
	if err := c.resumerPort.ResumeSensor(c.cmdData); err != nil {
		return err
	}
	c.simulatedSensor.Resume()
	return nil
}

func NewResumeSensorCmd(cmdData *commanddata.ResumeSensor, simulatedSensor sensor.SensorResumer, resumerPort configmanager.SensorResumerPort) *ResumeSensorCmd {
	return &ResumeSensorCmd{
		cmdData:         cmdData,
		simulatedSensor: simulatedSensor,
		resumerPort:     resumerPort,
	}
}

func (c *ResumeSensorCmd) String() string {
	return "ResumeSensorCmd"
}
