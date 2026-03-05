package gateway

type CreateGatewayUseCase interface {
	CreateGateway(dto CreateGateway) error
}

type CommissionGatewayUseCase interface {
	CommissionGateway(dto CommissionGateway) error
}

type DecommissionGatewayUseCase interface {
	DecommissionGateway(dto DecommissionGateway) error
}

type DeleteGatewayUseCase interface {
	DeleteGateway(dto DeleteGateway) error
}

type InterruptGatewayUseCase interface {
	InterruptGateway(dto InterruptGateway) error
}

type RebootGatewayUseCase interface {
	RebootGateway(dto RebootGateway) error
}

type ResetGatewayUseCase interface {
	ResetGateway(dto ResetGateway) error
}

type ResumeGatewayUseCase interface {
	ResumeGateway(dto ResumeGateway) error
}

type AddSensorUseCase interface {
	AddSensor(dto AddSensor) error
}

type DeleteSensorUseCase interface {
	DeleteSensor(dto DeleteSensor) error
}

type InterruptSensorUseCase interface {
	InterruptSensor(dto InterruptSensor) error
}

type ResumeSensorUseCase interface {
	ResumeSensor(dto ResumeSensor) error
}

type ChangeSensorFrequencyUseCase interface {
	ChangeSensorFrequency(dto ChangeSensorFrequency) error
}
