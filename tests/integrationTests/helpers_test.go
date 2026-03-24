package integration_tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"
	gatewayservices "Gateway/internal/gatewayManager/gatewayServices"
	sensorprofiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nkeys"
	"go.uber.org/zap"
)

type natsCreds struct {
	JWT    string
	Seed   string
	Public string
}

func cleanupSQLiteFiles(t *testing.T) {
	t.Helper()
	_ = os.Remove("gateway.db")
	_ = os.Remove("gateway.db-shm")
	_ = os.Remove("gateway.db-wal")
	_ = os.Remove("buffer.db")
	_ = os.Remove("buffer.db-shm")
	_ = os.Remove("buffer.db-wal")
}

func findRepoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working directory: %v", err)
	}

	cur := wd
	for {
		if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
			return cur
		}

		parent := filepath.Dir(cur)
		if parent == cur {
			t.Fatalf("cannot find repository root from %s", wd)
		}
		cur = parent
	}
}

func absolutePathInCmd(repoRoot, pathValue string) string {
	if filepath.IsAbs(pathValue) {
		return pathValue
	}
	return filepath.Join(repoRoot, "cmd", pathValue)
}

func absolutePathInInfrastructure(repoRoot, pathValue string) string {
	if filepath.IsAbs(pathValue) {
		return pathValue
	}
	return filepath.Join(repoRoot, "Infrastructure", pathValue)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func readNatsCreds(t *testing.T, credsPath string) natsCreds {
	t.Helper()

	content, err := os.ReadFile(credsPath)
	if err != nil {
		t.Fatalf("cannot read creds file %s: %v", credsPath, err)
	}

	jwt := extractCredsBlock(t, string(content), "-----BEGIN NATS USER JWT-----", "------END NATS USER JWT------")
	seed := extractCredsBlock(t, string(content), "-----BEGIN USER NKEY SEED-----", "------END USER NKEY SEED------")

	kp, err := nkeys.FromSeed([]byte(seed))
	if err != nil {
		t.Fatalf("cannot parse seed from creds %s: %v", credsPath, err)
	}
	public, err := kp.PublicKey()
	if err != nil {
		t.Fatalf("cannot derive public key from seed %s: %v", credsPath, err)
	}

	return natsCreds{
		JWT:    jwt,
		Seed:   seed,
		Public: public,
	}
}

func extractCredsBlock(t *testing.T, content, begin, end string) string {
	t.Helper()

	start := strings.Index(content, begin)
	if start < 0 {
		t.Fatalf("cannot find creds block %s", begin)
	}
	start += len(begin)
	endIndex := strings.Index(content[start:], end)
	if endIndex < 0 {
		t.Fatalf("cannot find creds block end %s", end)
	}

	block := content[start : start+endIndex]
	return strings.TrimSpace(block)
}

func startController(t *testing.T, fx *commissionFixture, controller commandcontrollers.NATSCommandController) {
	t.Helper()

	controller.Listen()
	if err := fx.controllerNc.FlushTimeout(2 * time.Second); err != nil {
		t.Fatalf("cannot flush controller subscription: %v", err)
	}
}

func sendCommand(t *testing.T, nc *nats.Conn, subject string, body any) gatewayservices.Response {
	t.Helper()

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("cannot marshal payload for %s: %v", subject, err)
	}

	msg, err := nc.Request(subject, payload, 3*time.Second)
	if err != nil {
		t.Fatalf("cannot request %s: %v", subject, err)
	}

	var res gatewayservices.Response
	if err := json.Unmarshal(msg.Data, &res); err != nil {
		t.Fatalf("cannot unmarshal response for %s: %v", subject, err)
	}

	return res
}

func countRows(t *testing.T, ctx context.Context, db *sql.DB, query string, args ...any) int {
	t.Helper()

	var count int
	if err := db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		t.Fatalf("cannot count rows: %v", err)
	}
	return count
}

func getGatewayState(t *testing.T, ctx context.Context, db *sql.DB, gatewayID uuid.UUID) (tenantID *string, token *string, status string, interval int64) {
	t.Helper()

	err := db.QueryRowContext(
		ctx,
		`SELECT tenantId, token, status, interval FROM gateways WHERE id = ?`,
		gatewayID.String(),
	).Scan(&tenantID, &token, &status, &interval)
	if err != nil {
		t.Fatalf("cannot get gateway state for %s: %v", gatewayID, err)
	}

	return tenantID, token, status, interval
}

func getSensorStatus(t *testing.T, ctx context.Context, db *sql.DB, gatewayID, sensorID uuid.UUID) string {
	t.Helper()

	var status string
	err := db.QueryRowContext(
		ctx,
		`SELECT status FROM sensors WHERE gatewayId = ? AND id = ?`,
		gatewayID.String(),
		sensorID.String(),
	).Scan(&status)
	if err != nil {
		t.Fatalf("cannot get sensor status for gateway=%s sensor=%s: %v", gatewayID, sensorID, err)
	}

	return status
}

func mustNoRespondersOrTimeout(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, nats.ErrNoResponders.Error()) || strings.Contains(errMsg, nats.ErrTimeout.Error())
}

func responseMustFailWithFormatError(t *testing.T, res gatewayservices.Response) {
	t.Helper()
	if res.Success {
		t.Fatalf("expected failure response, got %+v", res)
	}
	if !strings.Contains(res.Message, "Formato del comando incorretto") {
		t.Fatalf("expected format error, got: %q", res.Message)
	}
}

func responseMustSucceed(t *testing.T, res gatewayservices.Response, expectedMessage string) {
	t.Helper()
	if !res.Success {
		t.Fatalf("expected success response, got %+v", res)
	}
	if expectedMessage != "" && res.Message != expectedMessage {
		t.Fatalf("unexpected success message: %q", res.Message)
	}
}

func responseMustFailContaining(t *testing.T, res gatewayservices.Response, expected string) {
	t.Helper()
	if res.Success {
		t.Fatalf("expected failure response, got %+v", res)
	}
	if !strings.Contains(res.Message, expected) {
		t.Fatalf("expected message containing %q, got: %q", expected, res.Message)
	}
}

func nopLogger() *zap.Logger {
	return zap.NewNop()
}

func newRand() sensorprofiles.Rand {
	return sensorprofiles.NewRand()
}

func ensureHelloStream(t *testing.T, ctx context.Context, nc *nats.Conn) {
	t.Helper()

	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatalf("cannot create jetstream context for hello stream: %v", err)
	}

	if _, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "HELLO_STREAM",
		Subjects: []string{"gateway.hello.*"},
	}); err != nil {
		t.Fatalf("cannot create/update hello stream: %v", err)
	}
}
