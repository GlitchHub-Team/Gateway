package integration_tests

import (
	"testing"
	"time"

	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"
)

func TestNATSRebootGatewayIntegration(t *testing.T) {
	t.Run("reboot corretto con hello", func(t *testing.T) {
		t.Setenv("BASE_CREDS_PATH", "admin_test.creds")
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)
		ensureHelloStream(t, fx.ctx, fx.controllerNc)

		rebootController := commandcontrollers.NewNATSRebootGatewayController(
			fx.controllerNc,
			commandcontrollers.RebootGatewaySubject("commands.rebootgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, rebootController)

		commissionRes := fx.sendCommissionCommand(t, adminGatewayID, adminTenantID, fx.adminCreds(t).JWT)
		responseMustSucceed(t, commissionRes, "Gateway commissionato correttamente")

		helloSubject := "gateway.hello." + adminGatewayID.String()
		sub, err := fx.observerNc.SubscribeSync(helloSubject)
		if err != nil {
			t.Fatalf("cannot subscribe hello subject: %v", err)
		}
		if err := fx.observerNc.FlushTimeout(2 * time.Second); err != nil {
			t.Fatalf("cannot flush hello subscription: %v", err)
		}

		res := sendCommand(t, fx.publisherNc, "commands.rebootgateway", map[string]any{
			"gatewayId": adminGatewayID.String(),
		})
		responseMustSucceed(t, res, "Gateway riavviato con successo")

		if _, err := sub.NextMsg(3 * time.Second); err != nil {
			t.Fatalf("expected hello publish after reboot: %v", err)
		}
	})

	t.Run("gateway non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSRebootGatewayController(
			fx.controllerNc,
			commandcontrollers.RebootGatewaySubject("commands.rebootgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.rebootgateway", map[string]any{
			"gatewayId": "44444444-4444-4444-4444-444444444444",
		})
		responseMustFailContaining(t, res, "nessun gateway trovato")
	})

	t.Run("gatewayId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSRebootGatewayController(
			fx.controllerNc,
			commandcontrollers.RebootGatewaySubject("commands.rebootgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.rebootgateway", map[string]any{
			"gatewayId": "not-a-uuid",
		})
		responseMustFailWithFormatError(t, res)
	})
}
