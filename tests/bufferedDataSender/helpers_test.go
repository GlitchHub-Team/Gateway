package buffereddatasender_test

import (
	"context"
	"net"
	"testing"

	"github.com/nats-io/nkeys"
)

func getFreePort(t *testing.T) int {
	t.Helper()

	lc := net.ListenConfig{}
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate free tcp port: %v", err)
	}

	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener addr type: %T", ln.Addr())
	}

	err = ln.Close()
	if err != nil {
		t.Fatalf("failed to close listener: %v", err)
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

	// The mock server does not validate the JWT, but the client-side signer
	// still requires a syntactically valid user seed.
	return "test-user-jwt", string(seedBytes)
}
