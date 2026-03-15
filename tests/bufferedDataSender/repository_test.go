package buffereddatasendertests

import (
	"context"
	"strings"
	"testing"
	"time"

	buffereddatasender "Gateway/internal/bufferedDataSender"
	sensorpkg "Gateway/internal/sensor"

	"github.com/google/uuid"
)

func TestGetOrderedBufferedDataReturnsRowsOrderedByTimestamp(t *testing.T) {
	conn := newBufferTestDB(t)
	repo := buffereddatasender.NewBufferedDataRepository(context.Background(), conn)

	gatewayID := uuid.New()
	firstSensorID := uuid.New()
	secondSensorID := uuid.New()
	firstTimestamp := time.Date(2026, time.March, 14, 10, 0, 0, 0, time.UTC)
	secondTimestamp := time.Date(2026, time.March, 14, 11, 0, 0, 0, time.UTC)

	insertBufferRow(t, conn, gatewayID, secondSensorID, secondTimestamp, "HeartRate", `{"BpmValue":70}`)
	insertBufferRow(t, conn, gatewayID, firstSensorID, firstTimestamp, "HeartRate", `{"BpmValue":60}`)

	data, err := repo.GetOrderedBufferedData(gatewayID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(data) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(data))
	}

	if data[0].SensorId != firstSensorID || !data[0].Timestamp.Equal(firstTimestamp) {
		t.Fatalf("expected first row ordered by timestamp, got %+v", data[0])
	}

	if data[1].SensorId != secondSensorID || !data[1].Timestamp.Equal(secondTimestamp) {
		t.Fatalf("expected second row ordered by timestamp, got %+v", data[1])
	}

	if string(data[0].Data) != `{"BpmValue":60}` {
		t.Fatalf("expected json payload, got %v", data[0].Data)
	}
}

func TestGetOrderedBufferedDataReturnsQueryError(t *testing.T) {
	db := newNonBufferDB(t)
	repo := buffereddatasender.NewBufferedDataRepository(context.Background(), db)
	gatewayID := uuid.New()

	_, err := repo.GetOrderedBufferedData(gatewayID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "errore nell'eseguire la query per ottenere i dati del buffer") {
		t.Fatalf("expected query context, got %q", err.Error())
	}
}

func TestCleanBufferedDataReturnsNilWhenSliceIsEmpty(t *testing.T) {
	conn := newBufferTestDB(t)
	repo := buffereddatasender.NewBufferedDataRepository(context.Background(), conn)

	if err := repo.CleanBufferedData(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCleanWholeBufferDeletesRowsForGateway(t *testing.T) {
	conn := newBufferTestDB(t)
	repo := buffereddatasender.NewBufferedDataRepository(context.Background(), conn)

	gatewayID := uuid.New()
	otherGatewayID := uuid.New()
	insertBufferRow(t, conn, gatewayID, uuid.New(), time.Date(2026, time.March, 14, 14, 0, 0, 0, time.UTC), "HeartRate", `{"BpmValue":100}`)
	insertBufferRow(t, conn, otherGatewayID, uuid.New(), time.Date(2026, time.March, 14, 15, 0, 0, 0, time.UTC), "HeartRate", `{"BpmValue":110}`)

	if err := repo.CleanWholeBuffer(gatewayID); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var count int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer`).Scan(&count); err != nil {
		t.Fatalf("expected count query to succeed, got %v", err)
	}

	if count != 1 {
		t.Fatalf("expected one remaining row, got %d", count)
	}
}

func insertBufferRow(t *testing.T, conn sensorpkg.BufferDbConnection, gatewayID, sensorID uuid.UUID, timestamp time.Time, profile, payload string) {
	t.Helper()

	query := `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, jsonb(?))`
	if _, err := conn.ExecContext(context.Background(), query, gatewayID, sensorID, timestamp, profile, payload); err != nil {
		t.Fatalf("expected insert to succeed, got %v", err)
	}
}
