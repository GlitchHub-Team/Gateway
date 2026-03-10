package gatewayservices

import (
	"Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

func (s *GatewayManagerService) AddSensor(cmdData *commanddata.AddSensor) Response {
	//controllo che il gatewayId esista, se no errore
	//controllo che non ci sia un sensore nel gateway con lo stesso id
	commands.NewAddSensorCmd(
		cmdData,
		s.configPort,
		&s.sensors,
		s.saveSensorDataPort,
		s.ctx,
		s.logger,
	)
	return Response{}
}
