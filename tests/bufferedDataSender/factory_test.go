package buffereddatasendertests

import (
	"strings"
	"testing"

	buffereddatasender "Gateway/internal/bufferedDataSender"
)

func TestCreateReturnsPublisherThatUsesProvidedJetStream(t *testing.T) {
	js := &fakeJetStreamContext{}
	factory := newFactory(js)

	port := factory.Create()
	if err := port.Hello(validUUID(t), "public-key"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if js.publishCalls != 1 {
		t.Fatalf("expected one publish through provided jetstream, got %d", js.publishCalls)
	}
}

func TestReloadReturnsConnectionErrorWhenNATSIsUnavailable(t *testing.T) {
	factory := buffereddatasender.NewNATSDataPublisherFactory(nil, "127.0.0.1", 1)

	_, err := factory.Reload("token", validSeed(t))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "errore nella connessione a NATS") {
		t.Fatalf("expected connection error context, got %q", err.Error())
	}
}
