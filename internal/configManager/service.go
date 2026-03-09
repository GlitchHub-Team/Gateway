package configmanager

import (
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensor "Gateway/internal/sensor"

	"github.com/google/uuid"
)

type ConfigManagerService struct {
	configPort ConfigPort
}

func NewConfigManagerService(configPort ConfigPort) *ConfigManagerService {
	return &ConfigManagerService{
		configPort: configPort,
	}
}

func (s *ConfigManagerService) GetAllGateways() (map[uuid.UUID]*Gateway, error) {
	return s.configPort.GetAllGateways()
}

func (s *ConfigManagerService) GetAllGatewaysByTenantId(tenantId uuid.UUID) (map[uuid.UUID]Gateway, error) {
	return s.configPort.GetAllGatewaysByTenantId(tenantId)
}

func (s *ConfigManagerService) GetGatewayById(gatewayId uuid.UUID) (*Gateway, error) {
	return s.configPort.GetGatewayById(gatewayId)
}

func (s *ConfigManagerService) GetSensorById(gatewayId uuid.UUID, sensorId uuid.UUID) (*sensor.Sensor, error) {
	return s.configPort.GetSensorById(gatewayId, sensorId)
}

func (s *ConfigManagerService) ChangeSensorFrequency(cmdData *commanddata.ChangeSensorFrequency) error {
	return s.configPort.ChangeSensorFrequency(cmdData)
}

func (s *ConfigManagerService) CommissionGateway(cmdData *commanddata.CommissionGateway) error {
	return s.configPort.CommissionGateway(cmdData)
}

func (s *ConfigManagerService) CreateGateway(cmdData *commanddata.CreateGateway) error {
	return s.configPort.CreateGateway(cmdData)
}

func (s *ConfigManagerService) DecommissionGateway(cmdData *commanddata.DecommissionGateway) error {
	return s.configPort.DecommissionGateway(cmdData)
}

func (s *ConfigManagerService) DeleteGateway(cmdData *commanddata.DeleteGateway) error {
	return s.configPort.DeleteGateway(cmdData)
}

func (s *ConfigManagerService) InterruptGateway(cmdData *commanddata.InterruptGateway) error {
	return s.configPort.InterruptGateway(cmdData)
}

func (s *ConfigManagerService) RebootGateway(cmdData *commanddata.RebootGateway) error {
	return s.configPort.RebootGateway(cmdData)
}

func (s *ConfigManagerService) ResetGateway(cmdData *commanddata.ResetGateway) error {
	return s.configPort.ResetGateway(cmdData)
}

func (s *ConfigManagerService) ResumeGateway(cmdData *commanddata.ResumeGateway) error {
	return s.configPort.ResumeGateway(cmdData)
}

func (s *ConfigManagerService) InterruptSensor(cmdData *commanddata.InterruptSensor) error {
	return s.configPort.InterruptSensor(cmdData)
}

func (s *ConfigManagerService) ResumeSensor(cmdData *commanddata.ResumeSensor) error {
	return s.configPort.ResumeSensor(cmdData)
}

func (s *ConfigManagerService) AddSensor(cmdData *commanddata.AddSensor) error {
	return s.configPort.AddSensor(cmdData)
}

func (s *ConfigManagerService) DeleteSensor(cmdData *commanddata.DeleteSensor) error {
	return s.configPort.DeleteSensor(cmdData)
}
