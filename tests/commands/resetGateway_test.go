package commandstests

import (
	"errors"
	"testing"
	"time"

	commands "Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
)

func TestResetGatewayCmdExecute(t *testing.T) {
	// verifica che ResetGatewayCmd passi l'intervallo di default al port e al sender
	cmdData := &commanddata.ResetGateway{GatewayId: uuid.New()}
	port := &mockGatewayResetterPort{}
	resetter := &mockGatewayResetter{}

	cmd := commands.NewResetGatewayCmd(cmdData, resetter, port)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !port.called {
		t.Fatal("expected ResetGateway to be called")
	}

	if port.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, port.receivedCmd)
	}

	if port.receivedInterval != 5*time.Second {
		t.Fatalf("expected interval %v, got %v", 5*time.Second, port.receivedInterval)
	}

	if resetter.resetCalls != 1 {
		t.Fatalf("expected Reset to be called once, got %d", resetter.resetCalls)
	}

	if resetter.receivedInterval != 5*time.Second {
		t.Fatalf("expected interval %v, got %v", 5*time.Second, resetter.receivedInterval)
	}

	if cmd.String() != "ResetGatewayCmd" {
		t.Fatalf("expected String() to return ResetGatewayCmd, got %s", cmd.String())
	}
}

func TestResetGatewayCmdExecuteReturnsPortError(t *testing.T) {
	// verifica che ResetGatewayCmd non chiami il sender se il port fallisce
	expectedErr := errors.New("reset failed")
	port := &mockGatewayResetterPort{err: expectedErr}
	resetter := &mockGatewayResetter{}

	cmd := commands.NewResetGatewayCmd(&commanddata.ResetGateway{}, resetter, port)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if resetter.resetCalls != 0 {
		t.Fatalf("expected Reset not to be called, got %d", resetter.resetCalls)
	}
}

func TestResetGatewayCmdExecuteReturnsSenderError(t *testing.T) {
	// verifica che ResetGatewayCmd propaghi l'errore del sender
	expectedErr := errors.New("sender reset failed")
	port := &mockGatewayResetterPort{}
	resetter := &mockGatewayResetter{err: expectedErr}

	cmd := commands.NewResetGatewayCmd(&commanddata.ResetGateway{}, resetter, port)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
