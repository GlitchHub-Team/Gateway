package commands

import (
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type ResumeSensorCmd struct {
	cmdData       commanddata.ResumeSensor
	configService *configmanager.SensorResumerPort
	errChannel    chan error
}

func (c *ResumeSensorCmd) Execute() error {
	// Logic to resume a sensor using the configService
	return nil
}

func NewResumeSensorCmd(cmdData commanddata.ResumeSensor, configService *configmanager.SensorResumerPort, errChannel chan error) *ResumeSensorCmd {
	return &ResumeSensorCmd{
		cmdData:       cmdData,
		configService: configService,
		errChannel:    errChannel,
	}
}

func (c *ResumeSensorCmd) String() string {
	return "ResumeSensorCmd"
}
