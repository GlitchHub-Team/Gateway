package buffereddatasender_test

import (
	"context"
	"strings"
	"testing"
	_ "unsafe"

	natsserver "Gateway/cmd/external/natsServer"
	buffereddatasender "Gateway/internal/bufferedDataSender"
	"Gateway/internal/natsutil"
)

func TestNATSDataPublisherFactoryCreateReturnsRepository(t *testing.T) {
	factory := buffereddatasender.NewNATSDataPublisherFactory(nil, nil, "127.0.0.1", 4222, context.Background(), "")

	port := factory.Create()
	if _, ok := port.(*buffereddatasender.NATSDataPublisherRepository); !ok {
		t.Fatalf("expected *NATSDataPublisherRepository, got %T", port)
	}
}

func TestNATSDataPublisherFactoryReloadValidConnectionWithMockNATS(t *testing.T) {
	token, seed := newMockNATSCreds(t)

	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	nc := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = nc.Drain() })

	js, err := natsutil.NewJetStreamContext(nc)
	if err != nil {
		t.Fatalf("unable to create jetstream context: %v", err)
	}

	factory := buffereddatasender.NewNATSDataPublisherFactory(js, nc, natsutil.NatsAddress(host), natsutil.NatsPort(port), context.Background(), "")

	portImpl, err := factory.Reload(token, seed)
	if err != nil {
		t.Fatalf("expected reload to succeed, got: %v", err)
	}

	repo, ok := portImpl.(*buffereddatasender.NATSDataPublisherRepository)
	if !ok {
		t.Fatalf("expected *NATSDataPublisherRepository, got %T", portImpl)
	}
	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
}

func TestNATSDataPublisherFactoryReloadInvalidConnectionWithMockNATS(t *testing.T) {
	token, seed := newMockNATSCreds(t)

	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	nc := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = nc.Drain() })

	unreachablePort := natsutil.NatsPort(getFreePort(t))
	factory := buffereddatasender.NewNATSDataPublisherFactory(nil, nc, natsutil.NatsAddress(host), natsutil.NatsPort(unreachablePort), context.Background(), "")

	_, err := factory.Reload(token, seed)
	if err == nil {
		t.Fatal("expected reload to fail with unreachable mock server port")
	}
}

func TestNATSDataPublisherFactoryReloadInvalidAddress(t *testing.T) {
	factory := buffereddatasender.NewNATSDataPublisherFactory(nil, nil, "bad host", 4222, context.Background(), "")

	_, err := factory.Reload("token", "seed")
	if err == nil {
		t.Fatal("expected reload failure for invalid host")
	}
}

func TestNATSDataPublisherFactoryReloadInvalidPort(t *testing.T) {
	factory := buffereddatasender.NewNATSDataPublisherFactory(nil, nil, "127.0.0.1", -1, context.Background(), "")

	_, err := factory.Reload("token", "seed")
	if err == nil {
		t.Fatal("expected reload failure for invalid port")
	}
}

func TestNATSDataPublisherFactoryReloadInvalidCAPemPath(t *testing.T) {
	token, seed := newMockNATSCreds(t)

	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	nc := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = nc.Drain() })

	factory := buffereddatasender.NewNATSDataPublisherFactory(nil, nc, natsutil.NatsAddress(host), natsutil.NatsPort(port), context.Background(), "/does/not/exist/ca.pem")

	_, err := factory.Reload(token, seed)
	if err == nil {
		t.Fatal("expected reload failure for invalid ca pem path")
	}
}

func TestNATSDataPublisherFactoryReloadAcceptsCredsFormattedToken(t *testing.T) {
	token, seed := newMockNATSCreds(t)

	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	nc := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = nc.Drain() })

	factory := buffereddatasender.NewNATSDataPublisherFactory(nil, nc, natsutil.NatsAddress(host), natsutil.NatsPort(port), context.Background(), "")

	credsLikeToken := strings.Join([]string{
		"-----BEGIN NATS USER JWT-----",
		token,
		"------END NATS USER JWT------",
		"",
		"-----BEGIN USER NKEY SEED-----",
		seed,
		"------END USER NKEY SEED------",
	}, "\n")

	_, err := factory.Reload(credsLikeToken, seed)
	if err != nil {
		t.Fatalf("expected reload to succeed with creds-formatted token, got: %v", err)
	}
}

func TestNATSDataPublisherFactoryReloadEmptyTokenFailsFast(t *testing.T) {
	factory := buffereddatasender.NewNATSDataPublisherFactory(nil, nil, "127.0.0.1", 4222, context.Background(), "")

	_, err := factory.Reload(" \n\t ", "seed")
	if err == nil {
		t.Fatal("expected reload failure for empty token")
	}
	if !strings.Contains(err.Error(), "commissioned token vuoto") {
		t.Fatalf("expected empty-token error, got: %v", err)
	}
}
