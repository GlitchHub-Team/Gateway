package modules

import (
	configmanager "Gateway/internal/configManager"
	configrepositories "Gateway/internal/configManager/configRepositories"

	"go.uber.org/fx"
)

var ConfigManagerModule = fx.Module("ConfigManagerModule",
	fx.Provide(
		fx.Annotate(
			configrepositories.NewSQLiteConfigRepository,
			fx.As(new(configmanager.ConfigPort)),
			fx.As(new(configmanager.GatewaysFetcherPort)),
			fx.As(new(configmanager.GatewayCommissionerPort)),
			fx.As(new(configmanager.GatewayCreatorPort)),
			fx.As(new(configmanager.GatewayDecommissionerPort)),
			fx.As(new(configmanager.GatewayDeleterPort)),
			fx.As(new(configmanager.GatewayInterrupterPort)),
			fx.As(new(configmanager.GatewayResetterPort)),
			fx.As(new(configmanager.GatewayResumerPort)),
			fx.As(new(configmanager.SensorInterrupterPort)),
			fx.As(new(configmanager.SensorResumerPort)),
			fx.As(new(configmanager.SensorAdderPort)),
			fx.As(new(configmanager.SensorDeleterPort)),
		),
	),
)
