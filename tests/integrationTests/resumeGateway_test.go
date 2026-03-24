package integration_tests

import (
	"testing"
	"time"

	"Gateway/internal/domain"
	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"

	"github.com/google/uuid"
)

func TestNATSResumeGatewayIntegration(t *testing.T) {
	t.Run("resume corretto", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		interruptController := commandcontrollers.NewNATSInterruptGatewayController(
			fx.controllerNc,
			commandcontrollers.InterruptGatewaySubject("commands.interruptgateway"),
			fx.service,
			nopLogger(),
		)
		resumeController := commandcontrollers.NewNATSResumeGatewayController(
			fx.controllerNc,
			commandcontrollers.ResumeGatewaySubject("commands.resumegateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, interruptController)
		startController(t, fx, resumeController)

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
		if _, err := sub.NextMsg(300 * time.Millisecond); err != nil {
			t.Fatalf("expected sensor data publish before gateway interrupt: %v", err)
		}

		interruptRes := sendCommand(t, fx.publisherNc, "commands.interruptgateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustSucceed(t, interruptRes, "Gateway interrotto con successo")

		fx.insertBufferedData(t, gateway1ID, sensorID)
		if _, err := sub.NextMsg(150 * time.Millisecond); err == nil {
			t.Fatalf("expected no sensor data while gateway is interrupted")
		}

		resumeRes := sendCommand(t, fx.publisherNc, "commands.resumegateway", map[string]any{
			"gatewayId": gateway1ID.String(),
		})
		responseMustSucceed(t, resumeRes, "Gateway ripreso con successo")

		fx.insertBufferedData(t, gateway1ID, sensorID)
		if _, err := sub.NextMsg(300 * time.Millisecond); err != nil {
			t.Fatalf("expected sensor data publish after gateway resume: %v", err)
		}

		_, _, status, _ := getGatewayState(t, fx.ctx, fx.gatewayDb.DB, gateway1ID)
		if status != string(domain.Active) {
			t.Fatalf("expected active gateway after resume, got %q", status)
		}
	})

	t.Run("gateway non esistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSResumeGatewayController(
			fx.controllerNc,
			commandcontrollers.ResumeGatewaySubject("commands.resumegateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.resumegateway", map[string]any{
			"gatewayId": uuid.New().String(),
		})
		responseMustFailContaining(t, res, "nessun gateway trovato")
	})

	t.Run("gatewayId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		controller := commandcontrollers.NewNATSResumeGatewayController(
			fx.controllerNc,
			commandcontrollers.ResumeGatewaySubject("commands.resumegateway"),
			fx.service,
			nopLogger(),
		)
		startController(t, fx, controller)

		res := sendCommand(t, fx.publisherNc, "commands.resumegateway", map[string]any{
			"gatewayId": "not-a-uuid",
		})
		responseMustFailWithFormatError(t, res)
	})
}
