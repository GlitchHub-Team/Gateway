package sensorprofilestests

import (
	"testing"
	"time"

	profiles "Gateway/internal/sensor/sensorProfiles"
	testutils "Gateway/tests/utils"

	"github.com/google/uuid"
)

func TestPulseOximeterGenerate(t *testing.T) {
	// Verifica che il profilo produca insieme saturazione e battito con i valori attesi.
	sensorID := uuid.New()
	profile := profiles.NewPulseOximeterProfile(sensorID, &testutils.MockRandomGenerator{NInt: 10, NFloat: 0.5})

	got := profile.Generate()

	if got.SensorId != sensorID {
		t.Fatalf("expected SensorId %v, got %v", sensorID, got.SensorId)
	}

	if got.Profile != profile.String() {
		t.Fatalf("expected Profile %q, got %q", profile.String(), got.Profile)
	}

	if !got.Timestamp.After(time.Now().Add(-time.Minute)) {
		t.Fatalf("expected recent Timestamp, got %v", got.Timestamp)
	}

	gotData, ok := got.Data.(*profiles.PulseOximeterData)
	if !ok {
		t.Fatalf("expected data type *PulseOximeterData, got %T", got.Data)
	}

	if gotData.SpO2Value != 97.5 {
		t.Fatalf("expected SpO2Value 97.5, got %v", gotData.SpO2Value)
	}

	if gotData.PulseRateValue != 70 {
		t.Fatalf("expected PulseRateValue 70, got %v", gotData.PulseRateValue)
	}
}
