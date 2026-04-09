package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
)

func TestDeleteGatewayCmdExecute(t *testing.T) {
	// verifica che DeleteGatewayCmd elimini il gateway e poi fermi il sender
	cmdData := &commanddata.DeleteGateway{GatewayId: uuid.New()}
	deleter := &mockGatewayDeleter{}
	stopper := &mockGatewayStopper{}

	cmd := commands.NewDeleteGatewayCmd(cmdData, deleter, stopper)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !deleter.called {
		t.Fatal("expected DeleteGateway to be called")
	}

	if deleter.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, deleter.receivedCmd)
	}

	if stopper.stopCalls != 1 {
		t.Fatalf("expected Stop to be called once, got %d", stopper.stopCalls)
	}

	if cmd.String() != "DeleteGatewayCmd" {
		t.Fatalf("expected String() to return DeleteGatewayCmd, got %s", cmd.String())
	}
}

func TestDeleteGatewayCmdExecuteReturnsDeleteError(t *testing.T) {
	// verifica che DeleteGatewayCmd non fermi il sender se la delete fallisce
	expectedErr := errors.New("delete failed")
	deleter := &mockGatewayDeleter{err: expectedErr}
	stopper := &mockGatewayStopper{}

	cmd := commands.NewDeleteGatewayCmd(&commanddata.DeleteGateway{}, deleter, stopper)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if stopper.stopCalls != 0 {
		t.Fatalf("expected Stop not to be called, got %d", stopper.stopCalls)
	}
}

func TestDeleteGatewayCmdExecuteReturnsStopError(t *testing.T) {
	// verifica che DeleteGatewayCmd propaghi l'errore di stop del sender
	expectedErr := errors.New("stop failed")
	deleter := &mockGatewayDeleter{}
	stopper := &mockGatewayStopper{err: expectedErr}

	cmd := commands.NewDeleteGatewayCmd(&commanddata.DeleteGateway{}, deleter, stopper)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if !deleter.called {
		t.Fatal("expected DeleteGateway to be called")
	}

	if stopper.stopCalls != 1 {
		t.Fatalf("expected Stop to be called once, got %d", stopper.stopCalls)
	}
}
