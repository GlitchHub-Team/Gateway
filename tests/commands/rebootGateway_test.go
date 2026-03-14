package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
)

func TestRebootGatewayCmdExecute(t *testing.T) {
	//verifica che RebootGatewayCmd fermi il sender, saluti e lo riavvii
	stopper := &mockGatewayStopper{}
	greeter := &mockGatewayGreeter{}
	starter := &mockGatewayStarter{started: make(chan struct{}, 1)}

	cmd := commands.NewRebootGatewayCmd(stopper, greeter, starter)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if stopper.stopCalls != 1 {
		t.Fatalf("expected Stop to be called once, got %d", stopper.stopCalls)
	}

	if greeter.helloCalls != 1 {
		t.Fatalf("expected Hello to be called once, got %d", greeter.helloCalls)
	}

	waitForSignal(t, starter.started, "gateway start")

	if cmd.String() != "RebootGatewayCmd" {
		t.Fatalf("expected String() to return RebootGatewayCmd, got %s", cmd.String())
	}
}

func TestRebootGatewayCmdExecuteReturnsHelloError(t *testing.T) {
	//verifica che RebootGatewayCmd non avvii il sender se Hello fallisce
	expectedErr := errors.New("hello failed")
	stopper := &mockGatewayStopper{}
	greeter := &mockGatewayGreeter{err: expectedErr}
	starter := &mockGatewayStarter{started: make(chan struct{}, 1)}

	cmd := commands.NewRebootGatewayCmd(stopper, greeter, starter)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if stopper.stopCalls != 1 {
		t.Fatalf("expected Stop to be called once, got %d", stopper.stopCalls)
	}

	select {
	case <-starter.started:
		t.Fatal("expected Start not to be called on hello error")
	default:
	}
}
