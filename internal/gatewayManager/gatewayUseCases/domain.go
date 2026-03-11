package gatewayusecases

import (
	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayservices "Gateway/internal/gatewayManager/gatewayServices"
)

type CreateGatewayUseCase interface {
	CreateGateway(cmdData *commanddata.CreateGateway) gatewayservices.Response
}

type CommissionGatewayUseCase interface {
	CommissionGateway(cmdData *commanddata.CommissionGateway) gatewayservices.Response
}

type DecommissionGatewayUseCase interface {
	DecommissionGateway(cmdData *commanddata.DecommissionGateway) gatewayservices.Response
}

type DeleteGatewayUseCase interface {
	DeleteGateway(cmdData *commanddata.DeleteGateway) gatewayservices.Response
}

type InterruptGatewayUseCase interface {
	InterruptGateway(cmdData *commanddata.InterruptGateway) gatewayservices.Response
}

type RebootGatewayUseCase interface {
	RebootGateway(cmdData *commanddata.RebootGateway) gatewayservices.Response
}

type ResetGatewayUseCase interface {
	ResetGateway(cmdData *commanddata.ResetGateway) gatewayservices.Response
}

type ResumeGatewayUseCase interface {
	ResumeGateway(cmdData *commanddata.ResumeGateway) gatewayservices.Response
}

type AddSensorUseCase interface {
	AddSensor(cmdData *commanddata.AddSensor) gatewayservices.Response
}

type DeleteSensorUseCase interface {
	DeleteSensor(cmdData *commanddata.DeleteSensor) gatewayservices.Response
}

type InterruptSensorUseCase interface {
	InterruptSensor(cmdData *commanddata.InterruptSensor) gatewayservices.Response
}

type ResumeSensorUseCase interface {
	ResumeSensor(cmdData *commanddata.ResumeSensor) gatewayservices.Response
}
