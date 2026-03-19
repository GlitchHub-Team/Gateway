package sensorprofilestests

import (
	"testing"
	"time"

	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

func TestHealthThermometerGenerate(t *testing.T) {
	//verifica che la temperatura generata derivi correttamente dal valore del mock random
	sensorID := uuid.New()
	profile := profiles.NewHealthThermometerProfile(sensorID, &MockRandomGenerator{nFloat: 0.5})

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

	gotData, ok := got.Data.(*profiles.HealthThermometerData)
	if !ok {
		t.Fatalf("expected data type *HealthThermometerData, got %T", got.Data)
	}

	if gotData.TemperatureValue != 37 {
		t.Fatalf("expected TemperatureValue 37, got %v", gotData.TemperatureValue)
	}
}
