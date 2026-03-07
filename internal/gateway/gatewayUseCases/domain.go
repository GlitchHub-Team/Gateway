package gatewayusecases

import (
	gateway "Gateway/internal/gateway"
	commanddata "Gateway/internal/gateway/commandData"
)

type CreateGatewayUseCase interface {
	CreateGateway(cmdData *commanddata.CreateGateway) *gateway.Response
}

type CommissionGatewayUseCase interface {
	CommissionGateway(cmdData *commanddata.CommissionGateway) *gateway.Response
}

type DecommissionGatewayUseCase interface {
	DecommissionGateway(cmdData *commanddata.DecommissionGateway) *gateway.Response
}

type DeleteGatewayUseCase interface {
	DeleteGateway(cmdData *commanddata.DeleteGateway) *gateway.Response
}

type InterruptGatewayUseCase interface {
	InterruptGateway(cmdData *commanddata.InterruptGateway) *gateway.Response
}

type RebootGatewayUseCase interface {
	RebootGateway(cmdData *commanddata.RebootGateway) *gateway.Response
}

type ResetGatewayUseCase interface {
	ResetGateway(cmdData *commanddata.ResetGateway) *gateway.Response
}

type ResumeGatewayUseCase interface {
	ResumeGateway(cmdData *commanddata.ResumeGateway) *gateway.Response
}

type AddSensorUseCase interface {
	AddSensor(cmdData *commanddata.AddSensor) *gateway.Response
}

type DeleteSensorUseCase interface {
	DeleteSensor(cmdData *commanddata.DeleteSensor) *gateway.Response
}

type InterruptSensorUseCase interface {
	InterruptSensor(cmdData *commanddata.InterruptSensor) *gateway.Response
}

type ResumeSensorUseCase interface {
	ResumeSensor(cmdData *commanddata.ResumeSensor) *gateway.Response
}

type ChangeSensorFrequencyUseCase interface {
	ChangeSensorFrequency(cmdData *commanddata.ChangeSensorFrequency) *gateway.Response
}
