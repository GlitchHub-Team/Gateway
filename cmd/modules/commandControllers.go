package modules

import (
	commandControllers "Gateway/internal/gateway/commandControllers"

	"go.uber.org/fx"
)

var CommandControllersModule = fx.Module("CommandControllersModule",
	fx.Provide(
		commandControllers.NewNATSAddSensorController,
		commandControllers.NewNATSChangeSensorFrequencyController,
		commandControllers.NewNATSCommissionGatewayController,
		commandControllers.NewNATSDecommissionGatewayController,
		commandControllers.NewNATSDeleteGatewayController,
		commandControllers.NewNATSCreateGatewayController,
		commandControllers.NewNATSDeleteSensorController,
		commandControllers.NewNATSInterruptGatewayController,
		commandControllers.NewNATSInterruptSensorController,
		commandControllers.NewNATSResumeSensorController,
		commandControllers.NewNATSResetGatewayController,
		commandControllers.NewNATSResumeGatewayController,
		commandControllers.NewNATSRebootGatewayController,
	),
)
