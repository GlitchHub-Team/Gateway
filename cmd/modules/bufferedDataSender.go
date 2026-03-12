package modules

import (
	bufferedDataSender "Gateway/internal/bufferedDataSender"

	"go.uber.org/fx"
)

var BufferedDataSenderModule = fx.Module("BufferedDataSenderModule",
	fx.Provide(
		fx.Annotate(
			bufferedDataSender.NewBufferedDataRepository,
			fx.As(new(bufferedDataSender.BufferedDataPort)),
		),
		fx.Annotate(
			bufferedDataSender.NewNATSDataPublisherFactory,
			fx.As(new(bufferedDataSender.SendSensorDataPortFactory)),
		),
	),
)
