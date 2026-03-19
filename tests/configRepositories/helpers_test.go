package configrepositoriestests

import (
	"context"
	"testing"
	"time"

	gatewaydatabase "Gateway/cmd/external/gatewayDatabase"
	configrepositories "Gateway/internal/configManager/configRepositories"
	credentialsgenerator "Gateway/internal/credentialsGenerator"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensorpkg "Gateway/internal/sensor"
	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

type stubRand struct{}

func (s *stubRand) Intn(_ int) int { return 0 }
func (s *stubRand) Float64() float64 { return 0 }

func newRepository(t *testing.T) (*configrepositories.SQLiteConfigRepository, *configrepositories.ConfigDbConnection) {
	t.Helper()

	conn, err := gatewaydatabase.NewGatewayDatabase(context.Background())
	if err != nil {
		t.Fatalf("expected gateway db to open, got %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	return configrepositories.NewSQLiteConfigRepository(context.Background(), conn), conn
}

func insertGatewayRow(t *testing.T, conn *configrepositories.ConfigDbConnection, gatewayID uuid.UUID, tenantID *uuid.UUID, status domain.GatewayStatus, interval time.Duration, publicIdentifier, secretKey string, token *string) {
	t.Helper()

	query := `INSERT INTO gateways (id, tenantId, status, interval, publicIdentifier, secretKey, token) VALUES (?, ?, ?, ?, ?, ?, ?)`

	var tenantValue any
	if tenantID != nil {
		tenantValue = tenantID.String()
	}

	var tokenValue any
	if token != nil {
		tokenValue = *token
	}

	if _, err := conn.ExecContext(context.Background(), query, gatewayID.String(), tenantValue, string(status), interval.Milliseconds(), publicIdentifier, secretKey, tokenValue); err != nil {
		t.Fatalf("expected gateway insert to succeed, got %v", err)
	}
}

func insertSensorRow(t *testing.T, conn *configrepositories.ConfigDbConnection, sensorID, gatewayID uuid.UUID, profile string, status sensorpkg.SensorStatus, interval time.Duration) {
	t.Helper()

	query := `INSERT INTO sensors (id, gatewayId, profile, status, interval) VALUES (?, ?, ?, ?, ?)`
	if _, err := conn.ExecContext(context.Background(), query, sensorID.String(), gatewayID.String(), profile, string(status), interval.Milliseconds()); err != nil {
		t.Fatalf("expected sensor insert to succeed, got %v", err)
	}
}

func newCredentials() *credentialsgenerator.Credentials {
	return &credentialsgenerator.Credentials{
		PublicIdentifier: "public-id",
		SecretKey:        "secret-key",
	}
}

func newCreateGatewayCmd(gatewayID uuid.UUID, interval time.Duration) *commanddata.CreateGateway {
	return &commanddata.CreateGateway{
		GatewayId: gatewayID,
		Interval:  interval,
	}
}

func newHeartRateProfile(sensorID uuid.UUID) profiles.SensorProfile {
	return profiles.NewHeartRateProfile(sensorID, &stubRand{})
}
