package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"go.uber.org/zap"
)

func (s *GatewayManagerService) DeleteGateway(cmdData *commanddata.DeleteGateway) Response {
	s.gateways.Mu.RLock()
	worker, exists := s.gateways.Workers[cmdData.GatewayId]
	s.gateways.Mu.RUnlock()

	if !exists {
		return Response{Success: false, Message: fmt.Sprintf("gateway con Id %s non trovato nello stato del gateway manager", cmdData.GatewayId)}
	}

	cmd := commands.NewDeleteGatewayCmd(cmdData, s.configPort, worker.Sender)
	worker.CmdChannel <- cmd

	if err := <-worker.ErrChannel; err != nil {
		s.logger.Error("Errore nell'eliminazione del gateway",
			zap.String("gatewayId", cmdData.GatewayId.String()),
			zap.Error(err),
		)
		return Response{Success: false, Message: err.Error()}
	}

	s.gateways.Mu.Lock()
	delete(s.gateways.Workers, cmdData.GatewayId)
	s.gateways.Mu.Unlock()

	s.sensors.Mu.Lock()
	sensorMap := s.sensors.Workers[cmdData.GatewayId]
	delete(s.sensors.Workers, cmdData.GatewayId)
	s.sensors.Mu.Unlock()

	for sensorId, sensorWorker := range sensorMap {
		stopCmd := commands.NewStopSensorCmd(sensorWorker.SimulatedSensor)
		sensorWorker.CmdChannel <- stopCmd
		if err := <-sensorWorker.ErrChannel; err != nil {
			s.logger.Error("Errore nell'arresto del sensore",
				zap.String("gatewayId", cmdData.GatewayId.String()),
				zap.String("sensorId", sensorId.String()),
				zap.Error(err),
			)
			return Response{Success: false, Message: err.Error()}
		}
	}

	return Response{Success: true, Message: "Gateway eliminato con successo"}
}
