package sensortests

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	bufferdatabase "Gateway/cmd/external/bufferDatabase"
	sensor "Gateway/internal/sensor"
	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type serializeErrorData struct{}

func (s *serializeErrorData) Serialize() ([]byte, error) {
	return nil, errors.New("serialize failed")
}

func TestSaveReturnsErrorWhenSerializationFails(t *testing.T) {
	conn := bufferdatabase.NewMockBufferDatabase()
	t.Cleanup(func() { _ = conn.Close() })

	repo := sensor.NewSQLiteSaveSensorDataRepository(context.Background(), conn)
	data := &profiles.GeneratedSensorData{
		SensorId:  uuid.New(),
		Timestamp: time.Now().UTC(),
		Profile:   "BrokenProfile",
		Data:      &serializeErrorData{},
	}

	err := repo.Save(data, uuid.New())
	if err == nil {
		t.Fatal("expected serialization error")
	}
	if !strings.Contains(err.Error(), "errore nella serializzazione dei dati") {
		t.Fatalf("expected serialization error message, got %v", err)
	}
}

func TestSaveReturnsErrorWhenInsertFails(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := sensor.NewSQLiteSaveSensorDataRepository(context.Background(), sensor.BufferDbConnection{DB: db})
	data := &profiles.GeneratedSensorData{
		SensorId:  uuid.New(),
		Timestamp: time.Now().UTC(),
		Profile:   "heart_rate",
		Data:      &mockSerializableData{},
	}

	err = repo.Save(data, uuid.New())
	if err == nil {
		t.Fatal("expected insert error")
	}
	if !strings.Contains(err.Error(), "errore nel salvataggio del dato") {
		t.Fatalf("expected insert error message, got %v", err)
	}
}
