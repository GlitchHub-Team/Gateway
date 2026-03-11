package modules

import (
	sensor "Gateway/internal/sensor"

	"go.uber.org/fx"
)

var SensorModule = fx.Module("SensorModule",
	fx.Provide(
		fx.Annotate(
			sensor.NewSQLiteSaveSensorDataRepository,
			fx.As(new(sensor.SaveSensorDataPort)),
		),
	),
)
