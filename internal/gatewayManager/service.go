package gatewaymanager

import (
	"context"

	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GatewayId = uuid.UUID

type SensorId = uuid.UUID

type GatewayManagerService struct {
	gateways           map[GatewayId]buffereddatasender.DataSender
	sensors            map[GatewayId]map[SensorId]sensor.SimulatedSensor
	saveSensorDataPort sensor.SaveSensorDataPort
	bufferedDataPort   buffereddatasender.BufferedDataPort
	sendSensorDataPort buffereddatasender.SendSensorDataPort
	gatFetcher         configmanager.GatewaysFetcher
	ctx                context.Context
	logger             *zap.Logger
}

func NewGatewayManagerService(gateways map[GatewayId]buffereddatasender.DataSender, sensors map[GatewayId]map[SensorId]sensor.SimulatedSensor, saveSensorDataPort sensor.SaveSensorDataPort, bufferedDataPort buffereddatasender.BufferedDataPort, sendSensorDataPort buffereddatasender.SendSensorDataPort, gatFetcher configmanager.GatewaysFetcher, ctx context.Context, logger *zap.Logger) *GatewayManagerService {
	return &GatewayManagerService{
		gateways:           gateways,
		sensors:            sensors,
		saveSensorDataPort: saveSensorDataPort,
		bufferedDataPort:   bufferedDataPort,
		sendSensorDataPort: sendSensorDataPort,
		gatFetcher:         gatFetcher,
		ctx:                ctx,
		logger:             logger,
	}
}

func NewGatewayWorkers() map[GatewayId]buffereddatasender.DataSender {
	return make(map[GatewayId]buffereddatasender.DataSender)
}

func NewSensorWorkers() map[GatewayId]map[SensorId]sensor.SimulatedSensor {
	return make(map[GatewayId]map[SensorId]sensor.SimulatedSensor)
}

func (s *GatewayManagerService) CreateGateway(cmdData *commanddata.CreateGateway) Response {
	// Logic to create a gateway
	return Response{}
}

func (s *GatewayManagerService) CommissionGateway(cmdData *commanddata.CommissionGateway) Response {
	// Logic to commission a gateway
	return Response{}
}

func (s *GatewayManagerService) DecommissionGateway(cmdData *commanddata.DecommissionGateway) Response {
	// Logic to decommission a gateway
	return Response{}
}

func (s *GatewayManagerService) DeleteGateway(cmdData *commanddata.DeleteGateway) Response {
	// Logic to delete a gateway
	return Response{}
}

func (s *GatewayManagerService) InterruptGateway(cmdData *commanddata.InterruptGateway) Response {
	// Logic to interrupt a gateway
	return Response{}
}

func (s *GatewayManagerService) RebootGateway(cmdData *commanddata.RebootGateway) Response {
	// Logic to reboot a gateway
	return Response{}
}

func (s *GatewayManagerService) ResetGateway(cmdData *commanddata.ResetGateway) Response {
	// Logic to reset a gateway
	return Response{}
}

func (s *GatewayManagerService) ResumeGateway(cmdData *commanddata.ResumeGateway) Response {
	// Logic to resume a gateway
	return Response{}
}

func (s *GatewayManagerService) AddSensor(cmdData *commanddata.AddSensor) Response {
	// Logic to add a sensor to a gateway
	return Response{}
}

func (s *GatewayManagerService) DeleteSensor(cmdData *commanddata.DeleteSensor) Response {
	// Logic to delete a sensor from a gateway
	return Response{}
}

func (s *GatewayManagerService) InterruptSensor(cmdData *commanddata.InterruptSensor) Response {
	// Logic to interrupt a sensor
	return Response{}
}

func (s *GatewayManagerService) ResumeSensor(cmdData *commanddata.ResumeSensor) Response {
	// Logic to resume a sensor
	return Response{}
}

func (s *GatewayManagerService) ChangeSensorFrequency(cmdData *commanddata.ChangeSensorFrequency) Response {
	// Logic to change the frequency of a sensor
	return Response{}
}
