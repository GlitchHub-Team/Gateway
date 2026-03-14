package commands

import buffereddatasender "Gateway/internal/bufferedDataSender"

type RebootGatewayCmd struct {
	stopper buffereddatasender.DataSenderStopper
	greeter buffereddatasender.DataSenderGreeter
	starter buffereddatasender.DataSenderStarter
}

func (g *RebootGatewayCmd) Execute() error {
	g.stopper.Stop()
	if err := g.greeter.Hello(); err != nil {
		return err
	}
	go g.starter.Start()
	return nil
}

func NewRebootGatewayCmd(stopper buffereddatasender.DataSenderStopper, greeter buffereddatasender.DataSenderGreeter, starter buffereddatasender.DataSenderStarter) *RebootGatewayCmd {
	return &RebootGatewayCmd{stopper: stopper, greeter: greeter, starter: starter}
}

func (g *RebootGatewayCmd) String() string {
	return "RebootGatewayCmd"
}
