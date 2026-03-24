package integration_tests

import (
	"testing"
	"time"

	"Gateway/internal/domain"
	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"

	"github.com/google/uuid"
)

func TestNATSInterruptGatewayIntegration(t *testing.T) {
	t.Run("interrupt corretto", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSInterruptGatewayController(
			fx.controllerNc,
			commandcontrollers.InterruptGatewaySubject("commands.interruptgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		commissionRes := fx.sendCommissionCommand(t, gateway1ID, tenant1ID, fx.gateway1Creds(t).JWT)
		responseMustSucceed(t, commissionRes, "Gateway commissionato correttamente")

		sensorID := uuid.New()
		subject := sensorSubject(tenant1ID, gateway1ID, sensorID)
		sub, err := fx.observerNc.SubscribeSync(subject)
		if err != nil {
			t.Fatalf("cannot subscribe to sensor subject: %v", err)
		}
		if err := fx.observerNc.FlushTimeout(2 * time.Second); err != nil {
			t.Fatalf("cannot flush observer subscription: %v", err)
		}

		fx.insertBufferedData(t, gateway1ID, sensorID)

		if _, err := sub.NextMsg(2 * time.Second); err != nil {
			t.Fatalf("expected sensor data publish before interrupt: %v", err)
		}

		res := sendCommand(t, fx.publisherNc, "commands.interruptgateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustSucceed(t, res, "Gateway interrotto con successo")

		_, _, status, _ := getGatewayState(t, fx.ctx, fx.gatewayDb.DB, gateway1ID)
		if status != string(domain.Inactive) {
			t.Fatalf("expected inactive gateway after interrupt, got %q", status)
		}

		postInterruptSub, err := fx.observerNc.SubscribeSync(subject)
		if err != nil {
			t.Fatalf("cannot subscribe to sensor subject after interrupt: %v", err)
		}
		if err := fx.observerNc.FlushTimeout(2 * time.Second); err != nil {
			t.Fatalf("cannot flush post-interrupt subscription: %v", err)
		}

		fx.insertBufferedData(t, gateway1ID, sensorID)
		if _, err := postInterruptSub.NextMsg(250 * time.Millisecond); err == nil {
			t.Fatalf("expected no sensor data after gateway interrupt")
		}
	})

	t.Run("gateway non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSInterruptGatewayController(
			fx.controllerNc,
			commandcontrollers.InterruptGatewaySubject("commands.interruptgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.interruptgateway", map[string]any{
			"gatewayId": uuid.New().String(),
		})
		responseMustFailContaining(t, res, "nessun gateway trovato")
	})

	t.Run("gatewayId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSInterruptGatewayController(
			fx.controllerNc,
			commandcontrollers.InterruptGatewaySubject("commands.interruptgateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.interruptgateway", map[string]any{
			"gatewayId": "not-a-uuid",
		})
		responseMustFailWithFormatError(t, res)
	})
}
