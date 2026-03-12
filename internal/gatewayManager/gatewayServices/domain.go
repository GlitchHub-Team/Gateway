package gatewayservices

import (
	"context"

	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	credentialsgenerator "Gateway/internal/credentialsGenerator"
	gatewaymanager "Gateway/internal/gatewayManager"
	"Gateway/internal/sensor"

	"go.uber.org/zap"
)

type Response struct {
	Success bool
	Message string
}

type GatewaysLoader interface {
	LoadGatewayWorkers() error
}

type GatewayManagerService struct {
	gateways                  gatewaymanager.GatewayWorkers
	sensors                   gatewaymanager.SensorWorkers
	saveSensorDataPort        sensor.SaveSensorDataPort
	bufferedDataPort          buffereddatasender.BufferedDataPort
	sendSensorDataPortFactory buffereddatasender.SendSensorDataPortFactory
	credentialsGenPort        credentialsgenerator.CredentialsGeneratorPort
	configPort                configmanager.ConfigPort
	ctx                       context.Context
	logger                    *zap.Logger
}

func NewGatewayManagerService(gateways gatewaymanager.GatewayWorkers, sensors gatewaymanager.SensorWorkers, saveSensorDataPort sensor.SaveSensorDataPort, bufferedDataPort buffereddatasender.BufferedDataPort, sendSensorDataPortFactory buffereddatasender.SendSensorDataPortFactory, credentialsGenPort credentialsgenerator.CredentialsGeneratorPort, configPort configmanager.ConfigPort, ctx context.Context, logger *zap.Logger) *GatewayManagerService {
	return &GatewayManagerService{
		gateways:                  gateways,
		sensors:                   sensors,
		saveSensorDataPort:        saveSensorDataPort,
		bufferedDataPort:          bufferedDataPort,
		sendSensorDataPortFactory: sendSensorDataPortFactory,
		credentialsGenPort:        credentialsGenPort,
		configPort:                configPort,
		ctx:                       ctx,
		logger:                    logger,
	}
}
