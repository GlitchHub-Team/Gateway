package gatewaydatabase

import (
	"context"
	"database/sql"
	"log"

	configrepositories "Gateway/internal/configManager/configRepositories"

	_ "modernc.org/sqlite"
)

func createGatewayTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS gateways (
		id VARCHAR(255) PRIMARY KEY,
		tenantId VARCHAR(255),
		status VARCHAR(255) NOT NULL,
		interval INT NOT NULL,
		publicIdentifier VARCHAR(255) NOT NULL,
		secretKey VARCHAR(255) NOT NULL,
		token VARCHAR(255)
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
		interval INT NOT NULL,
		PRIMARY KEY (id, gatewayId),
		FOREIGN KEY (gatewayId) REFERENCES gateways(id)
	);
	`
	_, err := db.ExecContext(context.Background(), query)
	return err
}

func NewGatewayDatabase() configrepositories.ConfigDbConnection {
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

	return configrepositories.ConfigDbConnection{DB: db}
}

/*
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
*/
