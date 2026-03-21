package buffereddatasender_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"
	"Gateway/internal/natsutil"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

type mockCommand struct {
	err      error
	executed chan struct{}
}

func (m *mockCommand) Execute() error {
	if m.executed != nil {
		select {
		case m.executed <- struct{}{}:
		default:
		}
	}
	return m.err
}

func (m *mockCommand) String() string { return "mock-command" }

type mockSendSensorDataPortFactory struct {
	createFn     func() buffereddatasender.SendSensorDataPort
	reloadFn     func(token string, seed string) (buffereddatasender.SendSensorDataPort, error)
	createCalls  int
	reloadCalls  int
	lastToken    string
	lastSeed     string
	reloadErr    error
	reloadedPort buffereddatasender.SendSensorDataPort
	createdPort  buffereddatasender.SendSensorDataPort
}

func (m *mockSendSensorDataPortFactory) Create() buffereddatasender.SendSensorDataPort {
	m.createCalls++
	if m.createFn != nil {
		m.createdPort = m.createFn()
		return m.createdPort
	}
	return m.createdPort
}

func (m *mockSendSensorDataPortFactory) Reload(token string, seed string) (buffereddatasender.SendSensorDataPort, error) {
	m.reloadCalls++
	m.lastToken = token
	m.lastSeed = seed
	if m.reloadFn != nil {
		return m.reloadFn(token, seed)
	}
	if m.reloadErr != nil {
		return nil, m.reloadErr
	}
	return m.reloadedPort, nil
}

func newGateway(status domain.GatewayStatus, interval time.Duration) *configmanager.Gateway {
	tenant := uuid.New()
	token := "gateway-token"
	return &configmanager.Gateway{
		Id:               uuid.New(),
		TenantId:         &tenant,
		Status:           status,
		Interval:         interval,
		PublicIdentifier: "pub-id",
		SecretKey:        "secret-seed",
		Token:            &token,
	}
}

func createHelloStream(t *testing.T, js jetstream.JetStream) {
	t.Helper()

	_, err := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:     "HELLO_STREAM",
		Subjects: []string{"gateway.hello.*"},
	})
	if err != nil {
		t.Fatalf("unable to create hello stream: %v", err)
	}
}

func waitForSignal(t *testing.T, ch <-chan struct{}, label string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for %s", label)
	}
}

func waitForCondition(t *testing.T, label string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("condition not met: %s", label)
}

func insertBufferedRow(t *testing.T, conn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}, gatewayID uuid.UUID, sensorID uuid.UUID, ts time.Time, profile string, value string,
) {
	t.Helper()
	_, err := conn.ExecContext(context.Background(), `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`, gatewayID.String(), sensorID.String(), ts.UTC(), profile, value)
	if err != nil {
		t.Fatalf("insert setup failed: %v", err)
	}
}

func TestBufferedDataSenderStartExecutesCommandAndForwardsNilError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gateway := newGateway(domain.Inactive, time.Hour)
	nc := newMockNATSConnection(t)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())

	bufferRepo, _ := newMockBufferRepository(t)
	cmdChannel := make(chan domain.BaseCommand, 1)
	errChannel := make(chan error, 1)
	serviceDone := make(chan struct{})
	executed := make(chan struct{}, 1)

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		bufferRepo,
		&mockSendSensorDataPortFactory{},
		cmdChannel,
		errChannel,
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	cmdChannel <- &mockCommand{executed: executed}
	waitForSignal(t, executed, "command execution")

	select {
	case err := <-errChannel:
		if err != nil {
			t.Fatalf("expected nil command error, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for command error forwarding")
	}

	cancel()
	waitForSignal(t, serviceDone, "service stop")
}

func TestBufferedDataSenderStartExecutesCommandAndForwardsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expectedErr := errors.New("command failed")
	gateway := newGateway(domain.Inactive, time.Hour)
	nc := newMockNATSConnection(t)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())

	bufferRepo, _ := newMockBufferRepository(t)
	cmdChannel := make(chan domain.BaseCommand, 1)
	errChannel := make(chan error, 1)
	serviceDone := make(chan struct{})
	executed := make(chan struct{}, 1)

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		bufferRepo,
		&mockSendSensorDataPortFactory{},
		cmdChannel,
		errChannel,
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	cmdChannel <- &mockCommand{err: expectedErr, executed: executed}
	waitForSignal(t, executed, "command execution")

	select {
	case err := <-errChannel:
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for command error forwarding")
	}

	cancel()
	waitForSignal(t, serviceDone, "service stop")
}

func TestBufferedDataSenderStartTickerSendsDataAndCleansBuffer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bufferRepo, conn := newMockBufferRepository(t)
	gateway := newGateway(domain.Active, 15*time.Millisecond)

	sensorID := uuid.New()
	insertBufferedRow(t, conn, gateway.Id, sensorID, time.Now().UTC(), "HeartRate", `{"BpmValue":70}`)

	nc := newMockNATSConnection(t)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())

	subject := fmt.Sprintf("sensor.%s.%s", gateway.Id, sensorID)
	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		t.Fatalf("unable to subscribe on sensor subject: %v", err)
	}

	serviceDone := make(chan struct{})
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		bufferRepo,
		&mockSendSensorDataPortFactory{},
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	if _, err := sub.NextMsg(2 * time.Second); err != nil {
		t.Fatalf("expected published sensor message, got: %v", err)
	}

	waitForCondition(t, "buffer should be cleaned after successful send", func() bool {
		rows, err := bufferRepo.GetOrderedBufferedData(gateway.Id)
		return err == nil && len(rows) == 0
	})

	cancel()
	waitForSignal(t, serviceDone, "service stop")
}

func TestBufferedDataSenderStartTickerGetOrderedBufferedDataError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	brokenRepo := buffereddatasender.NewBufferedDataRepository(context.Background(), struct{ *sql.DB }{db})
	gateway := newGateway(domain.Active, 15*time.Millisecond)

	nc := newMockNATSConnection(t)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())

	sub, err := nc.SubscribeSync("sensor.>")
	if err != nil {
		t.Fatalf("unable to subscribe to wildcard sensor topic: %v", err)
	}

	serviceDone := make(chan struct{})
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		brokenRepo,
		&mockSendSensorDataPortFactory{},
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	if _, err := sub.NextMsg(180 * time.Millisecond); !errors.Is(err, nats.ErrTimeout) {
		t.Fatalf("expected no message due to getOrderedBufferedData error, got: %v", err)
	}

	cancel()
	waitForSignal(t, serviceDone, "service stop")
}

func TestBufferedDataSenderStartKeepsUnsentDataInBuffer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bufferRepo, conn := newMockBufferRepository(t)
	gateway := newGateway(domain.Active, 15*time.Millisecond)

	firstSensorID := uuid.New()
	secondSensorID := uuid.New()
	t0 := time.Now().UTC()
	insertBufferedRow(t, conn, gateway.Id, firstSensorID, t0, "HeartRate", `{"BpmValue":60}`)

	bigPayload := fmt.Sprintf(`{"blob":"%s"}`, strings.Repeat("a", 2*1024*1024))
	insertBufferedRow(t, conn, gateway.Id, secondSensorID, t0.Add(time.Millisecond), "Large", bigPayload)

	nc := newMockNATSConnection(t)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())

	firstSubject := fmt.Sprintf("sensor.%s.%s", gateway.Id, firstSensorID)
	secondSubject := fmt.Sprintf("sensor.%s.%s", gateway.Id, secondSensorID)
	firstSub, err := nc.SubscribeSync(firstSubject)
	if err != nil {
		t.Fatalf("unable to subscribe on first subject: %v", err)
	}
	secondSub, err := nc.SubscribeSync(secondSubject)
	if err != nil {
		t.Fatalf("unable to subscribe on second subject: %v", err)
	}

	serviceDone := make(chan struct{})
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		bufferRepo,
		&mockSendSensorDataPortFactory{},
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	if _, err := firstSub.NextMsg(2 * time.Second); err != nil {
		t.Fatalf("expected first message to be published, got: %v", err)
	}

	if _, err := secondSub.NextMsg(200 * time.Millisecond); !errors.Is(err, nats.ErrTimeout) {
		t.Fatalf("expected second message publish failure, got: %v", err)
	}

	waitForCondition(t, "unsent data should remain buffered", func() bool {
		rows, err := bufferRepo.GetOrderedBufferedData(gateway.Id)
		if err != nil {
			return false
		}
		return len(rows) == 1 && rows[0].SensorId == secondSensorID
	})

	cancel()
	waitForSignal(t, serviceDone, "service stop")
}

func TestBufferedDataSenderHelloSuccess(t *testing.T) {
	ctx := context.Background()
	bufferRepo, _ := newMockBufferRepository(t)
	gateway := newGateway(domain.Active, time.Second)

	nc := newMockNATSConnection(t)
	js, err := natsutil.NewJetStreamContext(nc)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}
	createHelloStream(t, js)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, js, ctx)

	subject := fmt.Sprintf("gateway.hello.%s", gateway.Id)
	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		t.Fatalf("unable to subscribe to hello topic: %v", err)
	}

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		bufferRepo,
		&mockSendSensorDataPortFactory{},
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	if err := service.Hello(); err != nil {
		t.Fatalf("expected hello to succeed, got: %v", err)
	}

	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("expected hello message, got: %v", err)
	}

	var dto buffereddatasender.HelloMessageDTO
	if err := json.Unmarshal(msg.Data, &dto); err != nil {
		t.Fatalf("unable to unmarshal hello dto: %v", err)
	}
	if dto.GatewayId != gateway.Id {
		t.Fatalf("unexpected gateway id: got %s want %s", dto.GatewayId, gateway.Id)
	}
	if dto.PublicIdentifier != gateway.PublicIdentifier {
		t.Fatalf("unexpected public identifier: got %s want %s", dto.PublicIdentifier, gateway.PublicIdentifier)
	}
}

func TestBufferedDataSenderHelloError(t *testing.T) {
	ctx := context.Background()
	bufferRepo, _ := newMockBufferRepository(t)
	gateway := newGateway(domain.Active, time.Second)

	nc := newMockNATSConnection(t)
	js, err := natsutil.NewJetStreamContext(nc)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}
	nc.Close()

	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, js, ctx)
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		bufferRepo,
		&mockSendSensorDataPortFactory{},
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	if err := service.Hello(); err == nil {
		t.Fatal("expected hello to fail with invalid jetstream context")
	}
}

func TestBufferedDataSenderDecommissionSuccess(t *testing.T) {
	ctx := context.Background()
	bufferRepo, conn := newMockBufferRepository(t)
	gateway := newGateway(domain.Active, time.Second)

	insertBufferedRow(t, conn, gateway.Id, uuid.New(), time.Now().UTC(), "HeartRate", `{"BpmValue":80}`)

	invalidNC := newMockNATSConnection(t)
	invalidJS, err := natsutil.NewJetStreamContext(invalidNC)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}
	invalidNC.Close()
	initialPort := buffereddatasender.NewNATSDataPublisherRepository(invalidNC, invalidJS, ctx)

	validNC := newMockNATSConnection(t)
	validJS, err := natsutil.NewJetStreamContext(validNC)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}
	createHelloStream(t, validJS)
	replacementPort := buffereddatasender.NewNATSDataPublisherRepository(validNC, validJS, ctx)

	factory := &mockSendSensorDataPortFactory{createFn: func() buffereddatasender.SendSensorDataPort { return replacementPort }}

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		initialPort,
		bufferRepo,
		factory,
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	if err := service.Hello(); err == nil {
		t.Fatal("expected initial hello failure before decommission")
	}

	if err := service.Decommission(); err != nil {
		t.Fatalf("expected decommission success, got: %v", err)
	}

	if factory.createCalls != 1 {
		t.Fatalf("expected factory Create to be called once, got %d", factory.createCalls)
	}

	rows, err := bufferRepo.GetOrderedBufferedData(gateway.Id)
	if err != nil {
		t.Fatalf("failed to read buffer after decommission: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected empty buffer after decommission, got %d rows", len(rows))
	}

	if gateway.Status != domain.Decommissioned {
		t.Fatalf("expected gateway status %q, got %q", domain.Decommissioned, gateway.Status)
	}
	if gateway.TenantId != nil {
		t.Fatal("expected tenant id to be nil after decommission")
	}
	if gateway.Token != nil {
		t.Fatal("expected token to be nil after decommission")
	}

	if err := service.Hello(); err != nil {
		t.Fatalf("expected hello to use replacement sendDataPort, got: %v", err)
	}
}

func TestBufferedDataSenderDecommissionError(t *testing.T) {
	ctx := context.Background()
	gateway := newGateway(domain.Active, time.Second)
	beforeTenant := *gateway.TenantId
	beforeToken := *gateway.Token

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	brokenRepo := buffereddatasender.NewBufferedDataRepository(context.Background(), struct{ *sql.DB }{db})
	nc := newMockNATSConnection(t)
	initialPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, ctx)
	factory := &mockSendSensorDataPortFactory{}

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		initialPort,
		brokenRepo,
		factory,
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	if err := service.Decommission(); err == nil {
		t.Fatal("expected decommission error when buffer clean fails")
	}

	if factory.createCalls != 0 {
		t.Fatalf("expected factory Create not to be called, got %d", factory.createCalls)
	}

	if gateway.Status != domain.Active {
		t.Fatalf("expected gateway status unchanged, got %q", gateway.Status)
	}
	if gateway.TenantId == nil || *gateway.TenantId != beforeTenant {
		t.Fatal("expected tenant id unchanged on decommission error")
	}
	if gateway.Token == nil || *gateway.Token != beforeToken {
		t.Fatal("expected token unchanged on decommission error")
	}
}

func TestBufferedDataSenderCommissionSuccess(t *testing.T) {
	ctx := context.Background()
	bufferRepo, _ := newMockBufferRepository(t)
	gateway := newGateway(domain.Inactive, time.Second)

	invalidNC := newMockNATSConnection(t)
	invalidJS, err := natsutil.NewJetStreamContext(invalidNC)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}
	invalidNC.Close()
	initialPort := buffereddatasender.NewNATSDataPublisherRepository(invalidNC, invalidJS, ctx)

	validNC := newMockNATSConnection(t)
	validJS, err := natsutil.NewJetStreamContext(validNC)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}
	createHelloStream(t, validJS)
	reloadedPort := buffereddatasender.NewNATSDataPublisherRepository(validNC, validJS, ctx)
	factory := &mockSendSensorDataPortFactory{reloadFn: func(token string, seed string) (buffereddatasender.SendSensorDataPort, error) {
		return reloadedPort, nil
	}}

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		initialPort,
		bufferRepo,
		factory,
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	if err := service.Hello(); err == nil {
		t.Fatal("expected initial hello failure before commission")
	}

	newTenant := uuid.New()
	newToken := "commissioned-token"
	if err := service.Commission(newTenant, newToken); err != nil {
		t.Fatalf("expected commission success, got: %v", err)
	}

	if factory.reloadCalls != 1 {
		t.Fatalf("expected one reload call, got %d", factory.reloadCalls)
	}
	if factory.lastToken != newToken {
		t.Fatalf("expected reload token %q, got %q", newToken, factory.lastToken)
	}
	if factory.lastSeed != gateway.SecretKey {
		t.Fatalf("expected reload seed %q, got %q", gateway.SecretKey, factory.lastSeed)
	}

	if gateway.Status != domain.Active {
		t.Fatalf("expected gateway status %q, got %q", domain.Active, gateway.Status)
	}
	if gateway.TenantId == nil || *gateway.TenantId != newTenant {
		t.Fatal("expected tenant id to be updated after commission")
	}
	if gateway.Token == nil || *gateway.Token != newToken {
		t.Fatal("expected token to be updated after commission")
	}

	if err := service.Hello(); err != nil {
		t.Fatalf("expected hello to use reloaded sendDataPort, got: %v", err)
	}
}

func TestBufferedDataSenderCommissionError(t *testing.T) {
	ctx := context.Background()
	bufferRepo, _ := newMockBufferRepository(t)
	gateway := newGateway(domain.Inactive, time.Second)
	originalTenant := *gateway.TenantId
	originalToken := *gateway.Token

	nc := newMockNATSConnection(t)
	initialPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, ctx)
	expectedErr := errors.New("reload failed")
	factory := &mockSendSensorDataPortFactory{reloadErr: expectedErr}

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		initialPort,
		bufferRepo,
		factory,
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	err := service.Commission(uuid.New(), "new-token")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected commission error %v, got %v", expectedErr, err)
	}

	if gateway.Status != domain.Inactive {
		t.Fatalf("expected gateway status unchanged, got %q", gateway.Status)
	}
	if gateway.TenantId == nil || *gateway.TenantId != originalTenant {
		t.Fatal("expected tenant id unchanged on commission error")
	}
	if gateway.Token == nil || *gateway.Token != originalToken {
		t.Fatal("expected token unchanged on commission error")
	}
}

func TestBufferedDataSenderResetSuccess(t *testing.T) {
	ctx := context.Background()
	bufferRepo, conn := newMockBufferRepository(t)
	gateway := newGateway(domain.Active, 100*time.Millisecond)

	insertBufferedRow(t, conn, gateway.Id, uuid.New(), time.Now().UTC(), "HeartRate", `{"BpmValue":90}`)

	nc := newMockNATSConnection(t)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, ctx)
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		bufferRepo,
		&mockSendSensorDataPortFactory{},
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	newInterval := 25 * time.Millisecond
	if err := service.Reset(newInterval); err != nil {
		t.Fatalf("expected reset success, got: %v", err)
	}

	if gateway.Interval != newInterval {
		t.Fatalf("expected interval %s, got %s", newInterval, gateway.Interval)
	}

	rows, err := bufferRepo.GetOrderedBufferedData(gateway.Id)
	if err != nil {
		t.Fatalf("failed to read buffer after reset: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected empty buffer after reset, got %d rows", len(rows))
	}
}

func TestBufferedDataSenderResetError(t *testing.T) {
	ctx := context.Background()
	gateway := newGateway(domain.Active, 100*time.Millisecond)

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	brokenRepo := buffereddatasender.NewBufferedDataRepository(context.Background(), struct{ *sql.DB }{db})
	nc := newMockNATSConnection(t)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, ctx)
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		brokenRepo,
		&mockSendSensorDataPortFactory{},
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	newInterval := 40 * time.Millisecond
	err = service.Reset(newInterval)
	if err == nil {
		t.Fatal("expected reset error when clean buffer fails")
	}

	if gateway.Interval != newInterval {
		t.Fatalf("expected interval to be updated before failure, got %s", gateway.Interval)
	}
}

func TestBufferedDataSenderStopInterruptResume(t *testing.T) {
	ctx := context.Background()
	bufferRepo, _ := newMockBufferRepository(t)
	gateway := newGateway(domain.Active, time.Second)
	nc := newMockNATSConnection(t)
	sendPort := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, ctx)

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		sendPort,
		bufferRepo,
		&mockSendSensorDataPortFactory{},
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	service.Interrupt()
	if gateway.Status != domain.Inactive {
		t.Fatalf("expected status %q after interrupt, got %q", domain.Inactive, gateway.Status)
	}

	service.Resume()
	if gateway.Status != domain.Active {
		t.Fatalf("expected status %q after resume, got %q", domain.Active, gateway.Status)
	}

	service.Stop()
	if gateway.Status != domain.Stopped {
		t.Fatalf("expected status %q after stop, got %q", domain.Stopped, gateway.Status)
	}
}
