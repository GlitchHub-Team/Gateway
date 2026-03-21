package buffereddatasender_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
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

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve caller path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func parseNATSCreds(t *testing.T, credsPath string) (token string, seed string) {
	t.Helper()
	content, err := os.ReadFile(credsPath)
	if err != nil {
		t.Fatalf("unable to read creds file %s: %v", credsPath, err)
	}
	text := string(content)

	token = between(text, "-----BEGIN NATS USER JWT-----", "------END NATS USER JWT------")
	seed = between(text, "-----BEGIN USER NKEY SEED-----", "------END USER NKEY SEED------")
	if token == "" || seed == "" {
		t.Fatalf("unable to parse token/seed from %s", credsPath)
	}
	return token, seed
}

func between(text, start, end string) string {
	startIdx := strings.Index(text, start)
	if startIdx < 0 {
		return ""
	}
	startIdx += len(start)

	endIdx := strings.Index(text[startIdx:], end)
	if endIdx < 0 {
		return ""
	}
	return strings.TrimSpace(text[startIdx : startIdx+endIdx])
}
