package commands

import buffereddatasender "Gateway/internal/bufferedDataSender"

type RebootGatewayCmd struct {
	greeter buffereddatasender.DataSenderGreeter
}

func (g *RebootGatewayCmd) Execute() error {
	return g.greeter.Hello()
}

func NewRebootGatewayCmd(greeter buffereddatasender.DataSenderGreeter) *RebootGatewayCmd {
	return &RebootGatewayCmd{greeter: greeter}
}

func (g *RebootGatewayCmd) String() string {
	return "RebootGatewayCmd"
}
