package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
)

func TestRebootGatewayCmdExecute(t *testing.T) {
	// verifica che RebootGatewayCmd invii solo hello
	greeter := &mockGatewayGreeter{}

	cmd := commands.NewRebootGatewayCmd(greeter)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if greeter.helloCalls != 1 {
		t.Fatalf("expected Hello to be called once, got %d", greeter.helloCalls)
	}

	if cmd.String() != "RebootGatewayCmd" {
		t.Fatalf("expected String() to return RebootGatewayCmd, got %s", cmd.String())
	}
}

func TestRebootGatewayCmdExecuteReturnsHelloError(t *testing.T) {
	// verifica che RebootGatewayCmd propaghi l'errore di Hello
	expectedErr := errors.New("hello failed")
	greeter := &mockGatewayGreeter{err: expectedErr}

	cmd := commands.NewRebootGatewayCmd(greeter)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if greeter.helloCalls != 1 {
		t.Fatalf("expected Hello to be called once, got %d", greeter.helloCalls)
	}
}
