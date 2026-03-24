package integration_tests

import (
	"testing"
	"time"
)

func TestNATSStopSensorCommandNotExposed(t *testing.T) {
	fx := newCommissionFixture(t, commissionFixtureOptions{})
	defer fx.close(t)

	_, err := fx.publisherNc.Request("commands.stopsensor", []byte(`{"gatewayId":"x","sensorId":"y"}`), 500*time.Millisecond)
	if !mustNoRespondersOrTimeout(err) {
		t.Fatalf("expected no responders or timeout for commands.stopsensor, got err=%v", err)
	}
}
