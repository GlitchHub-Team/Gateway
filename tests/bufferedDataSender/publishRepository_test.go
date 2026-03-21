package buffereddatasender_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	natsserver "Gateway/cmd/external/natsServer"
	buffereddatasender "Gateway/internal/bufferedDataSender"
	"Gateway/internal/natsutil"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func newMockNATSConnection(t *testing.T) *nats.Conn {
	t.Helper()

	root := moduleRoot(t)
	token, seed := parseNATSCreds(t, filepath.Join(root, "cmd", os.Getenv("GATEWAY_BASE_CREDS_PATH")))

	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	nc := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = nc.Drain() })

	return nc
}

func callSend(
	t *testing.T,
	repo *buffereddatasender.NATSDataPublisherRepository,
	sensorID uuid.UUID,
	gatewayID uuid.UUID,
	timestamp time.Time,
	profile string,
	data []byte,
	tenantID uuid.UUID,
) error {
	t.Helper()

	sendMethod := reflect.ValueOf(repo).MethodByName("Send")
	if !sendMethod.IsValid() {
		t.Fatal("Send method not found")
	}

	sendDataType := sendMethod.Type().In(0)
	sendDataValue := reflect.New(sendDataType.Elem())
	elem := sendDataValue.Elem()

	elem.FieldByName("SensorId").Set(reflect.ValueOf(sensorID))
	elem.FieldByName("GatewayId").Set(reflect.ValueOf(gatewayID))
	elem.FieldByName("Timestamp").Set(reflect.ValueOf(timestamp))
	elem.FieldByName("Profile").Set(reflect.ValueOf(profile))
	elem.FieldByName("Data").Set(reflect.ValueOf(data))

	returns := sendMethod.Call([]reflect.Value{sendDataValue, reflect.ValueOf(tenantID)})
	if returns[0].IsNil() {
		return nil
	}

	return returns[0].Interface().(error)
}

func TestPublishRepositorySendValidDataPublishesExpectedPayload(t *testing.T) {
	nc := newMockNATSConnection(t)

	sensorID := uuid.New()
	gatewayID := uuid.New()
	tenantID := uuid.New()
	ts := time.Now()

	subject := "sensor." + gatewayID.String() + "." + sensorID.String()
	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		t.Fatalf("unable to subscribe to test subject: %v", err)
	}

	repo := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())
	if err := callSend(t, repo, sensorID, gatewayID, ts, "HeartRate", []byte(`{"BpmValue":72}`), tenantID); err != nil {
		t.Fatalf("send returned unexpected error: %v", err)
	}

	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("expected published message, got: %v", err)
	}

	var got buffereddatasender.SensorDataDTO
	if err := json.Unmarshal(msg.Data, &got); err != nil {
		t.Fatalf("unable to unmarshal published dto: %v", err)
	}

	if got.SensorId != sensorID {
		t.Fatalf("unexpected sensor id: got %s want %s", got.SensorId, sensorID)
	}
	if got.GatewayId != gatewayID {
		t.Fatalf("unexpected gateway id: got %s want %s", got.GatewayId, gatewayID)
	}
	if got.TenantId != tenantID {
		t.Fatalf("unexpected tenant id: got %s want %s", got.TenantId, tenantID)
	}
	if !got.Timestamp.Equal(ts) {
		t.Fatalf("unexpected timestamp: got %s want %s", got.Timestamp, ts)
	}
	if got.Profile != "HeartRate" {
		t.Fatalf("unexpected profile: got %s", got.Profile)
	}

	var gotData map[string]any
	if err := json.Unmarshal(got.Data, &gotData); err != nil {
		t.Fatalf("unable to unmarshal data field: %v", err)
	}
	if bpm, ok := gotData["BpmValue"].(float64); !ok || bpm != 72 {
		t.Fatalf("unexpected data payload: %+v", gotData)
	}
}

func TestPublishRepositorySendInvalidDomainDataReturnsMarshalError(t *testing.T) {
	nc := newMockNATSConnection(t)
	repo := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())

	err := callSend(
		t,
		repo,
		uuid.New(),
		uuid.New(),
		time.Now(),
		"HeartRate",
		[]byte("{"),
		uuid.New(),
	)

	if err == nil {
		t.Fatal("expected marshal error, got nil")
	}
}

func TestPublishRepositorySendLogicalInvalidDataStillPublishes(t *testing.T) {
	nc := newMockNATSConnection(t)

	sensorID := uuid.Nil
	gatewayID := uuid.Nil
	tenantID := uuid.Nil
	ts := time.Time{}

	subject := "sensor." + gatewayID.String() + "." + sensorID.String()
	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		t.Fatalf("unable to subscribe to test subject: %v", err)
	}

	repo := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())
	if err := callSend(t, repo, sensorID, gatewayID, ts, "", []byte(`{}`), tenantID); err != nil {
		t.Fatalf("expected no error for structurally valid data, got: %v", err)
	}

	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("expected published message, got: %v", err)
	}

	var got buffereddatasender.SensorDataDTO
	if err := json.Unmarshal(msg.Data, &got); err != nil {
		t.Fatalf("unable to unmarshal published dto: %v", err)
	}
	if got.SensorId != uuid.Nil || got.GatewayId != uuid.Nil || got.TenantId != uuid.Nil {
		t.Fatalf("expected nil ids in payload, got sensor=%s gateway=%s tenant=%s", got.SensorId, got.GatewayId, got.TenantId)
	}
}

func TestPublishRepositorySendWithClosedNATSConnectionReturnsError(t *testing.T) {
	nc := newMockNATSConnection(t)
	repo := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())

	nc.Close()

	err := callSend(
		t,
		repo,
		uuid.New(),
		uuid.New(),
		time.Now().UTC(),
		"HeartRate",
		[]byte(`{"BpmValue":65}`),
		uuid.New(),
	)

	if err == nil {
		t.Fatal("expected publish error with closed nats connection")
	}
}

func TestPublishRepositorySendNilDataPanics(t *testing.T) {
	nc := newMockNATSConnection(t)
	repo := buffereddatasender.NewNATSDataPublisherRepository(nc, nil, context.Background())
	sendMethod := reflect.ValueOf(repo).MethodByName("Send")
	sendDataType := sendMethod.Type().In(0)

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected panic when sensor data is nil")
		}
	}()

	_ = sendMethod.Call([]reflect.Value{reflect.Zero(sendDataType), reflect.ValueOf(uuid.New())})
}

func TestPublishRepositoryHelloPublishesExpectedPayload(t *testing.T) {
	nc := newMockNATSConnection(t)
	ctx := context.Background()

	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}

	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "HELLO_STREAM",
		Subjects: []string{"gateway.hello.*"},
	})
	if err != nil {
		t.Fatalf("unable to create hello stream: %v", err)
	}

	gatewayID := uuid.New()
	publicID := "gateway-public-id"
	subject := "gateway.hello." + gatewayID.String()

	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		t.Fatalf("unable to subscribe to hello subject: %v", err)
	}

	repo := buffereddatasender.NewNATSDataPublisherRepository(nc, js, ctx)
	if err := repo.Hello(gatewayID, publicID); err != nil {
		t.Fatalf("hello returned unexpected error: %v", err)
	}

	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("expected hello message, got: %v", err)
	}

	var got buffereddatasender.HelloMessageDTO
	if err := json.Unmarshal(msg.Data, &got); err != nil {
		t.Fatalf("unable to unmarshal hello dto: %v", err)
	}

	if got.GatewayId != gatewayID {
		t.Fatalf("unexpected gateway id: got %s want %s", got.GatewayId, gatewayID)
	}
	if got.PublicIdentifier != publicID {
		t.Fatalf("unexpected public identifier: got %s want %s", got.PublicIdentifier, publicID)
	}
}

func TestPublishRepositoryHelloWithInvalidJetStreamContextReturnsError(t *testing.T) {
	nc := newMockNATSConnection(t)
	ctx := context.Background()

	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}

	_ = nc.Drain()

	repo := buffereddatasender.NewNATSDataPublisherRepository(nc, js, ctx)
	err = repo.Hello(uuid.New(), "gateway-public-id")
	if err == nil {
		t.Fatal("expected hello error with invalid jetstream context")
	}
}
