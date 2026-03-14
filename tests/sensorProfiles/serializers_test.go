package sensorprofilestests

import (
	"testing"

	profiles "Gateway/internal/sensor/sensorProfiles"
)

func TestSerialize(t *testing.T) {
	// Verifica che ogni payload sensore venga serializzato in JSON senza errori.
	tests := []struct {
		name string
		data profiles.SerializableData
		want string
	}{
		{
			name: "serializes ECG data",
			data: &profiles.EcgData{Waveform: []int{-400, 0, 1200}},
			want: "{\"Waveform\":[-400,0,1200]}",
		},
		{
			name: "serializes heart rate data",
			data: &profiles.HeartRateData{BpmValue: 70},
			want: "{\"BpmValue\":70}",
		},
		{
			name: "serializes pulse oximeter data",
			data: &profiles.PulseOximeterData{SpO2Value: 97.5, PulseRateValue: 70},
			want: "{\"SpO2Value\":97.5,\"PulseRateValue\":70}",
		},
		{
			name: "serializes health thermometer data",
			data: &profiles.HealthThermometerData{TemperatureValue: 37},
			want: "{\"TemperatureValue\":37}",
		},
		{
			name: "serializes environmental sensing data",
			data: &profiles.EnvironmentalSensingData{TemperatureValue: 22.5, HumidityValue: 50, PressureValue: 1005},
			want: "{\"TemperatureValue\":22.5,\"HumidityValue\":50,\"PressureValue\":1005}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.data.Serialize()
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			if string(got) != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, string(got))
			}
		})
	}
}
