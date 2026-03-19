package buffereddatasendertests

import (
	"context"
	"sync"
	"testing"
	"time"

	bufferdatabase "Gateway/cmd/external/bufferDatabase"
	gatewaydatabase "Gateway/cmd/external/gatewayDatabase"
	configmanager "Gateway/internal/configManager"
	credentialsgenerator "Gateway/internal/credentialsGenerator"
	"Gateway/internal/domain"
	buffereddatasender "Gateway/internal/bufferedDataSender"
	sensorpkg "Gateway/internal/sensor"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type fakeJetStreamContext struct {
	mu            sync.Mutex
	publishCalls  int
	publishedSubj string
	publishedData []byte
	publishErr    error
	publishSignal chan struct{}
}

func (f *fakeJetStreamContext) Publish(subj string, data []byte, _ ...nats.PubOpt) (*nats.PubAck, error) {
	f.mu.Lock()
	f.publishCalls++
	f.publishedSubj = subj
	f.publishedData = append([]byte(nil), data...)
	signal := f.publishSignal
	err := f.publishErr
	f.mu.Unlock()

	if signal != nil {
		select {
		case signal <- struct{}{}:
		default:
		}
	}

	if err != nil {
		return nil, err
	}

	return &nats.PubAck{}, nil
}

func (f *fakeJetStreamContext) PublishMsg(_ *nats.Msg, _ ...nats.PubOpt) (*nats.PubAck, error) {
	return &nats.PubAck{}, nil
}

func (f *fakeJetStreamContext) PublishAsync(_ string, _ []byte, _ ...nats.PubOpt) (nats.PubAckFuture, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) PublishMsgAsync(_ *nats.Msg, _ ...nats.PubOpt) (nats.PubAckFuture, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) PublishAsyncPending() int {
	return 0
}

func (f *fakeJetStreamContext) PublishAsyncComplete() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) CleanupPublisher() {}

func (f *fakeJetStreamContext) Subscribe(_ string, _ nats.MsgHandler, _ ...nats.SubOpt) (*nats.Subscription, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) SubscribeSync(_ string, _ ...nats.SubOpt) (*nats.Subscription, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) ChanSubscribe(_ string, _ chan *nats.Msg, _ ...nats.SubOpt) (*nats.Subscription, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) ChanQueueSubscribe(_ string, _ string, _ chan *nats.Msg, _ ...nats.SubOpt) (*nats.Subscription, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) QueueSubscribe(_ string, _ string, _ nats.MsgHandler, _ ...nats.SubOpt) (*nats.Subscription, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) QueueSubscribeSync(_ string, _ string, _ ...nats.SubOpt) (*nats.Subscription, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) PullSubscribe(_ string, _ string, _ ...nats.SubOpt) (*nats.Subscription, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) AddStream(_ *nats.StreamConfig, _ ...nats.JSOpt) (*nats.StreamInfo, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) UpdateStream(_ *nats.StreamConfig, _ ...nats.JSOpt) (*nats.StreamInfo, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) DeleteStream(_ string, _ ...nats.JSOpt) error {
	return nil
}

func (f *fakeJetStreamContext) StreamInfo(_ string, _ ...nats.JSOpt) (*nats.StreamInfo, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) PurgeStream(_ string, _ ...nats.JSOpt) error {
	return nil
}

func (f *fakeJetStreamContext) StreamsInfo(_ ...nats.JSOpt) <-chan *nats.StreamInfo {
	ch := make(chan *nats.StreamInfo)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) Streams(_ ...nats.JSOpt) <-chan *nats.StreamInfo {
	ch := make(chan *nats.StreamInfo)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) StreamNames(_ ...nats.JSOpt) <-chan string {
	ch := make(chan string)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) GetMsg(_ string, _ uint64, _ ...nats.JSOpt) (*nats.RawStreamMsg, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) GetLastMsg(_ string, _ string, _ ...nats.JSOpt) (*nats.RawStreamMsg, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) DeleteMsg(_ string, _ uint64, _ ...nats.JSOpt) error {
	return nil
}

func (f *fakeJetStreamContext) SecureDeleteMsg(_ string, _ uint64, _ ...nats.JSOpt) error {
	return nil
}

func (f *fakeJetStreamContext) AddConsumer(_ string, _ *nats.ConsumerConfig, _ ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) UpdateConsumer(_ string, _ *nats.ConsumerConfig, _ ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) DeleteConsumer(_ string, _ string, _ ...nats.JSOpt) error {
	return nil
}

func (f *fakeJetStreamContext) ConsumerInfo(_ string, _ string, _ ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) ConsumersInfo(_ string, _ ...nats.JSOpt) <-chan *nats.ConsumerInfo {
	ch := make(chan *nats.ConsumerInfo)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) Consumers(_ string, _ ...nats.JSOpt) <-chan *nats.ConsumerInfo {
	ch := make(chan *nats.ConsumerInfo)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) ConsumerNames(_ string, _ ...nats.JSOpt) <-chan string {
	ch := make(chan string)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) AccountInfo(_ ...nats.JSOpt) (*nats.AccountInfo, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) StreamNameBySubject(_ string, _ ...nats.JSOpt) (string, error) {
	return "", nil
}

func (f *fakeJetStreamContext) KeyValue(_ string) (nats.KeyValue, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) CreateKeyValue(_ *nats.KeyValueConfig) (nats.KeyValue, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) DeleteKeyValue(_ string) error {
	return nil
}

func (f *fakeJetStreamContext) KeyValueStoreNames() <-chan string {
	ch := make(chan string)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) KeyValueStores() <-chan nats.KeyValueStatus {
	ch := make(chan nats.KeyValueStatus)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) ObjectStore(_ string) (nats.ObjectStore, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) CreateObjectStore(_ *nats.ObjectStoreConfig) (nats.ObjectStore, error) {
	return nil, nil
}

func (f *fakeJetStreamContext) DeleteObjectStore(_ string) error {
	return nil
}

func (f *fakeJetStreamContext) ObjectStoreNames(_ ...nats.ObjectOpt) <-chan string {
	ch := make(chan string)
	close(ch)
	return ch
}

func (f *fakeJetStreamContext) ObjectStores(_ ...nats.ObjectOpt) <-chan nats.ObjectStoreStatus {
	ch := make(chan nats.ObjectStoreStatus)
	close(ch)
	return ch
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

func newBufferTestDB(t *testing.T) sensorpkg.BufferDbConnection {
	t.Helper()

	conn := bufferdatabase.NewBufferDatabase()
	t.Cleanup(func() {
		_ = conn.Close()
	})
	return conn
}

func newNonBufferDB(t *testing.T) sensorpkg.BufferDbConnection {
	t.Helper()

	conn, err := gatewaydatabase.NewGatewayDatabase(context.Background())
	if err != nil {
		t.Fatalf("expected gateway db to open, got %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	return sensorpkg.BufferDbConnection{DB: conn.DB}
}

func newGateway(status domain.GatewayStatus, interval time.Duration) *configmanager.Gateway {
	return &configmanager.Gateway{
		Id:               uuid.New(),
		Status:           status,
		Interval:         interval,
		Sensors:          make(map[uuid.UUID]*sensorpkg.Sensor),
		PublicIdentifier: "public-key",
		SecretKey:        "secret-key",
	}
}

func newFactory(js nats.JetStreamContext) *buffereddatasender.NATSDataPublisherFactory {
	return buffereddatasender.NewNATSDataPublisherFactory(js, "127.0.0.1", 4222)
}

func validSeed(t *testing.T) string {
	t.Helper()

	credentials, err := credentialsgenerator.GenerateCredentials()
	if err != nil {
		t.Fatalf("expected credentials generation, got %v", err)
	}

	return credentials.SecretKey
}

func validUUID(t *testing.T) uuid.UUID {
	t.Helper()
	return uuid.New()
}

func waitForSignal(t *testing.T, ch <-chan struct{}, label string) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", label)
	}
}
