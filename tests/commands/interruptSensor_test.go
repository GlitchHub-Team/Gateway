package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensor "Gateway/internal/sensor"

	"github.com/google/uuid"
)

func TestInterruptSensorCmdExecute(t *testing.T) {
	//verifica che InterruptSensorCmd aggiorni lo stato e interrompa il sensore
	cmdData := &commanddata.InterruptSensor{GatewayId: uuid.New(), SensorId: uuid.New()}
	port := &mockSensorInterrupterPort{}
	interrupter := &mockSensorInterrupter{}

	cmd := commands.NewInterruptSensorCmd(cmdData, interrupter, port, sensor.Inactive)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !port.called {
		t.Fatal("expected InterruptSensor to be called")
	}

	if port.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, port.receivedCmd)
	}

	if port.receivedStat != sensor.Inactive {
		t.Fatalf("expected status %q, got %q", sensor.Inactive, port.receivedStat)
	}

	if interrupter.interruptCalls != 1 {
		t.Fatalf("expected Interrupt to be called once, got %d", interrupter.interruptCalls)
	}

	if cmd.String() != "InterruptSensorCmd" {
		t.Fatalf("expected String() to return InterruptSensorCmd, got %s", cmd.String())
	}
}

func TestInterruptSensorCmdExecuteReturnsInterruptError(t *testing.T) {
	//verifica che InterruptSensorCmd non interrompa il sensore se il port fallisce
	expectedErr := errors.New("interrupt failed")
	port := &mockSensorInterrupterPort{err: expectedErr}
	interrupter := &mockSensorInterrupter{}

	cmd := commands.NewInterruptSensorCmd(&commanddata.InterruptSensor{}, interrupter, port, sensor.Inactive)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if interrupter.interruptCalls != 0 {
		t.Fatalf("expected Interrupt not to be called, got %d", interrupter.interruptCalls)
	}
}
