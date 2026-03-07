package modules

import (
	configManager "Gateway/internal/configManager"

	"go.uber.org/fx"
)

var ConfigManagerModule = fx.Module("ConfigManagerModule",
	fx.Provide(
		fx.Annotate(
			configManager.NewConfigManagerService,
			fx.As(new(configManager.SensorFrequencySetter)),
			fx.As(new(configManager.GatewayCommissioner)),
			fx.As(new(configManager.GatewayCreator)),
			fx.As(new(configManager.GatewayDecommissioner)),
			fx.As(new(configManager.GatewayDeleter)),
			fx.As(new(configManager.GatewayInterrupter)),
			fx.As(new(configManager.GatewayRebooter)),
			fx.As(new(configManager.GatewayResetter)),
			fx.As(new(configManager.GatewayResumer)),
			fx.As(new(configManager.SensorInterrupter)),
			fx.As(new(configManager.SensorResumer)),
			fx.As(new(configManager.SensorAdder)),
			fx.As(new(configManager.SensorDeleter)),
		),
		fx.Annotate(
			configManager.NewSQLiteConfigRepository,
			fx.As(new(configManager.ConfigPort)),
		),
	),
)
