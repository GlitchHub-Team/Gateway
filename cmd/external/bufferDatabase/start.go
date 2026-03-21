package bufferdatabase

import (
	"context"
	"database/sql"
	"log"

	sensor "Gateway/internal/sensor"

	_ "modernc.org/sqlite"
)

func createBufferTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS buffer (
		gatewayId VARCHAR(255) NOT NULL,
		sensorId VARCHAR(255) NOT NULL,
		timestamp DATETIME NOT NULL,
		profile VARCHAR(255) NOT NULL,
		value JSONB NOT NULL,
		PRIMARY KEY (gatewayId, sensorId, timestamp)
	);
	`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

func NewBufferDatabase() sensor.BufferDbConnection {
	db, err := sql.Open("sqlite", "file:buffer.db")
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
