package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
)

func TestDecommissionGatewayCmdExecute(t *testing.T) {
	//verifica che DecommissionGatewayCmd salvi il decommissioning, decommissioni e saluti
	cmdData := &commanddata.DecommissionGateway{GatewayId: uuid.New()}
	port := &mockGatewayDecommissionerPort{}
	decommissioner := &mockGatewayDecommissioner{}
	greeter := &mockGatewayGreeter{}

	cmd := commands.NewDecommissionGatewayCmd(cmdData, port, decommissioner, greeter, domain.Decommissioned)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !port.called {
		t.Fatal("expected DecommissionGateway to be called")
	}

	if port.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, port.receivedCmd)
	}

	if port.receivedStat != domain.Decommissioned {
		t.Fatalf("expected status %q, got %q", domain.Decommissioned, port.receivedStat)
	}

	if decommissioner.decommissionCalls != 1 {
		t.Fatalf("expected Decommission to be called once, got %d", decommissioner.decommissionCalls)
	}

	if greeter.helloCalls != 1 {
		t.Fatalf("expected Hello to be called once, got %d", greeter.helloCalls)
	}

	if cmd.String() != "DecommissionGatewayCmd" {
		t.Fatalf("expected String() to return DecommissionGatewayCmd, got %s", cmd.String())
	}
}

func TestDecommissionGatewayCmdExecuteReturnsDecommissioningError(t *testing.T) {
	//verifica che DecommissionGatewayCmd non chiami sender o greeter se il port fallisce
	expectedErr := errors.New("decommission failed")
	port := &mockGatewayDecommissionerPort{err: expectedErr}
	decommissioner := &mockGatewayDecommissioner{}
	greeter := &mockGatewayGreeter{}

	cmd := commands.NewDecommissionGatewayCmd(&commanddata.DecommissionGateway{}, port, decommissioner, greeter, domain.Decommissioned)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if decommissioner.decommissionCalls != 0 {
		t.Fatalf("expected Decommission not to be called, got %d", decommissioner.decommissionCalls)
	}

	if greeter.helloCalls != 0 {
		t.Fatalf("expected Hello not to be called, got %d", greeter.helloCalls)
	}
}

func TestDecommissionGatewayCmdExecuteReturnsSenderError(t *testing.T) {
	//verifica che DecommissionGatewayCmd non saluti se il sender fallisce
	expectedErr := errors.New("sender decommission failed")
	port := &mockGatewayDecommissionerPort{}
	decommissioner := &mockGatewayDecommissioner{err: expectedErr}
	greeter := &mockGatewayGreeter{}

	cmd := commands.NewDecommissionGatewayCmd(&commanddata.DecommissionGateway{}, port, decommissioner, greeter, domain.Decommissioned)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if greeter.helloCalls != 0 {
		t.Fatalf("expected Hello not to be called, got %d", greeter.helloCalls)
	}
}

func TestDecommissionGatewayCmdExecuteReturnsHelloError(t *testing.T) {
	//verifica che DecommissionGatewayCmd propaghi l'errore di Hello
	expectedErr := errors.New("hello failed")
	port := &mockGatewayDecommissionerPort{}
	decommissioner := &mockGatewayDecommissioner{}
	greeter := &mockGatewayGreeter{err: expectedErr}

	cmd := commands.NewDecommissionGatewayCmd(&commanddata.DecommissionGateway{}, port, decommissioner, greeter, domain.Decommissioned)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
