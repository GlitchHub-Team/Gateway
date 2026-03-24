package configrepositoriestests

import (
	"strings"
	"testing"
	"time"

	"Gateway/internal/domain"
	sensorpkg "Gateway/internal/sensor"

	"github.com/google/uuid"
)

func TestGetGatewayByIdLoadsGatewayAndSensors(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	tenantID := uuid.New()
	token := "jwt-token"
	sensorID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, &tenantID, domain.Active, 3*time.Second, "public", "secret", &token)
	insertSensorRow(t, conn, sensorID, gatewayID, "HeartRate", sensorpkg.Inactive, 2*time.Second)

	gateway, err := repo.GetGatewayById(gatewayID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if gateway.Id != gatewayID || gateway.TenantId == nil || *gateway.TenantId != tenantID || gateway.Token == nil || *gateway.Token != token {
		t.Fatalf("unexpected gateway identity fields: %+v", gateway)
	}

	if gateway.Status != domain.Active || gateway.Interval != 3*time.Second || gateway.PublicIdentifier != "public" || gateway.SecretKey != "secret" {
		t.Fatalf("unexpected gateway state fields: %+v", gateway)
	}

	sensorEntity, exists := gateway.Sensors[sensorID]
	if !exists {
		t.Fatal("expected sensor to be loaded")
	}

	if sensorEntity.GatewayId != gatewayID || sensorEntity.Status != sensorpkg.Inactive || sensorEntity.Interval != 2*time.Second {
		t.Fatalf("unexpected sensor loaded: %+v", sensorEntity)
	}

	if sensorEntity.Profile == nil || sensorEntity.Profile.String() != "HeartRate" {
		t.Fatalf("expected HeartRate profile, got %#v", sensorEntity.Profile)
	}
}

func TestGetGatewayByIdReturnsErrorWhenGatewayMissing(t *testing.T) {
	repo, _ := newRepository(t)
	gatewayID := uuid.New()

	_, err := repo.GetGatewayById(gatewayID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "fallito a recuperare il gateway") {
		t.Fatalf("expected missing gateway context, got %q", err.Error())
	}
}

func TestGetAllGatewaysLoadsAllGatewaysAndTheirSensors(t *testing.T) {
	repo, conn := newRepository(t)
	firstGatewayID := uuid.New()
	secondGatewayID := uuid.New()
	firstSensorID := uuid.New()
	secondSensorID := uuid.New()
	insertGatewayRow(t, conn, firstGatewayID, nil, domain.Decommissioned, time.Second, "public-1", "secret-1", nil)
	insertGatewayRow(t, conn, secondGatewayID, nil, domain.Active, 5*time.Second, "public-2", "secret-2", nil)
	insertSensorRow(t, conn, firstSensorID, firstGatewayID, "HeartRate", sensorpkg.Active, time.Second)
	insertSensorRow(t, conn, secondSensorID, secondGatewayID, "ECG", sensorpkg.Inactive, 2*time.Second)

	gateways, err := repo.GetAllGateways()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(gateways) != 2 {
		t.Fatalf("expected 2 gateways, got %d", len(gateways))
	}

	if gateways[firstGatewayID].Status != domain.Decommissioned || gateways[firstGatewayID].Sensors[firstSensorID].Profile.String() != "HeartRate" {
		t.Fatalf("unexpected first gateway loaded: %+v", gateways[firstGatewayID])
	}

	if gateways[secondGatewayID].Status != domain.Active || gateways[secondGatewayID].Sensors[secondSensorID].Profile.String() != "ECG" {
		t.Fatalf("unexpected second gateway loaded: %+v", gateways[secondGatewayID])
	}
}

func TestGetAllGatewaysReturnsErrorWhenSensorProfileIsInvalid(t *testing.T) {
	repo, conn := newRepository(t)
	gatewayID := uuid.New()
	sensorID := uuid.New()
	insertGatewayRow(t, conn, gatewayID, nil, domain.Active, time.Second, "public", "secret", nil)
	insertSensorRow(t, conn, sensorID, gatewayID, "UnknownProfile", sensorpkg.Active, time.Second)

	_, err := repo.GetAllGateways()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "fallito a recuperare i sensori dei gateway") {
		t.Fatalf("expected sensor loading context, got %q", err.Error())
	}
}
