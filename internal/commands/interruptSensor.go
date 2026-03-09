package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type InterruptSensorCmd struct {
	cmdData       commanddata.InterruptSensor
	configService *configmanager.SensorInterrupter
}

func (c *InterruptSensorCmd) Execute() error {
	// Logic to interrupt a sensor using the configService
	return nil
}

func NewInterruptSensorCmd(cmdData commanddata.InterruptSensor, configService *configmanager.SensorInterrupter) *InterruptSensorCmd {
	return &InterruptSensorCmd{
		cmdData:       cmdData,
		configService: configService,
	}
}

func (c *InterruptSensorCmd) String() string {
	return "InterruptSensorCmd"
}
