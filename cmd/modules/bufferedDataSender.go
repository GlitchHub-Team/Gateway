package modules

import (
	bufferedDataSender "Gateway/internal/bufferedDataSender"

	"go.uber.org/fx"
)

var BufferedDataSenderModule = fx.Module("BufferedDataSenderModule",
	fx.Provide(
		fx.Annotate(
			bufferedDataSender.NewNATSDataPublisherRepository,
			fx.As(new(bufferedDataSender.SendSensorDataPort)),
		),
		fx.Annotate(
			bufferedDataSender.NewBufferedDataRepository,
			fx.As(new(bufferedDataSender.BufferedDataPort)),
		),
	),
)
