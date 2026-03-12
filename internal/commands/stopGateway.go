package commands

import buffereddatasender "Gateway/internal/bufferedDataSender"

type StopGatewayCmd struct {
	gatewayStopper buffereddatasender.DataSenderStopper
}

func (g *StopGatewayCmd) Execute() error {
	g.gatewayStopper.Stop()
	return nil
}

func NewStopGatewayCmd(gatewayStopper buffereddatasender.DataSenderStopper) *StopGatewayCmd {
	return &StopGatewayCmd{gatewayStopper: gatewayStopper}
}

func (g *StopGatewayCmd) String() string {
	return "StopGatewayCmd"
}
