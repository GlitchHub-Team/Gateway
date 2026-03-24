package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (s *GatewayManagerService) RebootGateway(cmdData *commanddata.RebootGateway) Response {
	s.gateways.Mu.RLock()
	worker, exists := s.gateways.Workers[cmdData.GatewayId]
	s.gateways.Mu.RUnlock()

	if !exists {
		return Response{Success: false, Message: fmt.Sprintf("nessun gateway trovato per il riavvio, id %s", cmdData.GatewayId)}
	}

	cmd := commands.NewRebootGatewayCmd(worker.Sender)
	worker.CmdChannel <- cmd

	if err := <-worker.ErrChannel; err != nil {
		return Response{Success: false, Message: err.Error()}
	}

	return Response{Success: true, Message: "Gateway riavviato con successo"}
}
