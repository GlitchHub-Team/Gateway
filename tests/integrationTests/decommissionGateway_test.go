package integration_tests

import (
	"testing"
	"time"

	"Gateway/internal/domain"
	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"

	"github.com/google/uuid"
)

func TestNATSDecommissionGatewayIntegration(t *testing.T) {
	t.Run("decommission corretto", func(t *testing.T) {
		t.Setenv("BASE_CREDS_PATH", "admin_test.creds")
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)
		ensureHelloStream(t, fx.ctx, fx.controllerNc)

		interruptController := commandcontrollers.NewNATSInterruptGatewayController(
			fx.controllerNc,
			commandcontrollers.InterruptGatewaySubject("commands.interruptgateway"),
			fx.service,
			nopLogger(),
		)
		controller := commandcontrollers.NewNATSDecommissionGatewayController(
			fx.controllerNc,
			commandcontrollers.DecommissionGatewaySubject("commands.decommissiongateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, interruptController)
		startController(t, fx, controller)

		commissionRes := fx.sendCommissionCommand(t, gateway1ID, tenant1ID, fx.gateway1Creds(t).JWT)
		responseMustSucceed(t, commissionRes, "Gateway commissionato correttamente")

		interruptRes := sendCommand(t, fx.publisherNc, "commands.interruptgateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustSucceed(t, interruptRes, "Gateway interrotto con successo")

		fx.insertBufferedData(t, gateway1ID, uuid.New())
		bufferedBefore := countRows(t, fx.ctx, fx.bufferDb.DB, `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gateway1ID)
		if bufferedBefore == 0 {
			t.Fatalf("expected buffered rows before decommission")
		}

		helloSubject := "gateway.hello." + gateway1ID.String()
		sub, err := fx.observerNc.SubscribeSync(helloSubject)
		if err != nil {
			t.Fatalf("cannot subscribe hello subject: %v", err)
		}
		if err := fx.observerNc.FlushTimeout(2 * time.Second); err != nil {
			t.Fatalf("cannot flush hello subscription: %v", err)
		}

		res := sendCommand(t, fx.publisherNc, "commands.decommissiongateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustSucceed(t, res, "Gateway decommissionato correttamente")

		if _, err := sub.NextMsg(3 * time.Second); err != nil {
			t.Fatalf("expected hello publish after decommission: %v", err)
		}

		tenantID, token, status, _ := getGatewayState(t, fx.ctx, fx.gatewayDb.DB, gateway1ID)
		if tenantID != nil || token != nil || status != string(domain.Decommissioned) {
			t.Fatalf("unexpected gateway state after decommission")
		}

		bufferedAfter := countRows(t, fx.ctx, fx.bufferDb.DB, `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gateway1ID)
		if bufferedAfter != 0 {
			t.Fatalf("expected empty buffer after decommission, got %d rows", bufferedAfter)
		}
	})

	t.Run("gateway non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSDecommissionGatewayController(
			fx.controllerNc,
			commandcontrollers.DecommissionGatewaySubject("commands.decommissiongateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.decommissiongateway", map[string]any{
			"gatewayId": uuid.New().String(),
		})
		responseMustFailContaining(t, res, "nessun gateway trovato")
	})

	t.Run("gatewayId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSDecommissionGatewayController(
			fx.controllerNc,
			commandcontrollers.DecommissionGatewaySubject("commands.decommissiongateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.decommissiongateway", map[string]any{
			"gatewayId": "not-a-uuid",
		})
		responseMustFailWithFormatError(t, res)
	})

	t.Run("gateway gia decommissionato", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSDecommissionGatewayController(
			fx.controllerNc,
			commandcontrollers.DecommissionGatewaySubject("commands.decommissiongateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.decommissiongateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustFailContaining(t, res, "gia decommissionato")
	})
}
