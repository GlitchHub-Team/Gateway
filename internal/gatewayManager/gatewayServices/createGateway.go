package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"

	buffereddatasender "Gateway/internal/bufferedDataSender"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *GatewayManagerService) CreateGateway(cmdData *commanddata.CreateGateway) Response {
	s.gateways.Mu.RLock()
	_, exists := s.gateways.Workers[cmdData.GatewayId]
	s.gateways.Mu.RUnlock()

	if exists {
		return Response{Success: false, Message: fmt.Sprintf("gateway con Id %s già esistente", cmdData.GatewayId)}
	}

	gateway := &configmanager.Gateway{
		Id:       cmdData.GatewayId,
		TenantId: nil,
		Sensors:  make(map[uuid.UUID]*sensor.Sensor),
		Status:   configmanager.Decommissioned,
		Interval: cmdData.Interval,
	}

	errChannel := make(chan error)
	cmdChannel := make(chan domain.BaseCommand)

	dataSender := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		s.sendSensorDataPort,
		s.bufferedDataPort,
		cmdChannel,
		errChannel,
		s.ctx,
		s.logger,
	)

	cmd := commands.NewCreateGatewayCmd(cmdData, s.configPort, dataSender)
	if err := cmd.Execute(); err != nil {
		s.logger.Error("Errore nella creazione del gateway",
			zap.String("gatewayId", cmdData.GatewayId.String()),
			zap.Error(err),
		)
		return Response{Success: false, Message: err.Error()}
	}

	s.addGatewayToState(cmdData, dataSender, errChannel, cmdChannel)

	return Response{Success: true, Message: "Gateway creato con successo"}
}

func (s *GatewayManagerService) addGatewayToState(cmdData *commanddata.CreateGateway, dataSender buffereddatasender.DataSender, errChannel chan error, cmdChannel chan domain.BaseCommand) {
	s.gateways.Mu.Lock()
	s.gateways.Workers[cmdData.GatewayId] = gatewaymanager.GatewayWorker{
		Sender:     dataSender,
		ErrChannel: errChannel,
		CmdChannel: cmdChannel,
	}
	s.gateways.Mu.Unlock()

	s.sensors.Mu.Lock()
	if s.sensors.Workers[cmdData.GatewayId] == nil {
		s.sensors.Workers[cmdData.GatewayId] = make(map[uuid.UUID]gatewaymanager.SensorWorker)
	}
	s.sensors.Mu.Unlock()
}
