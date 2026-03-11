package commands

import "Gateway/internal/sensor"

type StopSensorCmd struct {
	sensorStopper sensor.SensorStopper
}

func (c *StopSensorCmd) Execute() error {
	c.sensorStopper.Stop()
	return nil
}

func NewStopSensorCmd(sensorStopper sensor.SensorStopper) *StopSensorCmd {
	return &StopSensorCmd{sensorStopper: sensorStopper}
}

func (c *StopSensorCmd) String() string {
	return "StopSensorCmd"
}
