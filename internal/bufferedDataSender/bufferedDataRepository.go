package buffereddatasender

import sensor "Gateway/internal/sensor"

type BufferedDataRepository struct {
	dbConnection sensor.BufferDbConnection
}

func NewBufferedDataRepository(conn sensor.BufferDbConnection) *BufferedDataRepository {
	return &BufferedDataRepository{
		dbConnection: conn,
	}
}

func (b *BufferedDataRepository) GetOrderedBufferedData() ([]*sensor.SensorData, error) {
	// Implementation for fetching buffered data
	return nil, nil
}
