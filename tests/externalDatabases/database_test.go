package externaldatabasetests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	bufferdatabase "Gateway/cmd/external/bufferDatabase"
	gatewaydatabase "Gateway/cmd/external/gatewayDatabase"
)

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change working dir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("failed to restore working dir: %v", err)
		}
	})
}

func TestNewBufferDatabaseCreatesFileAndTable(t *testing.T) {
	withWorkingDir(t, t.TempDir())

	conn := bufferdatabase.NewBufferDatabase()
	t.Cleanup(func() { _ = conn.Close() })

	if _, err := os.Stat(filepath.Join(".", "buffer.db")); err != nil {
		t.Fatalf("expected buffer.db to be created, got %v", err)
	}

	var tableName string
	if err := conn.QueryRowContext(context.Background(), "SELECT name FROM sqlite_master WHERE type = ? AND name = ?", "table", "buffer").Scan(&tableName); err != nil {
		t.Fatalf("expected buffer table to exist, got %v", err)
	}
	if tableName != "buffer" {
		t.Fatalf("expected buffer table, got %s", tableName)
	}
}

func TestNewGatewayDatabaseCreatesFileTablesAndEnablesForeignKeys(t *testing.T) {
	withWorkingDir(t, t.TempDir())

	conn, err := gatewaydatabase.NewGatewayDatabase(context.Background())
	if err != nil {
		t.Fatalf("expected gateway db to open, got %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	if _, err := os.Stat(filepath.Join(".", "gateway.db")); err != nil {
		t.Fatalf("expected gateway.db to be created, got %v", err)
	}

	for _, table := range []string{"gateways", "sensors"} {
		var tableName string
		if err := conn.QueryRowContext(context.Background(), "SELECT name FROM sqlite_master WHERE type = ? AND name = ?", "table", table).Scan(&tableName); err != nil {
			t.Fatalf("expected table %s to exist, got %v", table, err)
		}
		if tableName != table {
			t.Fatalf("expected table %s, got %s", table, tableName)
		}
	}

	var foreignKeysEnabled int
	if err := conn.QueryRowContext(context.Background(), "PRAGMA foreign_keys;").Scan(&foreignKeysEnabled); err != nil {
		t.Fatalf("expected foreign_keys pragma query to succeed, got %v", err)
	}
	if foreignKeysEnabled != 1 {
		t.Fatalf("expected foreign_keys pragma enabled, got %d", foreignKeysEnabled)
	}
}
