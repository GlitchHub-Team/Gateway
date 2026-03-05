package commands

import (
	configmanager "Gateway/internal/configManager"
	profiles "Gateway/internal/sensor/sensorProfiles"
)

type ChangeSensorFrequencyCmd struct {
	gatewayId     string
	profile       profiles.SensorProfile
	configService *configmanager.ConfigManagerService
}

func (c *ChangeSensorFrequencyCmd) Execute() error {
	// Logic to change sensor frequency using the configService
	return nil
}

func NewChangeSensorFrequencyCmd(gatewayId string, profile profiles.SensorProfile, configService *configmanager.ConfigManagerService) *ChangeSensorFrequencyCmd {
	return &ChangeSensorFrequencyCmd{
		gatewayId:     gatewayId,
		profile:       profile,
		configService: configService,
	}
}
