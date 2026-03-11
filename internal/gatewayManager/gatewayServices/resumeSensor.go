package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (s *GatewayManagerService) ResumeSensor(cmdData *commanddata.ResumeSensor) Response {
	s.sensors.Mu.RLock()
	sensors, exists := s.sensors.Workers[cmdData.GatewayId]

	if !exists {
		s.sensors.Mu.RUnlock()
		return Response{Success: false, Message: fmt.Sprintf("nessun gateway trovato per la ripresa del sensore, id gateway %s", cmdData.GatewayId)}
	}

	sensorWorker, sensorExists := sensors[cmdData.SensorId]
	s.sensors.Mu.RUnlock()

	if !sensorExists {
		return Response{Success: false, Message: fmt.Sprintf("nessun sensore trovato per la ripresa, id sensore %s", cmdData.SensorId)}
	}

	cmd := commands.NewResumeSensorCmd(cmdData, sensorWorker.SimulatedSensor, s.configPort)
	sensorWorker.CmdChannel <- cmd

	if err := <-sensorWorker.ErrChannel; err != nil {
		return Response{Success: false, Message: err.Error()}
	}

	return Response{Success: true, Message: "Sensore ripreso con successo"}
}
