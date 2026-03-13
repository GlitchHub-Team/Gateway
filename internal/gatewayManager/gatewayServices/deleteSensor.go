package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"go.uber.org/zap"
)

func (s *GatewayManagerService) DeleteSensor(cmdData *commanddata.DeleteSensor) Response {
	s.sensors.Mu.RLock()
	sensors, exists := s.sensors.Workers[cmdData.GatewayId]

	if !exists {
		s.sensors.Mu.RUnlock()
		return Response{Success: false, Message: fmt.Sprintf("gateway con Id %s non trovato nello stato del gateway manager", cmdData.GatewayId)}
	}

	sensorWorker, sensorExists := sensors[cmdData.SensorId]
	s.sensors.Mu.RUnlock()

	if !sensorExists {
		return Response{Success: false, Message: fmt.Sprintf("sensore con Id %s non trovato nello stato del gateway manager", cmdData.SensorId)}
	}

	cmd := commands.NewDeleteSensorCmd(cmdData, s.configPort, sensorWorker.SimulatedSensor)
	sensorWorker.CmdChannel <- cmd

	if err := <-sensorWorker.ErrChannel; err != nil {
		s.logger.Error("Errore nell'eliminazione del sensore",
			zap.String("gatewayId", cmdData.GatewayId.String()),
			zap.String("sensorId", cmdData.SensorId.String()),
			zap.Error(err),
		)
		return Response{Success: false, Message: err.Error()}
	}

	s.sensors.Mu.Lock()
	delete(sensors, cmdData.SensorId)
	s.sensors.Mu.Unlock()

	return Response{Success: true, Message: "Sensore eliminato con successo"}
}
