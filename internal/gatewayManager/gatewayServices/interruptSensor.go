package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
)

func (s *GatewayManagerService) InterruptSensor(cmdData *commanddata.InterruptSensor) Response {
	s.sensors.Mu.RLock()
	sensors, exists := s.sensors.Workers[cmdData.GatewayId]

	if !exists {
		s.sensors.Mu.RUnlock()
		return Response{Success: false, Message: fmt.Sprintf("nessun gateway trovato per l'interruzione del sensore, id gateway %s", cmdData.GatewayId)}
	}

	sensorWorker, sensorExists := sensors[cmdData.SensorId]
	s.sensors.Mu.RUnlock()

	if !sensorExists {
		return Response{Success: false, Message: fmt.Sprintf("nessun sensore trovato per l'interruzione, id sensore %s", cmdData.SensorId)}
	}

	cmd := commands.NewInterruptSensorCmd(cmdData, sensorWorker.SimulatedSensor, s.configPort, sensor.Inactive)
	sensorWorker.CmdChannel <- cmd

	if err := <-sensorWorker.ErrChannel; err != nil {
		return Response{Success: false, Message: err.Error()}
	}

	return Response{Success: true, Message: "Sensore interrotto con successo"}
}
