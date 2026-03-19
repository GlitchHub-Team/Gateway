package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	"Gateway/internal/domain"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"

	"go.uber.org/zap"
)

func (s *GatewayManagerService) AddSensor(cmdData *commanddata.AddSensor) Response {
	s.sensors.Mu.RLock()
	sensorMap, exists := s.sensors.Workers[cmdData.GatewayId]

	if !exists {
		return Response{Success: false, Message: fmt.Sprintf("gateway con Id %s non trovato", cmdData.GatewayId)}
	}

	_, exists = sensorMap[cmdData.SensorId]
	s.sensors.Mu.RUnlock()

	if exists {
		return Response{Success: false, Message: fmt.Sprintf("sensore con Id %s già presente nel gateway %s", cmdData.SensorId, cmdData.GatewayId)}
	}

	sensorEntity := &sensor.Sensor{
		Id:        cmdData.SensorId,
		GatewayId: cmdData.GatewayId,
		Profile:   cmdData.Profile,
		Interval:  cmdData.Interval,
		Status:    sensor.Active,
	}

	errChannel := make(chan error)
	cmdChannel := make(chan domain.BaseCommand)

	simulatedSensor := sensor.NewSensorService(
		sensorEntity,
		s.saveSensorDataPort,
		cmdChannel,
		errChannel,
		s.ctx,
		s.logger,
	)

	cmd := commands.NewAddSensorCmd(
		cmdData,
		s.configPort,
		simulatedSensor,
		sensor.Active,
	)

	if err := cmd.Execute(); err != nil {
		s.logger.Error("Errore nell'esecuzione del comando di aggiunta del sensore",
			zap.String("gatewayId", cmdData.GatewayId.String()),
			zap.String("sensorId", cmdData.SensorId.String()),
			zap.Error(err),
		)
		return Response{Success: false, Message: err.Error()}
	}

	s.addSensorToState(cmdData, simulatedSensor, errChannel, cmdChannel)

	return Response{Success: true, Message: "Sensore aggiunto con successo"}
}

func (s *GatewayManagerService) addSensorToState(cmdData *commanddata.AddSensor, simulatedSensor sensor.SimulatedSensor, errChannel chan error, cmdChannel chan domain.BaseCommand) {
	s.sensors.Mu.Lock()
	s.sensors.Workers[cmdData.GatewayId][cmdData.SensorId] = gatewaymanager.SensorWorker{
		SimulatedSensor: simulatedSensor,
		ErrChannel:      errChannel,
		CmdChannel:      cmdChannel,
	}
	s.sensors.Mu.Unlock()
}
