package sensorprofilestests

import (
	"slices"
	"testing"
	"time"

	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

func TestEcgGenerate(t *testing.T) {
	sensorId := uuid.New()
	tests := []struct {
		name    string
		profile *profiles.EcgProfile
		got     *profiles.GeneratedSensorData
		want    *profiles.GeneratedSensorData
	}{
		{
			name:    "Ecg correct generate data",
			profile: profiles.NewEcgProfile(sensorId, &MockRandomGenerator{nInt: 100}),
			want: &profiles.GeneratedSensorData{
				SensorId:  sensorId,
				Timestamp: time.Time{},
				Data:      &profiles.EcgData{Waveform: slices.Repeat([]int{-400}, 250)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.profile.Generate()

			if got.SensorId != tt.want.SensorId {
				t.Errorf("Expected SensorId %v, got %v", tt.want.SensorId, got.SensorId)
			}

			if !got.Timestamp.After(time.Now().Add(-time.Minute)) {
				t.Errorf("Expected Timestamp to be recent, got %v", got.Timestamp)
			}

			gotData, ok := got.Data.(*profiles.EcgData)
			if !ok {
				t.Fatalf("Expected data type *EcgData, got %T", got.Data)
			}

			wantData := tt.want.Data.(*profiles.EcgData)
			if !slices.Equal(gotData.Waveform, wantData.Waveform) {
				t.Errorf("Expected Waveform %v, got %v", wantData.Waveform, gotData.Waveform)
			}
		})
	}
}
