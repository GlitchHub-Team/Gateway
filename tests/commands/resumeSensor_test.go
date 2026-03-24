package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensor "Gateway/internal/sensor"

	"github.com/google/uuid"
)

func TestResumeSensorCmdExecute(t *testing.T) {
	// verifica che ResumeSensorCmd aggiorni lo stato e riprenda il sensore
	cmdData := &commanddata.ResumeSensor{GatewayId: uuid.New(), SensorId: uuid.New()}
	port := &mockSensorResumerPort{}
	resumer := &mockSensorResumer{}

	cmd := commands.NewResumeSensorCmd(cmdData, resumer, port, sensor.Active)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !port.called {
		t.Fatal("expected ResumeSensor to be called")
	}

	if port.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, port.receivedCmd)
	}

	if port.receivedStat != sensor.Active {
		t.Fatalf("expected status %q, got %q", sensor.Active, port.receivedStat)
	}

	if resumer.resumeCalls != 1 {
		t.Fatalf("expected Resume to be called once, got %d", resumer.resumeCalls)
	}

	if cmd.String() != "ResumeSensorCmd" {
		t.Fatalf("expected String() to return ResumeSensorCmd, got %s", cmd.String())
	}
}

func TestResumeSensorCmdExecuteReturnsResumeError(t *testing.T) {
	// verifica che ResumeSensorCmd non riprenda il sensore se il port fallisce
	expectedErr := errors.New("resume failed")
	port := &mockSensorResumerPort{err: expectedErr}
	resumer := &mockSensorResumer{}

	cmd := commands.NewResumeSensorCmd(&commanddata.ResumeSensor{}, resumer, port, sensor.Active)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if resumer.resumeCalls != 0 {
		t.Fatalf("expected Resume not to be called, got %d", resumer.resumeCalls)
	}
}
