package sensorprofilestests

import (
	"testing"
	"time"

	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

func TestEnvironmentalSensingGenerate(t *testing.T) {
	//verifica che il profilo generi il payload atteso e conservi i metadati del sensore
	sensorID := uuid.New()
	profile := profiles.NewEnvironmentalSensingProfile(sensorID, &MockRandomGenerator{nFloat: 0.5})

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

	gotData, ok := got.Data.(*profiles.EnvironmentalSensingData)
	if !ok {
		t.Fatalf("expected data type *EnvironmentalSensingData, got %T", got.Data)
	}

	if gotData.TemperatureValue != 22.5 {
		t.Fatalf("expected TemperatureValue 22.5, got %v", gotData.TemperatureValue)
	}

	if gotData.HumidityValue != 50 {
		t.Fatalf("expected HumidityValue 50, got %v", gotData.HumidityValue)
	}

	if gotData.PressureValue != 1005 {
		t.Fatalf("expected PressureValue 1005, got %v", gotData.PressureValue)
	}
}
