package sensortests

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	domainpkg "Gateway/internal/domain"
	sensor "Gateway/internal/sensor"
	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type mockProfile struct {
	mu            sync.Mutex
	generatedData *profiles.GeneratedSensorData
	generateCalls int
}

func (m *mockProfile) Generate() *profiles.GeneratedSensorData {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.generateCalls++
	return m.generatedData
}

func (m *mockProfile) String() string {
	return "mock-profile"
}

func (m *mockProfile) calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.generateCalls
}

type mockSaveSensorDataPort struct {
	mu            sync.Mutex
	err           error
	saveCalls     int
	savedData     *profiles.GeneratedSensorData
	savedGateway  uuid.UUID
	saveTriggered chan struct{}
}

func (m *mockSaveSensorDataPort) Save(data *profiles.GeneratedSensorData, gatewayID uuid.UUID) error {
	m.mu.Lock()
	m.saveCalls++
	m.savedData = data
	m.savedGateway = gatewayID
	trigger := m.saveTriggered
	err := m.err
	m.mu.Unlock()

	if trigger != nil {
		select {
		case trigger <- struct{}{}:
		default:
		}
	}

	return err
}

func (m *mockSaveSensorDataPort) calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.saveCalls
}

func (m *mockSaveSensorDataPort) lastSaved() (*profiles.GeneratedSensorData, uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.savedData, m.savedGateway
}

type mockCommand struct {
	executeErr error
	executed   chan struct{}
}

func (m *mockCommand) Execute() error {
	if m.executed != nil {
		select {
		case m.executed <- struct{}{}:
		default:
		}
	}

	return m.executeErr
}

func (m *mockCommand) String() string {
	return "mock-command"
}

type mockSerializableData struct{}

func (m *mockSerializableData) Serialize() ([]byte, error) {
	return []byte(`{"mock":true}`), nil
}

func TestStop(t *testing.T) {
	//* Verifica che Stop porti il sensore nello stato Stopped *
	sensorEntity := &sensor.Sensor{Status: sensor.Active}

	service := sensor.NewSensorService(
		sensorEntity,
		&mockSaveSensorDataPort{},
		make(chan domainpkg.BaseCommand, 1),
		make(chan error, 1),
		context.Background(),
		zap.NewNop(),
	)

	service.Stop()

	if sensorEntity.Status != sensor.Stopped {
		t.Fatalf("expected status %q, got %q", sensor.Stopped, sensorEntity.Status)
	}
}

func TestInterrupt(t *testing.T) {
	//* Verifica che Interrupt imposti lo stato Inactive *
	sensorEntity := &sensor.Sensor{Status: sensor.Active}

	service := sensor.NewSensorService(
		sensorEntity,
		&mockSaveSensorDataPort{},
		make(chan domainpkg.BaseCommand, 1),
		make(chan error, 1),
		context.Background(),
		zap.NewNop(),
	)

	service.Interrupt()

	if sensorEntity.Status != sensor.Inactive {
		t.Fatalf("expected status %q, got %q", sensor.Inactive, sensorEntity.Status)
	}
}

func TestResume(t *testing.T) {
	//* Verifica che Resume rimetta il sensore in stato Active *
	sensorEntity := &sensor.Sensor{Status: sensor.Inactive}

	service := sensor.NewSensorService(
		sensorEntity,
		&mockSaveSensorDataPort{},
		make(chan domainpkg.BaseCommand, 1),
		make(chan error, 1),
		context.Background(),
		zap.NewNop(),
	)

	service.Resume()

	if sensorEntity.Status != sensor.Active {
		t.Fatalf("expected status %q, got %q", sensor.Active, sensorEntity.Status)
	}
}

func TestStartExecutesCommandAndForwardsError(t *testing.T) {
	//* Verifica che Start esegua il comando e inoltri l'errore su errChannel *
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmdChannel := make(chan domainpkg.BaseCommand, 1)
	errChannel := make(chan error, 1)
	commandDone := make(chan struct{}, 1)
	serviceDone := make(chan struct{})
	expectedErr := errors.New("command failed")

	service := sensor.NewSensorService(
		newTestSensor(&mockProfile{generatedData: newGeneratedSensorData()}),
		&mockSaveSensorDataPort{},
		cmdChannel,
		errChannel,
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	cmdChannel <- &mockCommand{executeErr: expectedErr, executed: commandDone}

	waitForSignal(t, commandDone, "command execution")

	select {
	case err := <-errChannel:
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for forwarded command error")
	}

	cancel()
	waitForSignal(t, serviceDone, "service shutdown")
}

func TestStartGeneratesDataOnTick(t *testing.T) {
	//verifica che il ticker attivi il salvataggio quando il sensore e' active 
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	saveTriggered := make(chan struct{}, 1)
	serviceDone := make(chan struct{})
	generatedData := newGeneratedSensorData()
	profile := &mockProfile{generatedData: generatedData}
	savePort := &mockSaveSensorDataPort{saveTriggered: saveTriggered}
	sensorEntity := newTestSensor(profile)

	service := sensor.NewSensorService(
		sensorEntity,
		savePort,
		make(chan domainpkg.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	waitForSignal(t, saveTriggered, "sensor save")

	if profile.calls() == 0 {
		t.Fatal("expected Generate to be called at least once")
	}

	savedData, savedGatewayID := savePort.lastSaved()
	if savedData != generatedData {
		t.Fatalf("expected saved data %p, got %p", generatedData, savedData)
	}

	if savedGatewayID != sensorEntity.GatewayId {
		t.Fatalf("expected gateway id %v, got %v", sensorEntity.GatewayId, savedGatewayID)
	}

	cancel()
	waitForSignal(t, serviceDone, "service shutdown")
}

func TestStartSkipsGenerationWhenInactive(t *testing.T) {
	//erifica che il ticker non salvi dati quando il sensore e' inattivo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceDone := make(chan struct{})
	savePort := &mockSaveSensorDataPort{}

	sensorEntity := newTestSensor(&mockProfile{generatedData: newGeneratedSensorData()})
	sensorEntity.Status = sensor.Inactive

	service := sensor.NewSensorService(
		sensorEntity,
		savePort,
		make(chan domainpkg.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	time.Sleep(35 * time.Millisecond)

	if savePort.calls() != 0 {
		t.Fatalf("expected no save calls while inactive, got %d", savePort.calls())
	}

	cancel()
	waitForSignal(t, serviceDone, "service shutdown")
}

func TestStartStopsOnContextDone(t *testing.T) {
	//verifica che il servizio termini quando il context viene cancellato
	ctx, cancel := context.WithCancel(context.Background())
	serviceDone := make(chan struct{})

	service := sensor.NewSensorService(
		newTestSensor(&mockProfile{generatedData: newGeneratedSensorData()}),
		&mockSaveSensorDataPort{},
		make(chan domainpkg.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	cancel()
	waitForSignal(t, serviceDone, "service shutdown")
}

func newTestSensor(profile profiles.SensorProfile) *sensor.Sensor {
	return &sensor.Sensor{
		Id:        uuid.New(),
		GatewayId: uuid.New(),
		Profile:   profile,
		Interval:  10 * time.Millisecond,
		Status:    sensor.Active,
	}
}

func newGeneratedSensorData() *profiles.GeneratedSensorData {
	return &profiles.GeneratedSensorData{
		SensorId:  uuid.New(),
		Timestamp: time.Now(),
		Profile:   "mock-profile",
		Data:      &mockSerializableData{},
	}
}

func waitForSignal(t *testing.T, ch <-chan struct{}, label string) {
	//evita che i test concorrenti restino bloccati in attesa di un evento 
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", label)
	}
}
