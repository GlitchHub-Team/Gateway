package commandstests

import (
	"testing"

	commands "Gateway/internal/commands"
)

func TestStopGatewayCmdExecute(t *testing.T) {
	//verifica che StopGatewayCmd fermi il sender
	stopper := &mockGatewayStopper{}
	cmd := commands.NewStopGatewayCmd(stopper)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if stopper.stopCalls != 1 {
		t.Fatalf("expected Stop to be called once, got %d", stopper.stopCalls)
	}

	if cmd.String() != "StopGatewayCmd" {
		t.Fatalf("expected String() to return StopGatewayCmd, got %s", cmd.String())
	}
}
