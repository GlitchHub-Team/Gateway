package commands

import (
	configmanager "Gateway/internal/configManager"
	profiles "Gateway/internal/sensor/sensorProfiles"
)

type AddSensorCmd struct {
	gatewayId     string
	sensorId      string
	profile       profiles.SensorProfile
	frequency     int
	configService *configmanager.ConfigManagerService
}

func (c *AddSensorCmd) Execute() error {
	// Logic to add a new sensor using the configService
	return nil
}

func NewAddSensorCmd(gatewayId string, sensorId string, profile profiles.SensorProfile, frequency int, configService *configmanager.ConfigManagerService) *AddSensorCmd {
	return &AddSensorCmd{
		gatewayId:     gatewayId,
		sensorId:      sensorId,
		profile:       profile,
		frequency:     frequency,
		configService: configService,
	}
}
