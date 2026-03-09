package modules

import (
	gatmanager "Gateway/internal/gatewayManager"
	gatewayusecases "Gateway/internal/gatewayManager/gatewayUseCases"

	"go.uber.org/fx"
)

var GatewayModule = fx.Module("GatewayModule",
	fx.Provide(
		gatmanager.NewSensorWorkers,
		gatmanager.NewGatewayWorkers,
		fx.Annotate(
			gatmanager.NewGatewayManagerService,
			fx.As(new(gatmanager.GatewaysLoader)),
			fx.As(new(gatewayusecases.CreateGatewayUseCase)),
			fx.As(new(gatewayusecases.DeleteGatewayUseCase)),
			fx.As(new(gatewayusecases.CommissionGatewayUseCase)),
			fx.As(new(gatewayusecases.DecommissionGatewayUseCase)),
			fx.As(new(gatewayusecases.InterruptGatewayUseCase)),
			fx.As(new(gatewayusecases.ResumeGatewayUseCase)),
			fx.As(new(gatewayusecases.ResetGatewayUseCase)),
			fx.As(new(gatewayusecases.RebootGatewayUseCase)),
			fx.As(new(gatewayusecases.AddSensorUseCase)),
			fx.As(new(gatewayusecases.DeleteSensorUseCase)),
			fx.As(new(gatewayusecases.InterruptSensorUseCase)),
			fx.As(new(gatewayusecases.ResumeSensorUseCase)),
			fx.As(new(gatewayusecases.ChangeSensorFrequencyUseCase)),
		),
	),
)
