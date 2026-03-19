package configrepositoriestests

import (
	"context"
	"testing"
	"time"

	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensorpkg "Gateway/internal/sensor"

	"github.com/google/uuid"
)

func TestCreateGatewayPersistsGatewayRow(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	cmd := newCreateGatewayCmd(gatewayID, 3*time.Second)

	if err := repo.CreateGateway(cmd, newCredentials(), domain.Decommissioned); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var status string
	var interval int64
	var publicIdentifier, secretKey string
	var tenantID, token *string
	err := conn.QueryRowContext(context.Background(), `SELECT tenantId, status, interval, publicIdentifier, secretKey, token FROM gateways WHERE id = ?`, gatewayID.String()).
		Scan(&tenantID, &status, &interval, &publicIdentifier, &secretKey, &token)
	if err != nil {
		t.Fatalf("expected gateway row, got %v", err)
	}

	if tenantID != nil || token != nil {
		t.Fatalf("expected nil tenant and token, got %v %v", tenantID, token)
	}

	if status != string(domain.Decommissioned) || interval != cmd.Interval.Milliseconds() || publicIdentifier != "public-id" || secretKey != "secret-key" {
		t.Fatalf("unexpected persisted gateway row")
	}
}

func TestCommissionGatewayUpdatesTenantTokenAndStatus(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	tenantID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, nil, domain.Decommissioned, time.Second, "public", "secret", nil)

	cmd := &commanddata.CommissionGateway{
		GatewayId:         gatewayID,
		TenantId:          tenantID,
		CommissionedToken: "jwt-token",
	}

	if err := repo.CommissionGateway(cmd, domain.Active); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var gotTenant, gotToken, gotStatus string
	if err := conn.QueryRowContext(context.Background(), `SELECT tenantId, token, status FROM gateways WHERE id = ?`, gatewayID.String()).Scan(&gotTenant, &gotToken, &gotStatus); err != nil {
		t.Fatalf("expected updated gateway row, got %v", err)
	}

	if gotTenant != tenantID.String() || gotToken != "jwt-token" || gotStatus != string(domain.Active) {
		t.Fatalf("unexpected commissioned row")
	}
}

func TestDecommissionGatewayClearsTenantTokenAndUpdatesStatus(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	tenantID := uuid.New()
	token := "jwt-token"
	insertGatewayRow(t, conn, gatewayID, &tenantID, domain.Active, time.Second, "public", "secret", &token)

	if err := repo.DecommissionGateway(&commanddata.DecommissionGateway{GatewayId: gatewayID}, domain.Decommissioned); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var gotTenant, gotToken *string
	var gotStatus string
	if err := conn.QueryRowContext(context.Background(), `SELECT tenantId, token, status FROM gateways WHERE id = ?`, gatewayID.String()).Scan(&gotTenant, &gotToken, &gotStatus); err != nil {
		t.Fatalf("expected updated gateway row, got %v", err)
	}

	if gotTenant != nil || gotToken != nil || gotStatus != string(domain.Decommissioned) {
		t.Fatalf("unexpected decommissioned row")
	}
}

func TestInterruptAndResumeGatewayUpdateStatus(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, nil, domain.Active, time.Second, "public", "secret", nil)

	if err := repo.InterruptGateway(&commanddata.InterruptGateway{GatewayId: gatewayID}, domain.Inactive); err != nil {
		t.Fatalf("expected nil error interrupting gateway, got %v", err)
	}

	var status string
	if err := conn.QueryRowContext(context.Background(), `SELECT status FROM gateways WHERE id = ?`, gatewayID.String()).Scan(&status); err != nil {
		t.Fatalf("expected gateway row, got %v", err)
	}

	if status != string(domain.Inactive) {
		t.Fatalf("expected inactive status, got %q", status)
	}

	if err := repo.ResumeGateway(&commanddata.ResumeGateway{GatewayId: gatewayID}, domain.Active); err != nil {
		t.Fatalf("expected nil error resuming gateway, got %v", err)
	}

	if err := conn.QueryRowContext(context.Background(), `SELECT status FROM gateways WHERE id = ?`, gatewayID.String()).Scan(&status); err != nil {
		t.Fatalf("expected gateway row, got %v", err)
	}

	if status != string(domain.Active) {
		t.Fatalf("expected active status, got %q", status)
	}
}

func TestResetGatewayUpdatesInterval(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, nil, domain.Active, time.Second, "public", "secret", nil)

	if err := repo.ResetGateway(&commanddata.ResetGateway{GatewayId: gatewayID}, 7*time.Second); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var interval int64
	if err := conn.QueryRowContext(context.Background(), `SELECT interval FROM gateways WHERE id = ?`, gatewayID.String()).Scan(&interval); err != nil {
		t.Fatalf("expected gateway row, got %v", err)
	}

	if interval != 7000 {
		t.Fatalf("expected 7000, got %d", interval)
	}
}

func TestAddSensorPersistsSensorRow(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	sensorID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, nil, domain.Active, time.Second, "public", "secret", nil)

	cmd := &commanddata.AddSensor{
		GatewayId: gatewayID,
		SensorId:  sensorID,
		Profile:   newHeartRateProfile(sensorID),
		Interval:  2 * time.Second,
	}

	if err := repo.AddSensor(cmd, sensorpkg.Active); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var profile, status string
	var interval int64
	if err := conn.QueryRowContext(context.Background(), `SELECT profile, status, interval FROM sensors WHERE id = ? AND gatewayId = ?`, sensorID.String(), gatewayID.String()).Scan(&profile, &status, &interval); err != nil {
		t.Fatalf("expected sensor row, got %v", err)
	}

	if profile != "HeartRate" || status != string(sensorpkg.Active) || interval != 2000 {
		t.Fatalf("unexpected sensor row")
	}
}

func TestInterruptAndResumeSensorUpdateStatus(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	sensorID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, nil, domain.Active, time.Second, "public", "secret", nil)
	insertSensorRow(t, conn, sensorID, gatewayID, "HeartRate", sensorpkg.Active, time.Second)

	if err := repo.InterruptSensor(&commanddata.InterruptSensor{GatewayId: gatewayID, SensorId: sensorID}, sensorpkg.Inactive); err != nil {
		t.Fatalf("expected nil error interrupting sensor, got %v", err)
	}

	var status string
	if err := conn.QueryRowContext(context.Background(), `SELECT status FROM sensors WHERE id = ? AND gatewayId = ?`, sensorID.String(), gatewayID.String()).Scan(&status); err != nil {
		t.Fatalf("expected sensor row, got %v", err)
	}

	if status != string(sensorpkg.Inactive) {
		t.Fatalf("expected inactive status, got %q", status)
	}

	if err := repo.ResumeSensor(&commanddata.ResumeSensor{GatewayId: gatewayID, SensorId: sensorID}, sensorpkg.Active); err != nil {
		t.Fatalf("expected nil error resuming sensor, got %v", err)
	}

	if err := conn.QueryRowContext(context.Background(), `SELECT status FROM sensors WHERE id = ? AND gatewayId = ?`, sensorID.String(), gatewayID.String()).Scan(&status); err != nil {
		t.Fatalf("expected sensor row, got %v", err)
	}

	if status != string(sensorpkg.Active) {
		t.Fatalf("expected active status, got %q", status)
	}
}

func TestDeleteSensorDeletesRow(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	sensorID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, nil, domain.Active, time.Second, "public", "secret", nil)
	insertSensorRow(t, conn, sensorID, gatewayID, "HeartRate", sensorpkg.Active, time.Second)

	if err := repo.DeleteSensor(&commanddata.DeleteSensor{GatewayId: gatewayID, SensorId: sensorID}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var count int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM sensors WHERE id = ? AND gatewayId = ?`, sensorID.String(), gatewayID.String()).Scan(&count); err != nil {
		t.Fatalf("expected count query to succeed, got %v", err)
	}

	if count != 0 {
		t.Fatalf("expected deleted sensor, got %d", count)
	}
}

func TestDeleteGatewayDeletesGatewayAndCascadesSensors(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	sensorID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, nil, domain.Active, time.Second, "public", "secret", nil)
	insertSensorRow(t, conn, sensorID, gatewayID, "HeartRate", sensorpkg.Active, time.Second)

	if err := repo.DeleteGateway(&commanddata.DeleteGateway{GatewayId: gatewayID}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var gatewayCount int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM gateways WHERE id = ?`, gatewayID.String()).Scan(&gatewayCount); err != nil {
		t.Fatalf("expected gateway count query to succeed, got %v", err)
	}

	var sensorCount int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM sensors WHERE gatewayId = ?`, gatewayID.String()).Scan(&sensorCount); err != nil {
		t.Fatalf("expected sensor count query to succeed, got %v", err)
	}

	if gatewayCount != 0 || sensorCount != 0 {
		t.Fatalf("expected deleted gateway and cascaded sensors")
	}
}
