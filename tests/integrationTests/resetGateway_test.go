package integration_tests

import (
	"testing"

	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"

	"github.com/google/uuid"
)

func TestNATSResetGatewayIntegration(t *testing.T) {
	t.Run("reset corretto", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSResetGatewayController(
			fx.controllerNc,
			commandcontrollers.ResetGatewaySubject("commands.resetgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		fx.insertBufferedData(t, gateway1ID, uuid.New())
		bufferedBefore := countRows(t, fx.ctx, fx.bufferDb.DB, `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gateway1ID)
		if bufferedBefore == 0 {
			t.Fatalf("expected buffered rows before reset")
		}

		res := sendCommand(t, fx.publisherNc, "commands.resetgateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustSucceed(t, res, "Gateway resettato con successo")

		_, _, _, interval := getGatewayState(t, fx.ctx, fx.gatewayDb.DB, gateway1ID)
		if interval != 5000 {
			t.Fatalf("expected reset default interval 5000ms, got %d", interval)
		}

		bufferedAfter := countRows(t, fx.ctx, fx.bufferDb.DB, `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gateway1ID)
		if bufferedAfter != 0 {
			t.Fatalf("expected empty buffer after reset, got %d rows", bufferedAfter)
		}
	})

	t.Run("gateway non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSResetGatewayController(
			fx.controllerNc,
			commandcontrollers.ResetGatewaySubject("commands.resetgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.resetgateway", map[string]any{
			"gatewayId": uuid.New().String(),
		})
		responseMustFailContaining(t, res, "nessun gateway trovato")
	})

	t.Run("gatewayId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSResetGatewayController(
			fx.controllerNc,
			commandcontrollers.ResetGatewaySubject("commands.resetgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.resetgateway", map[string]any{
			"gatewayId": "not-a-uuid",
		})
		responseMustFailWithFormatError(t, res)
	})
}
