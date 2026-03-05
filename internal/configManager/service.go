package configmanager

import (
	gateway "Gateway/internal/gateway"
)

type ConfigManagerService struct{}

func NewConfigManagerService() *ConfigManagerService {
	return &ConfigManagerService{}
}

// TODO: separate implementations into different files for better organization
func (s *ConfigManagerService) CreateGateway(dto gateway.CreateGateway) error {
	// TODO: implement gateway creation logic
	return nil
}

func (s *ConfigManagerService) CommissionGateway(dto gateway.CommissionGateway) error {
	// TODO: implement gateway commissioning logic
	return nil
}

func (s *ConfigManagerService) DecommissionGateway(dto gateway.DecommissionGateway) error {
	// TODO: implement gateway decommissioning logic
	return nil
}

func (s *ConfigManagerService) DeleteGateway(dto gateway.DeleteGateway) error {
	// TODO: implement gateway deletion logic
	return nil
}

func (s *ConfigManagerService) InterruptGateway(dto gateway.InterruptGateway) error {
	// TODO: implement gateway interruption logic
	return nil
}

func (s *ConfigManagerService) RebootGateway(dto gateway.RebootGateway) error {
	// TODO: implement gateway reboot logic
	return nil
}

func (s *ConfigManagerService) ResetGateway(dto gateway.ResetGateway) error {
	// TODO: implement gateway reset logic
	return nil
}

func (s *ConfigManagerService) ResumeGateway(dto gateway.ResumeGateway) error {
	// TODO: implement gateway resume logic
	return nil
}

func (s *ConfigManagerService) AddSensor(dto gateway.AddSensor) error {
	// TODO: implement sensor addition logic
	return nil
}

func (s *ConfigManagerService) DeleteSensor(dto gateway.DeleteSensor) error {
	// TODO: implement sensor deletion logic
	return nil
}

func (s *ConfigManagerService) InterruptSensor(dto gateway.InterruptSensor) error {
	// TODO: implement sensor interruption logic
	return nil
}

func (s *ConfigManagerService) ResumeSensor(dto gateway.ResumeSensor) error {
	// TODO: implement sensor resume logic
	return nil
}

func (s *ConfigManagerService) ChangeSensorFrequency(dto gateway.ChangeSensorFrequency) error {
	// TODO: implement sensor frequency change logic
	return nil
}
