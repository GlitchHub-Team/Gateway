package buffereddatasender

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	sensor "Gateway/internal/sensor"

	"github.com/google/uuid"
)

type BufferedDataRepository struct {
	ctx          context.Context
	dbConnection sensor.BufferDbConnection
}

func NewBufferedDataRepository(ctx context.Context, conn sensor.BufferDbConnection) *BufferedDataRepository {
	return &BufferedDataRepository{
		ctx:          ctx,
		dbConnection: conn,
	}
}

func (b *BufferedDataRepository) GetOrderedBufferedData(gatewayId uuid.UUID) ([]*sensorData, error) {
	query := `SELECT sensor_id, gateway_id, timestamp, value
				FROM buffered_data 
				WHERE gateway_id = ? 
				ORDER BY timestamp ASC`
	rows, err := b.dbConnection.QueryContext(b.ctx, query, gatewayId)
	if err != nil {
		return nil, fmt.Errorf("errore nell'eseguire la query per ottenere i dati del buffer: %w, gatewayId: %s", err, gatewayId.String())
	}

	var data []*sensorData
	for rows.Next() {
		var sensorId, gatewayId uuid.UUID
		var timestamp time.Time
		var value []byte
		if err := rows.Scan(&sensorId, &gatewayId, &timestamp, &value); err != nil {
			return nil, fmt.Errorf("errore nello scan della riga del buffer: %w, gatewayId: %s", err, gatewayId.String())
		}
		data = append(data, &sensorData{
			SensorId:  sensorId,
			GatewayId: gatewayId,
			Timestamp: timestamp,
			Data:      value,
		})
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("errore nella chiusura delle righe del buffer: %w, gatewayId: %s", err, gatewayId.String())
	}

	return data, nil
}

func (b *BufferedDataRepository) CleanBufferedData(data []*sensorData) error {
	if len(data) <= 0 {
		return nil
	}

	placeholders := slices.Repeat([]string{"(?, ?, ?)"}, len(data))
	query := `DELETE FROM buffered_data WHERE (gatewayId, sensorId, timestamp) IN (%s)`
	generatedQuery := fmt.Sprintf(query, strings.Join(placeholders, ", "))

	args := make([]any, 0, len(data)*3)
	for _, d := range data {
		args = append(args, d.GatewayId, d.SensorId, d.Timestamp)
	}

	_, err := b.dbConnection.ExecContext(b.ctx, generatedQuery, args...)
	if err != nil {
		return fmt.Errorf("errore nell'eseguire la query per pulire i dati del buffer: %w", err)
	}

	return nil
}

func (b *BufferedDataRepository) CleanWholeBuffer(gatewayId uuid.UUID) error {
	query := `DELETE FROM buffered WHERE gatewayId = ?`
	_, err := b.dbConnection.ExecContext(b.ctx, query, gatewayId)
	if err != nil {
		return fmt.Errorf("errore nell'eseguire la query per pulire il buffer del gateway %s: %w", gatewayId.String(), err)
	}
	return nil
}
