package buffereddatasender

import (
	configmanager "Gateway/internal/configManager"
)

type BufferedDataSenderService struct {
	status       configmanager.GatewayStatus
	sendDataRepo SendSensorDataRepository
}

func NewBufferedDataSenderService(status configmanager.GatewayStatus, sendDataRepo SendSensorDataRepository) *BufferedDataSenderService {
	return &BufferedDataSenderService{
		status:       status,
		sendDataRepo: sendDataRepo,
	}
}

func (s *BufferedDataSenderService) Start() error {
	// Logic to start the data publisher service
	return nil
}

func (s *BufferedDataSenderService) Stop() error {
	// Logic to stop the data publisher service
	return nil
}
