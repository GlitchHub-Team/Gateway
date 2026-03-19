package sensorprofilestests

import (
	"testing"

	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

func TestParseSensorProfile(t *testing.T) {
	//verifica che ogni stringa valida venga tradotta nel tipo concreto di profilo corretto
	tests := []struct {
		name        string
		profileType string
		wantType    any
		wantString  string
	}{
		{
			name:        "returns ECG profile",
			profileType: "ECG",
			wantType:    &profiles.EcgProfile{},
			wantString:  "ECG",
		},
		{
			name:        "returns EnvironmentalSensing profile",
			profileType: "EnvironmentalSensing",
			wantType:    &profiles.EnvironmentalSensingProfile{},
			wantString:  "EnvironmentalSensing",
		},
		{
			name:        "returns HealthThermometer profile",
			profileType: "HealthThermometer",
			wantType:    &profiles.HealthThermometerProfile{},
			wantString:  "HealthThermometer",
		},
		{
			name:        "returns HeartRate profile",
			profileType: "HeartRate",
			wantType:    &profiles.HeartRateProfile{},
			wantString:  "HeartRate",
		},
		{
			name:        "returns PulseOximeter profile",
			profileType: "PulseOximeter",
			wantType:    &profiles.PulseOximeterProfile{},
			wantString:  "PulseOximeter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := profiles.ParseSensorProfile(tt.profileType, &MockRandomGenerator{})
			if got == nil {
				t.Fatal("expected non-nil profile")
			}

			switch tt.wantType.(type) {
			case *profiles.EcgProfile:
				if _, ok := got.(*profiles.EcgProfile); !ok {
					t.Fatalf("expected *EcgProfile, got %T", got)
				}
			case *profiles.EnvironmentalSensingProfile:
				if _, ok := got.(*profiles.EnvironmentalSensingProfile); !ok {
					t.Fatalf("expected *EnvironmentalSensingProfile, got %T", got)
				}
			case *profiles.HealthThermometerProfile:
				if _, ok := got.(*profiles.HealthThermometerProfile); !ok {
					t.Fatalf("expected *HealthThermometerProfile, got %T", got)
				}
			case *profiles.HeartRateProfile:
				if _, ok := got.(*profiles.HeartRateProfile); !ok {
					t.Fatalf("expected *HeartRateProfile, got %T", got)
				}
			case *profiles.PulseOximeterProfile:
				if _, ok := got.(*profiles.PulseOximeterProfile); !ok {
					t.Fatalf("expected *PulseOximeterProfile, got %T", got)
				}
			}

			generatedData := got.Generate()
			if generatedData.SensorId == uuid.Nil {
				t.Fatal("expected generated sensor id to be populated")
			}

			if generatedData.Profile != tt.wantString {
				t.Fatalf("expected generated Profile %q, got %q", tt.wantString, generatedData.Profile)
			}
		})
	}
}

func TestParseSensorProfileUnknown(t *testing.T) {
	//verifica che un nome profilo non supportato ritorni null
	got := profiles.ParseSensorProfile("Unknown", &MockRandomGenerator{})
	if got != nil {
		t.Fatalf("expected nil profile, got %T", got)
	}
}
