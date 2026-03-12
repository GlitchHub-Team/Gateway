package gatewayservices

import (
	"fmt"

	"Gateway/internal/commands"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (s *GatewayManagerService) DecommissionGateway(cmdData *commanddata.DecommissionGateway) Response {
	s.gateways.Mu.RLock()
	worker, exists := s.gateways.Workers[cmdData.GatewayId]
	s.gateways.Mu.RUnlock()

	if !exists {
		return Response{Success: false, Message: fmt.Sprintf("nessun gateway trovato per il decommissioning, id %s", cmdData.GatewayId)}
	}

	cmd := commands.NewDecommissionGatewayCmd(cmdData, s.configPort, worker.Sender, worker.Sender, domain.Decommissioned)
	worker.CmdChannel <- cmd

	if err := <-worker.ErrChannel; err != nil {
		return Response{Success: false, Message: err.Error()}
	}

	return Response{Success: true, Message: "Gateway decommissionato correttamente"}
}
