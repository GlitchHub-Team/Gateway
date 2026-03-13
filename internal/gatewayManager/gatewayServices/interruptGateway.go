package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (s *GatewayManagerService) InterruptGateway(cmdData *commanddata.InterruptGateway) Response {
	s.gateways.Mu.RLock()
	worker, exists := s.gateways.Workers[cmdData.GatewayId]
	s.gateways.Mu.RUnlock()

	if !exists {
		return Response{Success: false, Message: fmt.Sprintf("nessun gateway trovato per l'interruzione, id %s", cmdData.GatewayId)}
	}

	cmd := commands.NewInterruptGatewayCmd(cmdData, worker.Sender, s.configPort, domain.Inactive)
	worker.CmdChannel <- cmd

	if err := <-worker.ErrChannel; err != nil {
		return Response{Success: false, Message: err.Error()}
	}

	return Response{Success: true, Message: "Gateway interrotto con successo"}
}
