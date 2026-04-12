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

const (
	sqliteMaxVariables = 999
	deleteTupleSize    = 3
	maxRowsPerBatch    = sqliteMaxVariables / deleteTupleSize
)

func NewBufferedDataRepository(ctx context.Context, conn sensor.BufferDbConnection) *BufferedDataRepository {
	return &BufferedDataRepository{
		ctx:          ctx,
		dbConnection: conn,
	}
}

func (b *BufferedDataRepository) GetOrderedBufferedData(gatewayId uuid.UUID) ([]*sensorData, error) {
	query := `SELECT sensorId, gatewayId, timestamp, profile, json(value)
				FROM buffer 
				WHERE gatewayId = ? 
				ORDER BY timestamp ASC
				LIMIT ?`
	rows, err := b.dbConnection.QueryContext(b.ctx, query, gatewayId, maxRowsPerBatch)
	if err != nil {
		return nil, fmt.Errorf("errore nell'eseguire la query per ottenere i dati del buffer: %w, gatewayId: %s", err, gatewayId.String())
	}

	var data []*sensorData
	for rows.Next() {
		var sensorId, gatewayId uuid.UUID
		var timestamp time.Time
		var profile string
		var value []byte
		if err := rows.Scan(&sensorId, &gatewayId, &timestamp, &profile, &value); err != nil {
			return nil, fmt.Errorf("errore nello scan della riga del buffer: %w, gatewayId: %s", err, gatewayId.String())
		}
		data = append(data, &sensorData{
			SensorId:  sensorId,
			GatewayId: gatewayId,
			Timestamp: timestamp,
			Profile:   profile,
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

	for start := 0; start < len(data); start += maxRowsPerBatch {
		end := start + maxRowsPerBatch
		if end > len(data) {
			end = len(data)
		}

		chunk := data[start:end]
		placeholders := slices.Repeat([]string{"(?, ?, ?)"}, len(chunk))
		query := `DELETE FROM buffer WHERE (gatewayId, sensorId, timestamp) IN (%s)`
		generatedQuery := fmt.Sprintf(query, strings.Join(placeholders, ", "))

		args := make([]any, 0, len(chunk)*deleteTupleSize)
		for _, d := range chunk {
			args = append(args, d.GatewayId, d.SensorId, d.Timestamp)
		}

		if _, err := b.dbConnection.ExecContext(b.ctx, generatedQuery, args...); err != nil {
			return fmt.Errorf("errore nell'eseguire la query per pulire i dati del buffer: %w", err)
		}
	}

	return nil
}

func (b *BufferedDataRepository) CleanWholeBuffer(gatewayId uuid.UUID) error {
	query := `DELETE FROM buffer WHERE gatewayId = ?`
	_, err := b.dbConnection.ExecContext(b.ctx, query, gatewayId)
	if err != nil {
		return fmt.Errorf("errore nell'eseguire la query per pulire il buffer del gateway %s: %w", gatewayId.String(), err)
	}
	return nil
}
