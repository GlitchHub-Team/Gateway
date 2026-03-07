package gatewaydatabase

import (
	"context"
	"database/sql"
	"log"

	configManager "Gateway/internal/configManager"

	_ "modernc.org/sqlite"
)

func createGatewayTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS gateways (
		id VARCHAR(255) PRIMARY KEY,
		tenantId VARCHAR(255) NOT NULL,
		status VARCHAR(255) NOT NULL
	);
	`
	_, err := db.ExecContext(context.Background(), query)
	return err
}

func createSensorTypeFrequencyTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sensor_type_frequencies (
		gatewayId VARCHAR(255) NOT NULL,
		sensorType VARCHAR(255) NOT NULL,
		frequency INT NOT NULL,
		PRIMARY KEY (gatewayId, sensorType),
		FOREIGN KEY (gatewayId) REFERENCES gateways(id)
	);
	`
	_, err := db.ExecContext(context.Background(), query)
	return err
}

func createSensorTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sensors (
		id VARCHAR(255),
		gatewayId VARCHAR(255) NOT NULL,
		profile VARCHAR(255) NOT NULL,
		status VARCHAR(255) NOT NULL,
		PRIMARY KEY (id, gatewayId),
		FOREIGN KEY (gatewayId) REFERENCES gateways(id)
	);
	`
	_, err := db.ExecContext(context.Background(), query)
	return err
}

func NewGatewayDatabase() configManager.ConfigDbConnection {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Error while opening DB: %v", err)
	}

	err = createGatewayTable(db)
	if err != nil {
		log.Fatalf("Error while creating gateways table: %v", err)
	}

	err = createSensorTable(db)
	if err != nil {
		log.Fatalf("Error while creating sensors table: %v", err)
	}

	err = createSensorTypeFrequencyTable(db)
	if err != nil {
		log.Fatalf("Error while creating sensor type frequencies table: %v", err)
	}

	return configManager.ConfigDbConnection{DB: db}
}
