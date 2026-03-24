package integration_tests

import (
	"strings"
	"testing"
	"time"

	"Gateway/internal/domain"
	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"

	"github.com/google/uuid"
)

func TestNATSCreateGatewayIntegration(t *testing.T) {
	t.Run("creazione gateway corretta", func(t *testing.T) {
		t.Setenv("BASE_CREDS_PATH", "admin_test.creds")
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)
		ensureHelloStream(t, fx.ctx, fx.controllerNc)

		controller := commandcontrollers.NewNATSCreateGatewayController(
			fx.controllerNc,
			commandcontrollers.CreateGatewaySubject("commands.creategateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		gatewayID := uuid.New()
		helloSubject := "gateway.hello." + gatewayID.String()
		sub, err := fx.observerNc.SubscribeSync(helloSubject)
		if err != nil {
			t.Fatalf("cannot subscribe hello subject: %v", err)
		}
		if err := fx.observerNc.FlushTimeout(2 * time.Second); err != nil {
			t.Fatalf("cannot flush hello subscription: %v", err)
		}

		res := sendCommand(t, fx.publisherNc, "commands.creategateway", map[string]any{
			"gatewayId": gatewayID.String(),
			"interval":  750,
		})
		responseMustSucceed(t, res, "Gateway creato con successo")

		if _, err := sub.NextMsg(3 * time.Second); err != nil {
			t.Fatalf("expected hello publish for created gateway: %v", err)
		}

		tenantID, token, status, interval := getGatewayState(t, fx.ctx, fx.gatewayDb.DB, gatewayID)
		if tenantID != nil || token != nil {
			t.Fatalf("expected uncommissioned gateway persistence")
		}
		if status != string(domain.Decommissioned) || interval != 750 {
			t.Fatalf("unexpected gateway persistence values")
		}

		duplicate := sendCommand(t, fx.publisherNc, "commands.creategateway", map[string]any{
			"gatewayId": gatewayID.String(),
			"interval":  750,
		})
		responseMustFailContaining(t, duplicate, "già esistente")
	})

	t.Run("gatewayId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSCreateGatewayController(
			fx.controllerNc,
			commandcontrollers.CreateGatewaySubject("commands.creategateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.creategateway", map[string]any{
			"gatewayId": "not-a-uuid",
			"interval":  100,
		})
		responseMustFailWithFormatError(t, res)
	})

	t.Run("interval non valido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSCreateGatewayController(
			fx.controllerNc,
			commandcontrollers.CreateGatewaySubject("commands.creategateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.creategateway", map[string]any{
			"gatewayId": uuid.New().String(),
			"interval":  0,
		})
		responseMustFailWithFormatError(t, res)
		if !strings.Contains(res.Message, "intervallo non valido") {
			t.Fatalf("expected interval validation in response, got: %q", res.Message)
		}
	})
}
