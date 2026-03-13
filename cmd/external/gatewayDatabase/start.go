package gatewaydatabase

import (
	"context"
	"database/sql"

	configrepositories "Gateway/internal/configManager/configRepositories"

	_ "modernc.org/sqlite"
)

func createGatewayTable(db *sql.DB, ctx context.Context) error {
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
	_, err := db.ExecContext(ctx, query)
	return err
}

func createSensorTable(db *sql.DB, ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS sensors (
		id VARCHAR(255),
		gatewayId VARCHAR(255) NOT NULL,
		profile VARCHAR(255) NOT NULL,
		status VARCHAR(255) NOT NULL,
		interval INT NOT NULL,
		PRIMARY KEY (id, gatewayId),
		FOREIGN KEY (gatewayId) REFERENCES gateways(id) ON DELETE CASCADE
	);
	`
	_, err := db.ExecContext(ctx, query)
	return err
}

func NewGatewayDatabase(ctx context.Context) (*configrepositories.ConfigDbConnection, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, "PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	err = createGatewayTable(db, ctx)
	if err != nil {
		return nil, err
	}

	err = createSensorTable(db, ctx)
	if err != nil {
		return nil, err
	}

	return &configrepositories.ConfigDbConnection{DB: db}, nil
}

/*
func createSensorTypeFrequencyTable(db *sql.DB, ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS sensor_type_frequencies (
		gatewayId VARCHAR(255) NOT NULL,
		sensorType VARCHAR(255) NOT NULL,
		frequency INT NOT NULL,
		PRIMARY KEY (gatewayId, sensorType),
		FOREIGN KEY (gatewayId) REFERENCES gateways(id) ON DELETE CASCADE
	);
	`
	_, err := db.ExecContext(ctx, query)
	return err
}
*/
