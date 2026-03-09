package modules

import (
	configmanager "Gateway/internal/configManager"
	configrepositories "Gateway/internal/configManager/configRepositories"

	"go.uber.org/fx"
)

var ConfigManagerModule = fx.Module("ConfigManagerModule",
	fx.Provide(
		fx.Annotate(
			configmanager.NewConfigManagerService,
			fx.As(new(configmanager.GatewaysFetcher)),
			fx.As(new(configmanager.SensorFrequencySetter)),
			fx.As(new(configmanager.GatewayCommissioner)),
			fx.As(new(configmanager.GatewayCreator)),
			fx.As(new(configmanager.GatewayDecommissioner)),
			fx.As(new(configmanager.GatewayDeleter)),
			fx.As(new(configmanager.GatewayInterrupter)),
			fx.As(new(configmanager.GatewayRebooter)),
			fx.As(new(configmanager.GatewayResetter)),
			fx.As(new(configmanager.GatewayResumer)),
			fx.As(new(configmanager.SensorInterrupter)),
			fx.As(new(configmanager.SensorResumer)),
			fx.As(new(configmanager.SensorAdder)),
			fx.As(new(configmanager.SensorDeleter)),
		),
		fx.Annotate(
			configrepositories.NewSQLiteConfigRepository,
			fx.As(new(configmanager.ConfigPort)),
		),
	),
)
