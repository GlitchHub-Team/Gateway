package commandstests

import (
	"errors"
	"testing"
	"time"

	commands "Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensor "Gateway/internal/sensor"

	"github.com/google/uuid"
)

func TestAddSensorCmdExecute(t *testing.T) {
	// verifica che AddSensorCmd salvi il sensore e avvii lo starter
	cmdData := &commanddata.AddSensor{
		GatewayId: uuid.New(),
		SensorId:  uuid.New(),
		Interval:  time.Second,
	}
	adder := &mockSensorAdder{}
	starter := &mockSensorStarter{started: make(chan struct{}, 1)}

	cmd := commands.NewAddSensorCmd(cmdData, adder, starter, sensor.Active)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !adder.called {
		t.Fatal("expected AddSensor to be called")
	}

	if adder.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, adder.receivedCmd)
	}

	if adder.receivedStat != sensor.Active {
		t.Fatalf("expected status %q, got %q", sensor.Active, adder.receivedStat)
	}

	waitForSignal(t, starter.started, "sensor start")

	if cmd.String() != "AddSensorCmd" {
		t.Fatalf("expected String() to return AddSensorCmd, got %s", cmd.String())
	}
}

func TestAddSensorCmdExecuteReturnsAdderError(t *testing.T) {
	// verifica che AddSensorCmd propaghi l'errore del port di add
	expectedErr := errors.New("add failed")
	adder := &mockSensorAdder{err: expectedErr}
	starter := &mockSensorStarter{started: make(chan struct{}, 1)}

	cmd := commands.NewAddSensorCmd(&commanddata.AddSensor{}, adder, starter, sensor.Active)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	select {
	case <-starter.started:
		t.Fatal("expected Start not to be called on add error")
	default:
	}
}
