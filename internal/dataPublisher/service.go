package datapublisher

import (
	gateway "Gateway/internal/gateway"
)

type DataPublisherService struct {
	status       gateway.GatewayStatus
	sendDataRepo SendSensorDataRepository
}

func NewDataPublisherService(status gateway.GatewayStatus, sendDataRepo SendSensorDataRepository) *DataPublisherService {
	return &DataPublisherService{
		status:       status,
		sendDataRepo: sendDataRepo,
	}
}

func (s *DataPublisherService) Start() error {
	// Logic to start the data publisher service
	return nil
}

func (s *DataPublisherService) Stop() error {
	// Logic to stop the data publisher service
	return nil
}
