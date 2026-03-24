package integration_tests

import (
	"testing"
	"time"

	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"
	"Gateway/internal/sensor"

	"github.com/google/uuid"
)

func TestNATSInterruptSensorIntegration(t *testing.T) {
	t.Run("interrupt sensore corretto", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		addController := commandcontrollers.NewNATSAddSensorController(
			fx.controllerNc,
			commandcontrollers.AddSensorSubject("commands.addsensor"),
			fx.service,
			newRand(),
			nopLogger(),
		)
		interruptController := commandcontrollers.NewNATSInterruptSensorController(
			fx.controllerNc,
			commandcontrollers.InterruptSensorSubject("commands.interruptsensor"),
			fx.service,
			nopLogger(),
		)
		interruptGatewayController := commandcontrollers.NewNATSInterruptGatewayController(
			fx.controllerNc,
			commandcontrollers.InterruptGatewaySubject("commands.interruptgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, addController)
		startController(t, fx, interruptController)
		startController(t, fx, interruptGatewayController)

		commissionRes := fx.sendCommissionCommand(t, gateway1ID, tenant1ID, fx.gateway1Creds(t).JWT)
		responseMustSucceed(t, commissionRes, "Gateway commissionato correttamente")

		interruptGatewayRes := sendCommand(t, fx.publisherNc, "commands.interruptgateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustSucceed(t, interruptGatewayRes, "Gateway interrotto con successo")

		sensorID := uuid.New()
		addRes := sendCommand(t, fx.publisherNc, "commands.addsensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  sensorID.String(),
			"profile":   "HeartRate",
			"interval":  60,
		})
		responseMustSucceed(t, addRes, "Sensore aggiunto con successo")

		time.Sleep(260 * time.Millisecond)
		rowsBeforeInterrupt := countRows(
			t,
			fx.ctx,
			fx.bufferDb.DB,
			`SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`,
			gateway1ID,
		)
		if rowsBeforeInterrupt == 0 {
			t.Fatalf("expected buffered data before sensor interrupt")
		}

		res := sendCommand(t, fx.publisherNc, "commands.interruptsensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  sensorID.String(),
		})
		responseMustSucceed(t, res, "Sensore interrotto con successo")

		time.Sleep(260 * time.Millisecond)
		rowsAfterInterrupt := countRows(
			t,
			fx.ctx,
			fx.bufferDb.DB,
			`SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`,
			gateway1ID,
		)
		if rowsAfterInterrupt != rowsBeforeInterrupt {
			t.Fatalf("expected no new buffered data after sensor interrupt, before=%d after=%d", rowsBeforeInterrupt, rowsAfterInterrupt)
		}

		status := getSensorStatus(t, fx.ctx, fx.gatewayDb.DB, gateway1ID, sensorID)
		if status != string(sensor.Inactive) {
			t.Fatalf("expected inactive sensor after interrupt, got %q", status)
		}
	})

	t.Run("gateway non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSInterruptSensorController(
			fx.controllerNc,
			commandcontrollers.InterruptSensorSubject("commands.interruptsensor"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.interruptsensor", map[string]any{
			"gatewayId": uuid.New().String(),
			"sensorId":  uuid.New().String(),
		})
		responseMustFailContaining(t, res, "nessun gateway trovato")
	})

	t.Run("sensore non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSInterruptSensorController(
			fx.controllerNc,
			commandcontrollers.InterruptSensorSubject("commands.interruptsensor"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.interruptsensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  uuid.New().String(),
		})
		responseMustFailContaining(t, res, "nessun sensore trovato")
	})

	t.Run("id invalidi", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSInterruptSensorController(
			fx.controllerNc,
			commandcontrollers.InterruptSensorSubject("commands.interruptsensor"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		resGateway := sendCommand(t, fx.publisherNc, "commands.interruptsensor", map[string]any{
			"gatewayId": "not-a-uuid",
			"sensorId":  uuid.New().String(),
		})
		responseMustFailWithFormatError(t, resGateway)

		resSensor := sendCommand(t, fx.publisherNc, "commands.interruptsensor", map[string]any{
			"gatewayId": gateway1ID.String(),
			"sensorId":  "not-a-uuid",
		})
		responseMustFailWithFormatError(t, resSensor)
	})
}
