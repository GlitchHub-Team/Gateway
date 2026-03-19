package buffereddatasendertests

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	buffereddatasender "Gateway/internal/bufferedDataSender"

	"github.com/google/uuid"
)

func TestSendPublishesExpectedSubjectAndPayload(t *testing.T) {
	js := &fakeJetStreamContext{}
	repo := buffereddatasender.NewNATSDataPublisherRepository(js)
	gatewayID := uuid.New()
	sensorID := uuid.New()
	timestamp := time.Date(2026, time.March, 15, 10, 0, 0, 0, time.UTC)
	tenantID := uuid.New()

	conn := newBufferTestDB(t)
	insertBufferRow(t, conn, gatewayID, sensorID, timestamp, "HeartRate", `{"BpmValue":75}`)
	loaded, err := buffereddatasender.NewBufferedDataRepository(context.Background(), conn).GetOrderedBufferedData(gatewayID)
	if err != nil {
		t.Fatalf("expected nil error loading data, got %v", err)
	}

	if err := repo.Send(loaded[0], tenantID); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedSubject := "sensor." + gatewayID.String() + "." + sensorID.String()
	if js.publishedSubj != expectedSubject {
		t.Fatalf("expected subject %q, got %q", expectedSubject, js.publishedSubj)
	}

	var dto buffereddatasender.SensorDataDTO
	if err := json.Unmarshal(js.publishedData, &dto); err != nil {
		t.Fatalf("expected valid published json, got %v", err)
	}

	if dto.GatewayId != gatewayID || dto.SensorId != sensorID || dto.TenantId != tenantID || dto.Profile != "HeartRate" {
		t.Fatalf("unexpected dto published: %+v", dto)
	}
}

func TestHelloPublishesExpectedPayload(t *testing.T) {
	js := &fakeJetStreamContext{}
	repo := buffereddatasender.NewNATSDataPublisherRepository(js)
	gatewayID := uuid.New()

	if err := repo.Hello(gatewayID, "public-key"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedSubject := "gateway.hello." + gatewayID.String()
	if js.publishedSubj != expectedSubject {
		t.Fatalf("expected subject %q, got %q", expectedSubject, js.publishedSubj)
	}

	var dto buffereddatasender.HelloMessageDTO
	if err := json.Unmarshal(js.publishedData, &dto); err != nil {
		t.Fatalf("expected valid hello payload, got %v", err)
	}

	if dto.GatewayId != gatewayID || dto.PublicIdentifier != "public-key" {
		t.Fatalf("unexpected hello dto: %+v", dto)
	}
}

func TestHelloReturnsWrappedPublishError(t *testing.T) {
	expectedErr := errors.New("publish failed")
	js := &fakeJetStreamContext{publishErr: expectedErr}
	repo := buffereddatasender.NewNATSDataPublisherRepository(js)
	gatewayID := uuid.New()

	err := repo.Hello(gatewayID, "public-key")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}

	if !strings.Contains(err.Error(), gatewayID.String()) {
		t.Fatalf("expected gateway id %s in error, got %q", gatewayID, err.Error())
	}
}
