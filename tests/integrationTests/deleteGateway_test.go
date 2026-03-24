package integration_tests

import (
	"testing"

	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"

	"github.com/google/uuid"
)

func TestNATSDeleteGatewayIntegration(t *testing.T) {
	t.Run("delete gateway corretto", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		addController := commandcontrollers.NewNATSAddSensorController(
			fx.controllerNc,
			commandcontrollers.AddSensorSubject("commands.addsensor"),
			fx.service,
			newRand(),
			nopLogger(),
		)
		deleteController := commandcontrollers.NewNATSDeleteGatewayController(
			fx.controllerNc,
			commandcontrollers.DeleteGatewaySubject("commands.deletegateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, addController)
		startController(t, fx, deleteController)

		sensorID := uuid.New()
		addRes := sendCommand(t, fx.publisherNc, "commands.addsensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  sensorID.String(),
			"profile":   "HeartRate",
			"interval":  50,
		})
		responseMustSucceed(t, addRes, "Sensore aggiunto con successo")

		res := sendCommand(t, fx.publisherNc, "commands.deletegateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustSucceed(t, res, "Gateway eliminato con successo")

		if countRows(t, fx.ctx, fx.gatewayDb.DB, `SELECT COUNT(*) FROM gateways WHERE id = ?`, gateway1ID.String()) != 0 {
			t.Fatalf("expected deleted gateway row")
		}
		if countRows(t, fx.ctx, fx.gatewayDb.DB, `SELECT COUNT(*) FROM sensors WHERE gatewayId = ?`, gateway1ID.String()) != 0 {
			t.Fatalf("expected deleted sensors by cascade")
		}

		secondDelete := sendCommand(t, fx.publisherNc, "commands.deletegateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustFailContaining(t, secondDelete, "non trovato")
	})

	t.Run("gateway non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		deleteController := commandcontrollers.NewNATSDeleteGatewayController(
			fx.controllerNc,
			commandcontrollers.DeleteGatewaySubject("commands.deletegateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, deleteController)

		res := sendCommand(t, fx.publisherNc, "commands.deletegateway", map[string]any{
			"gatewayId": uuid.New().String(),
		})
		responseMustFailContaining(t, res, "non trovato")
	})

	t.Run("gatewayId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		deleteController := commandcontrollers.NewNATSDeleteGatewayController(
			fx.controllerNc,
			commandcontrollers.DeleteGatewaySubject("commands.deletegateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, deleteController)

		res := sendCommand(t, fx.publisherNc, "commands.deletegateway", map[string]any{
			"gatewayId": "not-a-uuid",
		})
		responseMustFailWithFormatError(t, res)
	})
}
