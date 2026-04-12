package buffereddatasender_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	bufferdatabase "Gateway/cmd/external/bufferDatabase"
	buffereddatasender "Gateway/internal/bufferedDataSender"
	sensor "Gateway/internal/sensor"
	sensorprofiles "Gateway/internal/sensor/sensorProfiles"
	testutils "Gateway/tests/utils"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

func newMockBufferRepository(t *testing.T) (*buffereddatasender.BufferedDataRepository, sensor.BufferDbConnection) {
	t.Helper()
	conn := bufferdatabase.NewMockBufferDatabase()
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	t.Cleanup(func() { _ = conn.Close() })
	return buffereddatasender.NewBufferedDataRepository(context.Background(), conn), conn
}

func seedGeneratedSensorData(t *testing.T, conn sensor.BufferDbConnection, gatewayID uuid.UUID, profile sensorprofiles.SensorProfile) (*sensorprofiles.GeneratedSensorData, []byte) {
	t.Helper()
	saveRepo := sensor.NewSQLiteSaveSensorDataRepository(context.Background(), conn)
	generated := profile.Generate()
	serialized, err := generated.Data.Serialize()
	if err != nil {
		t.Fatalf("failed to serialize generated data: %v", err)
	}
	if err := saveRepo.Save(generated, gatewayID); err != nil {
		t.Fatalf("failed to save generated data: %v", err)
	}
	return generated, serialized
}

func TestGetOrderedBufferedDataByGatewayAndValueColumn(t *testing.T) {
	repo, conn := newMockBufferRepository(t)
	gatewayID := uuid.New()

	rand := &testutils.MockRandomGenerator{NInt: 1, NFloat: 0.5}
	profiles := []sensorprofiles.SensorProfile{
		sensorprofiles.NewEcgProfile(uuid.New(), rand),
		sensorprofiles.NewEnvironmentalSensingProfile(uuid.New(), rand),
		sensorprofiles.NewHealthThermometerProfile(uuid.New(), rand),
		sensorprofiles.NewHeartRateProfile(uuid.New(), rand),
		sensorprofiles.NewPulseOximeterProfile(uuid.New(), rand),
	}

	expectedBySensor := make(map[uuid.UUID][]byte)
	for _, p := range profiles {
		generated, serialized := seedGeneratedSensorData(t, conn, gatewayID, p)
		expectedBySensor[generated.SensorId] = serialized
	}

	otherGateway := uuid.New()
	_, _ = seedGeneratedSensorData(t, conn, otherGateway, sensorprofiles.NewHeartRateProfile(uuid.New(), rand))

	data, err := repo.GetOrderedBufferedData(gatewayID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(data) != len(profiles) {
		t.Fatalf("expected %d rows, got %d", len(profiles), len(data))
	}

	for i := 1; i < len(data); i++ {
		if data[i-1].Timestamp.After(data[i].Timestamp) {
			t.Fatalf("expected ordered timestamps, got %s before %s", data[i].Timestamp, data[i-1].Timestamp)
		}
	}

	for _, d := range data {
		expected, ok := expectedBySensor[d.SensorId]
		if !ok {
			t.Fatalf("unexpected sensor id returned: %s", d.SensorId)
		}

		var expectedJSON any
		var actualJSON any
		if err := json.Unmarshal(expected, &expectedJSON); err != nil {
			t.Fatalf("expected json unmarshal failed: %v", err)
		}
		if err := json.Unmarshal(d.Data, &actualJSON); err != nil {
			t.Fatalf("actual json unmarshal failed: %v", err)
		}

		expectedCanonical, _ := json.Marshal(expectedJSON)
		actualCanonical, _ := json.Marshal(actualJSON)
		if !bytes.Equal(expectedCanonical, actualCanonical) {
			t.Fatalf("value column mismatch for sensor %s", d.SensorId)
		}
	}
}

func TestGetOrderedBufferedDataWrongColumnData(t *testing.T) {
	tests := []struct {
		name   string
		insert string
		argsFn func(gatewayID uuid.UUID) []any
	}{
		{
			name:   "wrong sensorId",
			insert: `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`,
			argsFn: func(gatewayID uuid.UUID) []any {
				return []any{gatewayID.String(), "not-a-uuid", time.Now().UTC(), "heart_rate", `{"BpmValue":70}`}
			},
		},
		{
			name:   "wrong timestamp",
			insert: `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`,
			argsFn: func(gatewayID uuid.UUID) []any {
				return []any{gatewayID.String(), uuid.New().String(), "not-a-time", "heart_rate", `{"BpmValue":70}`}
			},
		},
		{
			name:   "wrong value",
			insert: `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`,
			argsFn: func(gatewayID uuid.UUID) []any {
				return []any{gatewayID.String(), uuid.New().String(), time.Now().UTC(), "heart_rate", "{"}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, conn := newMockBufferRepository(t)
			gatewayID := uuid.New()

			if _, err := conn.ExecContext(context.Background(), tt.insert, tt.argsFn(gatewayID)...); err != nil {
				t.Fatalf("insert setup failed: %v", err)
			}

			_, err := repo.GetOrderedBufferedData(gatewayID)
			if err == nil {
				t.Fatal("expected error from wrong buffered data")
			}
		})
	}
}

func TestInsertWrongDataForEachColumn(t *testing.T) {
	_, conn := newMockBufferRepository(t)

	tests := []struct {
		name string
		args []any
	}{
		{name: "gatewayId null", args: []any{nil, uuid.New().String(), time.Now().UTC(), "heart_rate", `{"BpmValue":70}`}},
		{name: "sensorId null", args: []any{uuid.New().String(), nil, time.Now().UTC(), "heart_rate", `{"BpmValue":70}`}},
		{name: "timestamp null", args: []any{uuid.New().String(), uuid.New().String(), nil, "heart_rate", `{"BpmValue":70}`}},
		{name: "profile null", args: []any{uuid.New().String(), uuid.New().String(), time.Now().UTC(), nil, `{"BpmValue":70}`}},
		{name: "value null", args: []any{uuid.New().String(), uuid.New().String(), time.Now().UTC(), "heart_rate", nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conn.ExecContext(context.Background(), `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`, tt.args...)
			if err == nil {
				t.Fatalf("expected insert error for %s", tt.name)
			}
		})
	}
}

func TestGetOrderedBufferedDataWithWrongDBConnection(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := buffereddatasender.NewBufferedDataRepository(context.Background(), sensor.BufferDbConnection{DB: db})
	_, err = repo.GetOrderedBufferedData(uuid.New())
	if err == nil {
		t.Fatal("expected query error on wrong db schema")
	}
}

func TestCleanBufferedDataNilWhenEmptySlice(t *testing.T) {
	repo, _ := newMockBufferRepository(t)
	if err := repo.CleanBufferedData(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCleanBufferedDataDeletesOnlyPassedRows(t *testing.T) {
	repo, conn := newMockBufferRepository(t)
	gatewayID := uuid.New()

	t1 := time.Now().UTC().Add(-2 * time.Minute)
	t2 := time.Now().UTC().Add(-1 * time.Minute)
	t3 := time.Now().UTC()
	s1 := uuid.New()
	s2 := uuid.New()
	s3 := uuid.New()

	insert := `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`
	for _, row := range []struct {
		sensorID uuid.UUID
		ts       time.Time
	}{
		{sensorID: s1, ts: t1},
		{sensorID: s2, ts: t2},
		{sensorID: s3, ts: t3},
	} {
		if _, err := conn.ExecContext(context.Background(), insert, gatewayID.String(), row.sensorID.String(), row.ts, "heart_rate", `{"BpmValue":70}`); err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	data, err := repo.GetOrderedBufferedData(gatewayID)
	if err != nil {
		t.Fatalf("get ordered data failed: %v", err)
	}

	if err := repo.CleanBufferedData(data[:2]); err != nil {
		t.Fatalf("clean buffered data failed: %v", err)
	}

	remaining, err := repo.GetOrderedBufferedData(gatewayID)
	if err != nil {
		t.Fatalf("get remaining data failed: %v", err)
	}

	if len(remaining) != 1 {
		t.Fatalf("expected 1 remaining row, got %d", len(remaining))
	}
	if remaining[0].SensorId != s3 {
		t.Fatalf("expected remaining sensor %s, got %s", s3, remaining[0].SensorId)
	}
}

func TestCleanBufferedDataWithWrongDBConnection(t *testing.T) {
	goodRepo, goodConn := newMockBufferRepository(t)
	gatewayID := uuid.New()
	if _, err := goodConn.ExecContext(context.Background(), `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`, gatewayID.String(), uuid.New().String(), time.Now().UTC(), "heart_rate", `{"BpmValue":70}`); err != nil {
		t.Fatalf("failed to seed valid row: %v", err)
	}
	seedData, err := goodRepo.GetOrderedBufferedData(gatewayID)
	if err != nil {
		t.Fatalf("failed to load seed data: %v", err)
	}
	if len(seedData) == 0 {
		t.Fatal("expected seed data to be non-empty")
	}

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := buffereddatasender.NewBufferedDataRepository(context.Background(), sensor.BufferDbConnection{DB: db})
	err = repo.CleanBufferedData(seedData)
	if err == nil {
		t.Fatal("expected clean buffered data error on wrong db schema")
	}
}

func TestGetOrderedBufferedDataRespectsBatchLimitAndAllowsFullCleanup(t *testing.T) {
	repo, conn := newMockBufferRepository(t)
	gatewayID := uuid.New()

	const (
		rowsToInsert    = 700
		maxRowsPerBatch = 999 / 3
	)
	insert := `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`
	baseTs := time.Now().UTC().Add(-time.Hour)
	for i := 0; i < rowsToInsert; i++ {
		if _, err := conn.ExecContext(
			context.Background(),
			insert,
			gatewayID.String(),
			uuid.New().String(),
			baseTs.Add(time.Duration(i)*time.Millisecond),
			"heart_rate",
			`{"BpmValue":70}`,
		); err != nil {
			t.Fatalf("insert failed at row %d: %v", i, err)
		}
	}

	data, err := repo.GetOrderedBufferedData(gatewayID)
	if err != nil {
		t.Fatalf("get ordered data failed: %v", err)
	}
	if len(data) != maxRowsPerBatch {
		t.Fatalf("expected first batch of %d rows, got %d", maxRowsPerBatch, len(data))
	}

	cleaned := 0
	for {
		batch, err := repo.GetOrderedBufferedData(gatewayID)
		if err != nil {
			t.Fatalf("get ordered data failed during cleanup loop: %v", err)
		}
		if len(batch) == 0 {
			break
		}
		if len(batch) > maxRowsPerBatch {
			t.Fatalf("query returned %d rows, expected at most %d", len(batch), maxRowsPerBatch)
		}
		if err := repo.CleanBufferedData(batch); err != nil {
			t.Fatalf("clean buffered data failed: %v", err)
		}
		cleaned += len(batch)
	}

	if cleaned != rowsToInsert {
		t.Fatalf("expected cleaned rows %d, got %d", rowsToInsert, cleaned)
	}

	remaining, err := repo.GetOrderedBufferedData(gatewayID)
	if err != nil {
		t.Fatalf("get remaining rows failed: %v", err)
	}
	if len(remaining) != 0 {
		t.Fatalf("expected empty buffer after cleanup, got %d rows", len(remaining))
	}
}

func TestBufferTriggerLimitIsAppliedPerGateway(t *testing.T) {
	_, conn := newMockBufferRepository(t)
	gatewayA := uuid.New()
	gatewayB := uuid.New()

	insert := `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`
	baseTsA := time.Now().UTC().Add(-2 * time.Hour)
	baseTsB := time.Now().UTC().Add(-time.Hour)

	oldestTsA := baseTsA
	secondOldestTsA := baseTsA.Add(time.Millisecond)
	oldestTsB := baseTsB

	for i := 0; i < 1000; i++ {
		if _, err := conn.ExecContext(
			context.Background(),
			insert,
			gatewayA.String(),
			uuid.New().String(),
			baseTsA.Add(time.Duration(i)*time.Millisecond),
			"heart_rate",
			`{"BpmValue":70}`,
		); err != nil {
			t.Fatalf("insert first 1000 gateway A failed at row %d: %v", i, err)
		}
	}

	var countAFirst int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gatewayA.String()).Scan(&countAFirst); err != nil {
		t.Fatalf("count first 1000 gateway A failed: %v", err)
	}
	if countAFirst != 1000 {
		t.Fatalf("expected 1000 rows for gateway A before overflow, got %d", countAFirst)
	}

	var oldestAPresentBefore int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ? AND timestamp = ?`, gatewayA.String(), oldestTsA).Scan(&oldestAPresentBefore); err != nil {
		t.Fatalf("count oldest A before overflow failed: %v", err)
	}
	if oldestAPresentBefore != 1 {
		t.Fatalf("expected oldest gateway A row to exist before overflow, got %d", oldestAPresentBefore)
	}

	for i := 0; i < 600; i++ {
		if _, err := conn.ExecContext(
			context.Background(),
			insert,
			gatewayB.String(),
			uuid.New().String(),
			baseTsB.Add(time.Duration(i)*time.Millisecond),
			"heart_rate",
			`{"BpmValue":71}`,
		); err != nil {
			t.Fatalf("insert gateway B failed at row %d: %v", i, err)
		}
	}

	// Insert the 1001st row for gateway A: trigger must drop only its oldest row.
	if _, err := conn.ExecContext(
		context.Background(),
		insert,
		gatewayA.String(),
		uuid.New().String(),
		baseTsA.Add(1000*time.Millisecond),
		"heart_rate",
		`{"BpmValue":72}`,
	); err != nil {
		t.Fatalf("insert 1001st gateway A failed: %v", err)
	}

	var countA int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gatewayA.String()).Scan(&countA); err != nil {
		t.Fatalf("count gateway A failed: %v", err)
	}

	var countB int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`, gatewayB.String()).Scan(&countB); err != nil {
		t.Fatalf("count gateway B failed: %v", err)
	}

	var oldestAStillPresent int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ? AND timestamp = ?`, gatewayA.String(), oldestTsA).Scan(&oldestAStillPresent); err != nil {
		t.Fatalf("count oldest A after overflow failed: %v", err)
	}

	var secondOldestAPresent int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ? AND timestamp = ?`, gatewayA.String(), secondOldestTsA).Scan(&secondOldestAPresent); err != nil {
		t.Fatalf("count second oldest A after overflow failed: %v", err)
	}

	var oldestBPresent int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer WHERE gatewayId = ? AND timestamp = ?`, gatewayB.String(), oldestTsB).Scan(&oldestBPresent); err != nil {
		t.Fatalf("count oldest B after A overflow failed: %v", err)
	}

	var countTotal int
	if err := conn.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM buffer`).Scan(&countTotal); err != nil {
		t.Fatalf("count total failed: %v", err)
	}

	if countA != 1000 {
		t.Fatalf("expected 1000 rows for gateway A, got %d", countA)
	}
	if oldestAStillPresent != 0 {
		t.Fatalf("expected oldest gateway A row to be deleted at overflow, got %d", oldestAStillPresent)
	}
	if secondOldestAPresent != 1 {
		t.Fatalf("expected second oldest gateway A row to remain, got %d", secondOldestAPresent)
	}
	if countB != 600 {
		t.Fatalf("expected 600 rows for gateway B, got %d", countB)
	}
	if oldestBPresent != 1 {
		t.Fatalf("expected oldest gateway B row to remain untouched, got %d", oldestBPresent)
	}
	if countTotal != 1600 {
		t.Fatalf("expected 1600 total rows (>1000 globally), got %d", countTotal)
	}
}

func TestCleanWholeBufferForValidAndInvalidGatewayID(t *testing.T) {
	repo, conn := newMockBufferRepository(t)
	validGateway := uuid.New()
	otherGateway := uuid.New()
	invalidGateway := uuid.New()

	insert := `INSERT INTO buffer (gatewayId, sensorId, timestamp, profile, value) VALUES (?, ?, ?, ?, ?)`
	if _, err := conn.ExecContext(context.Background(), insert, validGateway.String(), uuid.New().String(), time.Now().UTC(), "heart_rate", `{"BpmValue":70}`); err != nil {
		t.Fatalf("insert valid gateway failed: %v", err)
	}
	if _, err := conn.ExecContext(context.Background(), insert, otherGateway.String(), uuid.New().String(), time.Now().UTC(), "heart_rate", `{"BpmValue":71}`); err != nil {
		t.Fatalf("insert other gateway failed: %v", err)
	}

	if err := repo.CleanWholeBuffer(validGateway); err != nil {
		t.Fatalf("clean whole buffer for valid gateway failed: %v", err)
	}

	rowsOther, err := repo.GetOrderedBufferedData(otherGateway)
	if err != nil {
		t.Fatalf("get other gateway rows failed: %v", err)
	}
	if len(rowsOther) != 1 {
		t.Fatalf("expected other gateway row to remain, got %d", len(rowsOther))
	}

	if err := repo.CleanWholeBuffer(invalidGateway); err != nil {
		t.Fatalf("clean whole buffer for invalid gateway should not fail, got %v", err)
	}
}

func TestCleanWholeBufferWithWrongDBConnection(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := buffereddatasender.NewBufferedDataRepository(context.Background(), sensor.BufferDbConnection{DB: db})
	err = repo.CleanWholeBuffer(uuid.New())
	if err == nil {
		t.Fatal("expected clean whole buffer error on wrong db schema")
	}
}
