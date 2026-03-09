package gatewayusecases

import (
	gatManager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
)

type CreateGatewayUseCase interface {
	CreateGateway(cmdData *commanddata.CreateGateway) gatManager.Response
}

type CommissionGatewayUseCase interface {
	CommissionGateway(cmdData *commanddata.CommissionGateway) gatManager.Response
}

type DecommissionGatewayUseCase interface {
	DecommissionGateway(cmdData *commanddata.DecommissionGateway) gatManager.Response
}

type DeleteGatewayUseCase interface {
	DeleteGateway(cmdData *commanddata.DeleteGateway) gatManager.Response
}

type InterruptGatewayUseCase interface {
	InterruptGateway(cmdData *commanddata.InterruptGateway) gatManager.Response
}

type RebootGatewayUseCase interface {
	RebootGateway(cmdData *commanddata.RebootGateway) gatManager.Response
}

type ResetGatewayUseCase interface {
	ResetGateway(cmdData *commanddata.ResetGateway) gatManager.Response
}

type ResumeGatewayUseCase interface {
	ResumeGateway(cmdData *commanddata.ResumeGateway) gatManager.Response
}

type AddSensorUseCase interface {
	AddSensor(cmdData *commanddata.AddSensor) gatManager.Response
}

type DeleteSensorUseCase interface {
	DeleteSensor(cmdData *commanddata.DeleteSensor) gatManager.Response
}

type InterruptSensorUseCase interface {
	InterruptSensor(cmdData *commanddata.InterruptSensor) gatManager.Response
}

type ResumeSensorUseCase interface {
	ResumeSensor(cmdData *commanddata.ResumeSensor) gatManager.Response
}

type ChangeSensorFrequencyUseCase interface {
	ChangeSensorFrequency(cmdData *commanddata.ChangeSensorFrequency) gatManager.Response
}
