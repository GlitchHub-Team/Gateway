package natsutiltests

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	natsserver "Gateway/cmd/external/natsServer"
	"Gateway/internal/natsutil"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func getFreePort(t *testing.T) int {
	t.Helper()

	lc := net.ListenConfig{}
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate free tcp port: %v", err)
	}
	defer func() { _ = ln.Close() }()

	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener addr type: %T", ln.Addr())
	}

	return tcpAddr.Port
}

func newMockNATSCreds(t *testing.T) (token string, seed string) {
	t.Helper()

	kp, err := nkeys.CreateUser()
	if err != nil {
		t.Fatalf("failed to create nkey pair: %v", err)
	}

	seedBytes, err := kp.Seed()
	if err != nil {
		t.Fatalf("failed to extract nkey seed: %v", err)
	}

	return "test-user-jwt", string(seedBytes)
}

func writeCredsFile(t *testing.T, dir string, token string, seed string) string {
	t.Helper()

	path := filepath.Join(dir, "test.creds")
	content := fmt.Sprintf("-----BEGIN NATS USER JWT-----\n%s\n------END NATS USER JWT------\n\n************************* IMPORTANT *************************\nNKEY Seed printed below can be used to sign and prove identity.\nNKEYs are sensitive and should be treated as secrets.\n\n-----BEGIN USER NKEY SEED-----\n%s\n------END USER NKEY SEED------\n", token, seed)

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write creds file: %v", err)
	}

	return path
}

func TestNewNATSConnectionWithoutCredsConnectsToMockServer(t *testing.T) {
	token, seed := newMockNATSCreds(t)
	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	serverConn := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = serverConn.Drain() })

	conn := natsutil.NewNATSConnection(host, port, "", "")
	t.Cleanup(func() { _ = conn.Drain() })

	if !conn.IsConnected() {
		t.Fatal("expected NATS connection to be established")
	}
}

func TestNewNATSConnectionWithCredsFileConnectsToMockServer(t *testing.T) {
	token, seed := newMockNATSCreds(t)
	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	serverConn := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = serverConn.Drain() })

	credsPath := writeCredsFile(t, t.TempDir(), token, seed)
	conn := natsutil.NewNATSConnection(host, port, natsutil.NatsCredsPath(credsPath), "")
	t.Cleanup(func() { _ = conn.Drain() })

	if !conn.IsConnected() {
		t.Fatal("expected NATS connection with creds file to be established")
	}
}

func TestNewJetStreamContextReturnsContextForValidConnection(t *testing.T) {
	token, seed := newMockNATSCreds(t)
	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	conn := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = conn.Drain() })

	js, err := natsutil.NewJetStreamContext(conn)
	if err != nil {
		t.Fatalf("expected jetstream context creation to succeed, got %v", err)
	}

	if _, err := js.AccountInfo(context.Background()); err != nil {
		t.Fatalf("expected jetstream context to be usable, got %v", err)
	}
}

func TestJWTAuthWithInvalidSeedFailsToConnect(t *testing.T) {
	token, seed := newMockNATSCreds(t)
	host := natsutil.NatsAddress("127.0.0.1")
	port := natsutil.NatsPort(getFreePort(t))

	serverConn := natsserver.NewMockNATSConnection(host, port, natsutil.NatsToken(token), natsutil.NatsSeed(seed))
	t.Cleanup(func() { _ = serverConn.Drain() })

	url := fmt.Sprintf("nats://%s:%d", host, port)
	conn, err := nats.Connect(url, natsutil.JWTAuth(token, "not-a-valid-seed"))
	if err == nil {
		_ = conn.Drain()
		t.Fatal("expected connection with invalid seed to fail")
	}
}
