package commandstests

import (
	"testing"

	commands "Gateway/internal/commands"
)

func TestStopSensorCmdExecute(t *testing.T) {
	// verifica che StopSensorCmd fermi il sensore
	stopper := &mockSensorStopper{}
	cmd := commands.NewStopSensorCmd(stopper)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if stopper.stopCalls != 1 {
		t.Fatalf("expected Stop to be called once, got %d", stopper.stopCalls)
	}

	if cmd.String() != "StopSensorCmd" {
		t.Fatalf("expected String() to return StopSensorCmd, got %s", cmd.String())
	}
}
