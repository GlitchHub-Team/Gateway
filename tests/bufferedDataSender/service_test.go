package buffereddatasendertests

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"Gateway/internal/domain"
	buffereddatasender "Gateway/internal/bufferedDataSender"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestStartExecutesCommandAndForwardsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gateway := newGateway(domain.Active, time.Hour)
	js := &fakeJetStreamContext{}
	serviceDone := make(chan struct{})
	cmdCh := make(chan domain.BaseCommand, 1)
	errCh := make(chan error, 1)
	executed := make(chan struct{}, 1)
	expectedErr := errors.New("command failed")

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		buffereddatasender.NewNATSDataPublisherRepository(js),
		buffereddatasender.NewBufferedDataRepository(ctx, newBufferTestDB(t)),
		newFactory(js),
		cmdCh,
		errCh,
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	cmdCh <- &mockCommand{executeErr: expectedErr, executed: executed}
	waitForSignal(t, executed, "command execution")

	select {
	case err := <-errCh:
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for forwarded error")
	}

	cancel()
	waitForSignal(t, serviceDone, "service shutdown")
}

func TestStartPublishesAndCleansBufferedDataOnTickWhenActive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn := newBufferTestDB(t)
	gateway := newGateway(domain.Active, 10*time.Millisecond)
	tenantID := uuid.New()
	gateway.TenantId = &tenantID
	sensorID := uuid.New()
	js := &fakeJetStreamContext{publishSignal: make(chan struct{}, 1)}

	insertBufferRow(t, conn, gateway.Id, sensorID, time.Date(2026, time.March, 15, 11, 0, 0, 0, time.UTC), "HeartRate", `{"BpmValue":65}`)

	serviceDone := make(chan struct{})
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		buffereddatasender.NewNATSDataPublisherRepository(js),
		buffereddatasender.NewBufferedDataRepository(ctx, conn),
		newFactory(js),
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	waitForSignal(t, js.publishSignal, "publish")
	waitForBufferCount(t, conn, gateway.Id, 0)
	cancel()
	waitForSignal(t, serviceDone, "service shutdown")

	var count int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gateway.Id.String()).Scan(&count); err != nil {
		t.Fatalf("expected count query to succeed, got %v", err)
	}

	var dto buffereddatasender.SensorDataDTO
	if err := json.Unmarshal(js.publishedData, &dto); err != nil {
		t.Fatalf("expected valid published json, got %v", err)
	}

	if dto.GatewayId != gateway.Id || dto.SensorId != sensorID || dto.TenantId != tenantID {
		t.Fatalf("unexpected published dto: %+v", dto)
	}
}

func TestStartSkipsPublishingWhenGatewayIsInactive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn := newBufferTestDB(t)
	gateway := newGateway(domain.Inactive, 10*time.Millisecond)
	insertBufferRow(t, conn, gateway.Id, uuid.New(), time.Now().UTC(), "HeartRate", `{"BpmValue":65}`)
	js := &fakeJetStreamContext{}
	serviceDone := make(chan struct{})

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		buffereddatasender.NewNATSDataPublisherRepository(js),
		buffereddatasender.NewBufferedDataRepository(ctx, conn),
		newFactory(js),
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	go func() {
		defer close(serviceDone)
		service.Start()
	}()

	time.Sleep(40 * time.Millisecond)
	cancel()
	waitForSignal(t, serviceDone, "service shutdown")

	if js.publishCalls != 0 {
		t.Fatalf("expected no publish while inactive, got %d", js.publishCalls)
	}

	var count int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gateway.Id.String()).Scan(&count); err != nil {
		t.Fatalf("expected count query to succeed, got %v", err)
	}

	if count != 1 {
		t.Fatalf("expected buffered row to remain, got %d", count)
	}
}

func TestHelloPublishesUsingSendPort(t *testing.T) {
	js := &fakeJetStreamContext{}
	gateway := newGateway(domain.Active, time.Second)
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		buffereddatasender.NewNATSDataPublisherRepository(js),
		buffereddatasender.NewBufferedDataRepository(context.Background(), newBufferTestDB(t)),
		newFactory(js),
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		context.Background(),
		zap.NewNop(),
	)

	if err := service.Hello(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedSubject := "gateway.hello." + gateway.Id.String()
	if js.publishedSubj != expectedSubject {
		t.Fatalf("expected subject %q, got %q", expectedSubject, js.publishedSubj)
	}
}

func TestDecommissionCleansBufferAndSwitchesSender(t *testing.T) {
	ctx := context.Background()
	conn := newBufferTestDB(t)
	gateway := newGateway(domain.Active, time.Second)
	tenantID := uuid.New()
	token := "jwt-token"
	gateway.TenantId = &tenantID
	gateway.Token = &token
	insertBufferRow(t, conn, gateway.Id, uuid.New(), time.Now().UTC(), "HeartRate", `{"BpmValue":65}`)

	initialJS := &fakeJetStreamContext{}
	replacementJS := &fakeJetStreamContext{}
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		buffereddatasender.NewNATSDataPublisherRepository(initialJS),
		buffereddatasender.NewBufferedDataRepository(ctx, conn),
		newFactory(replacementJS),
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		ctx,
		zap.NewNop(),
	)

	if err := service.Decommission(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if gateway.Status != domain.Decommissioned || gateway.TenantId != nil || gateway.Token != nil {
		t.Fatalf("expected gateway decommissioned, got %+v", gateway)
	}

	var count int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gateway.Id.String()).Scan(&count); err != nil {
		t.Fatalf("expected count query to succeed, got %v", err)
	}

	if count != 0 {
		t.Fatalf("expected empty buffer after decommission, got %d rows", count)
	}

	if err := service.Hello(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if replacementJS.publishCalls != 1 || initialJS.publishCalls != 0 {
		t.Fatalf("expected Hello to use replacement sender, got replacement=%d initial=%d", replacementJS.publishCalls, initialJS.publishCalls)
	}
}

func TestDecommissionReturnsErrorWhenBufferCleanupFails(t *testing.T) {
	service := buffereddatasender.NewBufferedDataSenderService(
		newGateway(domain.Active, time.Second),
		buffereddatasender.NewNATSDataPublisherRepository(&fakeJetStreamContext{}),
		buffereddatasender.NewBufferedDataRepository(context.Background(), newNonBufferDB(t)),
		newFactory(&fakeJetStreamContext{}),
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		context.Background(),
		zap.NewNop(),
	)

	err := service.Decommission()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCommissionReturnsErrorWhenReloadFails(t *testing.T) {
	gateway := newGateway(domain.Decommissioned, time.Second)
	gateway.SecretKey = validSeed(t)
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		buffereddatasender.NewNATSDataPublisherRepository(&fakeJetStreamContext{}),
		buffereddatasender.NewBufferedDataRepository(context.Background(), newBufferTestDB(t)),
		buffereddatasender.NewNATSDataPublisherFactory(nil, "127.0.0.1", 1),
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		context.Background(),
		zap.NewNop(),
	)

	err := service.Commission(uuid.New(), "jwt-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if gateway.Status != domain.Decommissioned || gateway.TenantId != nil || gateway.Token != nil {
		t.Fatalf("expected gateway state unchanged, got %+v", gateway)
	}
}

func TestResetUpdatesIntervalAndCleansBuffer(t *testing.T) {
	conn := newBufferTestDB(t)
	gateway := newGateway(domain.Active, time.Hour)
	insertBufferRow(t, conn, gateway.Id, uuid.New(), time.Now().UTC(), "HeartRate", `{"BpmValue":65}`)

	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		buffereddatasender.NewNATSDataPublisherRepository(&fakeJetStreamContext{}),
		buffereddatasender.NewBufferedDataRepository(context.Background(), conn),
		newFactory(&fakeJetStreamContext{}),
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		context.Background(),
		zap.NewNop(),
	)

	if err := service.Reset(5 * time.Second); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if gateway.Interval != 5*time.Second {
		t.Fatalf("expected interval 5s, got %v", gateway.Interval)
	}

	var count int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gateway.Id.String()).Scan(&count); err != nil {
		t.Fatalf("expected count query to succeed, got %v", err)
	}

	if count != 0 {
		t.Fatalf("expected cleaned buffer, got %d rows", count)
	}
}

func TestStopInterruptAndResumeUpdateGatewayStatus(t *testing.T) {
	gateway := newGateway(domain.Active, time.Second)
	service := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		buffereddatasender.NewNATSDataPublisherRepository(&fakeJetStreamContext{}),
		buffereddatasender.NewBufferedDataRepository(context.Background(), newBufferTestDB(t)),
		newFactory(&fakeJetStreamContext{}),
		make(chan domain.BaseCommand, 1),
		make(chan error, 1),
		context.Background(),
		zap.NewNop(),
	)

	service.Interrupt()
	if gateway.Status != domain.Inactive {
		t.Fatalf("expected inactive status, got %q", gateway.Status)
	}

	service.Resume()
	if gateway.Status != domain.Active {
		t.Fatalf("expected active status, got %q", gateway.Status)
	}

	service.Stop()
	if gateway.Status != domain.Stopped {
		t.Fatalf("expected stopped status, got %q", gateway.Status)
	}
}

func waitForBufferCount(t *testing.T, conn interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}, gatewayID uuid.UUID, want int) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		var count int
		if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gatewayID.String()).Scan(&count); err != nil {
			t.Fatalf("expected count query to succeed, got %v", err)
		}
		if count == want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for buffer count %d", want)
}
