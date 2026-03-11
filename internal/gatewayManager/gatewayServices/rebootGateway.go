package gatewayservices

import (
	"fmt"
	"time"

	"Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

const defaultRebootDuration = 2 * time.Second

func (s *GatewayManagerService) RebootGateway(cmdData *commanddata.RebootGateway) Response {
	s.gateways.Mu.RLock()
	worker, exists := s.gateways.Workers[cmdData.GatewayId]
	s.gateways.Mu.RUnlock()

	if !exists {
		return Response{Success: false, Message: fmt.Sprintf("nessun gateway trovato per il riavvio, id %s", cmdData.GatewayId)}
	}

	cmd := commands.NewRebootGatewayCmd(*cmdData, defaultRebootDuration, s.ctx, s.logger)
	worker.CmdChannel <- cmd

	if err := <-worker.ErrChannel; err != nil {
		return Response{Success: false, Message: err.Error()}
	}

	return Response{Success: true, Message: "Gateway riavviato con successo"}
}
