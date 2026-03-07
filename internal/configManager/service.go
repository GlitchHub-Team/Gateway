package configmanager

import commanddata "Gateway/internal/gateway/commandData"

type ConfigManagerService struct {
	configPort *ConfigPort
}

func NewConfigManagerService(configPort *ConfigPort) *ConfigManagerService {
	return &ConfigManagerService{
		configPort: configPort,
	}
}

// TODO: separate implementations into different files for better organization
func (s *ConfigManagerService) CreateGateway(cmdData *commanddata.CreateGateway) error {
	// TODO: implement gateway creation logic
	return nil
}

func (s *ConfigManagerService) CommissionGateway(cmdData *commanddata.CommissionGateway) error {
	// TODO: implement gateway commissioning logic
	return nil
}

func (s *ConfigManagerService) DecommissionGateway(cmdData *commanddata.DecommissionGateway) error {
	// TODO: implement gateway decommissioning logic
	return nil
}

func (s *ConfigManagerService) DeleteGateway(cmdData *commanddata.DeleteGateway) error {
	// TODO: implement gateway deletion logic
	return nil
}

func (s *ConfigManagerService) InterruptGateway(cmdData *commanddata.InterruptGateway) error {
	// TODO: implement gateway interruption logic
	return nil
}

func (s *ConfigManagerService) RebootGateway(cmdData *commanddata.RebootGateway) error {
	// TODO: implement gateway reboot logic
	return nil
}

func (s *ConfigManagerService) ResetGateway(cmdData *commanddata.ResetGateway) error {
	// TODO: implement gateway reset logic
	return nil
}

func (s *ConfigManagerService) ResumeGateway(cmdData *commanddata.ResumeGateway) error {
	// TODO: implement gateway resume logic
	return nil
}

func (s *ConfigManagerService) AddSensor(cmdData *commanddata.AddSensor) error {
	// TODO: implement sensor addition logic
	return nil
}

func (s *ConfigManagerService) DeleteSensor(cmdData *commanddata.DeleteSensor) error {
	// TODO: implement sensor deletion logic
	return nil
}

func (s *ConfigManagerService) InterruptSensor(cmdData *commanddata.InterruptSensor) error {
	// TODO: implement sensor interruption logic
	return nil
}

func (s *ConfigManagerService) ResumeSensor(cmdData *commanddata.ResumeSensor) error {
	// TODO: implement sensor resume logic
	return nil
}

func (s *ConfigManagerService) ChangeSensorFrequency(cmdData *commanddata.ChangeSensorFrequency) error {
	// TODO: implement sensor frequency change logic
	return nil
}
