package sensorprofilestests

import (
	"testing"
	"time"

	profiles "Gateway/internal/sensor/sensorProfiles"
	testutils "Gateway/tests/utils"

	"github.com/google/uuid"
)

func TestHeartRateGenerate(t *testing.T) {
	// Verifica che il profilo converta il valore casuale nel BPM previsto.
	sensorID := uuid.New()
	profile := profiles.NewHeartRateProfile(sensorID, &testutils.MockRandomGenerator{NInt: 10})

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

	gotData, ok := got.Data.(*profiles.HeartRateData)
	if !ok {
		t.Fatalf("expected data type *HeartRateData, got %T", got.Data)
	}

	if gotData.BpmValue != 70 {
		t.Fatalf("expected BpmValue 70, got %v", gotData.BpmValue)
	}
}
