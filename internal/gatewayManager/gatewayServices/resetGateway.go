package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (s *GatewayManagerService) ResetGateway(cmdData *commanddata.ResetGateway) Response {
	s.gateways.Mu.RLock()
	worker, exists := s.gateways.Workers[cmdData.GatewayId]
	s.gateways.Mu.RUnlock()

	if !exists {
		return Response{Success: false, Message: fmt.Sprintf("nessun gateway trovato per il reset, id %s", cmdData.GatewayId)}
	}

	cmd := commands.NewResetGatewayCmd(cmdData, worker.Sender, s.configPort)
	worker.CmdChannel <- cmd

	if err := <-worker.ErrChannel; err != nil {
		return Response{Success: false, Message: err.Error()}
	}

	return Response{Success: true, Message: "Gateway resettato con successo"}
}
