package bufferdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	sensor "Gateway/internal/sensor"

	_ "modernc.org/sqlite"
)

const (
	MAX_BUFFER_SIZE = 1000
)

func createBufferTable(db *sql.DB) error {
	query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS buffer (
        gatewayId TEXT NOT NULL,
        sensorId TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        profile TEXT NOT NULL,
        value BLOB NOT NULL,
        PRIMARY KEY (gatewayId, sensorId, timestamp)
    );

	CREATE INDEX IF NOT EXISTS idx_buffer_cleanup 
	ON buffer (gatewayId, timestamp DESC);

	DROP TRIGGER IF EXISTS check_buffer_limit;

	CREATE TRIGGER check_buffer_limit
    AFTER INSERT ON buffer
    BEGIN
        DELETE FROM buffer
		WHERE rowid IN (
			  SELECT rowid
              FROM buffer
              WHERE gatewayId = NEW.gatewayId
              ORDER BY timestamp DESC
			  LIMIT -1 OFFSET %d
          );
    END;
    `, MAX_BUFFER_SIZE)

	_, err := db.ExecContext(context.Background(), query)
	return err
}

func NewBufferDatabase() sensor.BufferDbConnection {
	dsn := "file:buffer.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("Error while opening DB: %v", err)
	}

	err = createBufferTable(db)
	if err != nil {
		log.Fatalf("Error while creating buffer table: %v", err)
	}
	log.Println("Buffer database initialized successfully")

	return sensor.BufferDbConnection{DB: db}
}

func NewMockBufferDatabase() sensor.BufferDbConnection {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Error while opening DB: %v", err)
	}

	err = createBufferTable(db)
	if err != nil {
		log.Fatalf("Error while creating buffer table: %v", err)
	}
	log.Println("Buffer database initialized successfully")

	return sensor.BufferDbConnection{DB: db}
}
