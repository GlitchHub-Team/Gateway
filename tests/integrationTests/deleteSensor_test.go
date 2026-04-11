package integration_tests

import (
	"testing"

	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"

	"github.com/google/uuid"
)

func TestNATSDeleteSensorIntegration(t *testing.T) {
	t.Run("delete sensore corretto", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		addController := commandcontrollers.NewNATSAddSensorController(
			fx.controllerNc,
			commandcontrollers.AddSensorSubject("commands.addsensor"),
			fx.service,
			newRand(),
			nopLogger(),
		)
		deleteController := commandcontrollers.NewNATSDeleteSensorController(
			fx.controllerNc,
			commandcontrollers.DeleteSensorSubject("commands.deletesensor"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, addController)
		startController(t, fx, deleteController)

		sensorID := uuid.New()
		addRes := sendCommand(t, fx.publisherNc, "commands.addsensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  sensorID.String(),
			"profile":   "heart_rate",
			"interval":  40,
		})
		responseMustSucceed(t, addRes, "Sensore aggiunto con successo")

		res := sendCommand(t, fx.publisherNc, "commands.deletesensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  sensorID.String(),
		})
		responseMustSucceed(t, res, "Sensore eliminato con successo")

		if countRows(t, fx.ctx, fx.gatewayDb.DB, `SELECT COUNT(*) FROM sensors WHERE gatewayId = ? AND id = ?`, gateway1ID.String(), sensorID.String()) != 0 {
			t.Fatalf("expected deleted sensor row")
		}

		secondDelete := sendCommand(t, fx.publisherNc, "commands.deletesensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  sensorID.String(),
		})
		responseMustFailContaining(t, secondDelete, "non trovato")
	})

	t.Run("gateway non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		deleteController := commandcontrollers.NewNATSDeleteSensorController(
			fx.controllerNc,
			commandcontrollers.DeleteSensorSubject("commands.deletesensor"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, deleteController)

		res := sendCommand(t, fx.publisherNc, "commands.deletesensor", map[string]any{
			"gatewayId": uuid.New().String(),
			"sensorId":  uuid.New().String(),
		})
		responseMustFailContaining(t, res, "non trovato")
	})

	t.Run("id invalidi", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		deleteController := commandcontrollers.NewNATSDeleteSensorController(
			fx.controllerNc,
			commandcontrollers.DeleteSensorSubject("commands.deletesensor"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, deleteController)

		resGateway := sendCommand(t, fx.publisherNc, "commands.deletesensor", map[string]any{
			"gatewayId": "not-a-uuid",
			"sensorId":  uuid.New().String(),
		})
		responseMustFailWithFormatError(t, resGateway)

		resSensor := sendCommand(t, fx.publisherNc, "commands.deletesensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  "not-a-uuid",
		})
		responseMustFailWithFormatError(t, resSensor)
	})
}
