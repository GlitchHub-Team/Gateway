package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
)

func TestInterruptGatewayCmdExecute(t *testing.T) {
	// verifica che InterruptGatewayCmd aggiorni lo stato e interrompa il sender
	cmdData := &commanddata.InterruptGateway{GatewayId: uuid.New()}
	port := &mockGatewayInterrupterPort{}
	interrupter := &mockGatewayInterrupter{}

	cmd := commands.NewInterruptGatewayCmd(cmdData, interrupter, port, domain.Inactive)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !port.called {
		t.Fatal("expected InterruptGateway to be called")
	}

	if port.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, port.receivedCmd)
	}

	if port.receivedStat != domain.Inactive {
		t.Fatalf("expected status %q, got %q", domain.Inactive, port.receivedStat)
	}

	if interrupter.interruptCalls != 1 {
		t.Fatalf("expected Interrupt to be called once, got %d", interrupter.interruptCalls)
	}

	if cmd.String() != "InterruptGatewayCmd" {
		t.Fatalf("expected String() to return InterruptGatewayCmd, got %s", cmd.String())
	}
}

func TestInterruptGatewayCmdExecuteReturnsInterruptError(t *testing.T) {
	// verifica che InterruptGatewayCmd non interrompa il sender se il port fallisce
	expectedErr := errors.New("interrupt failed")
	port := &mockGatewayInterrupterPort{err: expectedErr}
	interrupter := &mockGatewayInterrupter{}

	cmd := commands.NewInterruptGatewayCmd(&commanddata.InterruptGateway{}, interrupter, port, domain.Inactive)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if interrupter.interruptCalls != 0 {
		t.Fatalf("expected Interrupt not to be called, got %d", interrupter.interruptCalls)
	}
}
