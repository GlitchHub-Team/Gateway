package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
)

func TestDeleteSensorCmdExecute(t *testing.T) {
	// verifica che DeleteSensorCmd elimini il sensore e poi lo fermi
	cmdData := &commanddata.DeleteSensor{GatewayId: uuid.New(), SensorId: uuid.New()}
	deleter := &mockSensorDeleter{}
	stopper := &mockSensorStopper{}

	cmd := commands.NewDeleteSensorCmd(cmdData, deleter, stopper)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !deleter.called {
		t.Fatal("expected DeleteSensor to be called")
	}

	if deleter.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, deleter.receivedCmd)
	}

	if stopper.stopCalls != 1 {
		t.Fatalf("expected Stop to be called once, got %d", stopper.stopCalls)
	}

	if cmd.String() != "DeleteSensorCmd" {
		t.Fatalf("expected String() to return DeleteSensorCmd, got %s", cmd.String())
	}
}

func TestDeleteSensorCmdExecuteReturnsDeleteError(t *testing.T) {
	// verifica che DeleteSensorCmd non fermi il sensore se la delete fallisce
	expectedErr := errors.New("delete failed")
	deleter := &mockSensorDeleter{err: expectedErr}
	stopper := &mockSensorStopper{}

	cmd := commands.NewDeleteSensorCmd(&commanddata.DeleteSensor{}, deleter, stopper)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if stopper.stopCalls != 0 {
		t.Fatalf("expected Stop not to be called, got %d", stopper.stopCalls)
	}
}
