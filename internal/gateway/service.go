package gateway

type GatewayManagerService struct{}

func NewGatewayManagerService() *GatewayManagerService {
	return &GatewayManagerService{}
}

func (s *GatewayManagerService) CreateGateway(dto CreateGateway) error {
	// Logic to create a gateway
	return nil
}

func (s *GatewayManagerService) CommissionGateway(dto CommissionGateway) error {
	// Logic to commission a gateway
	return nil
}

func (s *GatewayManagerService) DecommissionGateway(dto DecommissionGateway) error {
	// Logic to decommission a gateway
	return nil
}

func (s *GatewayManagerService) DeleteGateway(dto DeleteGateway) error {
	// Logic to delete a gateway
	return nil
}

func (s *GatewayManagerService) InterruptGateway(dto InterruptGateway) error {
	// Logic to interrupt a gateway
	return nil
}

func (s *GatewayManagerService) RebootGateway(dto RebootGateway) error {
	// Logic to reboot a gateway
	return nil
}

func (s *GatewayManagerService) ResetGateway(dto ResetGateway) error {
	// Logic to reset a gateway
	return nil
}

func (s *GatewayManagerService) ResumeGateway(dto ResumeGateway) error {
	// Logic to resume a gateway
	return nil
}

func (s *GatewayManagerService) AddSensor(dto AddSensor) error {
	// Logic to add a sensor to a gateway
	return nil
}

func (s *GatewayManagerService) DeleteSensor(dto DeleteSensor) error {
	// Logic to delete a sensor from a gateway
	return nil
}

func (s *GatewayManagerService) InterruptSensor(dto InterruptSensor) error {
	// Logic to interrupt a sensor
	return nil
}

func (s *GatewayManagerService) ResumeSensor(dto ResumeSensor) error {
	// Logic to resume a sensor
	return nil
}

func (s *GatewayManagerService) ChangeSensorFrequency(dto ChangeSensorFrequency) error {
	// Logic to change the frequency of a sensor
	return nil
}
