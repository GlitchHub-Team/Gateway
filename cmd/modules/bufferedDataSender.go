package modules

import (
	bufferedDataSender "Gateway/internal/bufferedDataSender"

	"go.uber.org/fx"
)

var BufferedDataSenderModule = fx.Module("BufferedDataSenderModule",
	fx.Provide(
		fx.Annotate(
			bufferedDataSender.NewBufferedDataSenderService,
			fx.As(new(bufferedDataSender.DataSender)),
		),
		fx.Annotate(
			bufferedDataSender.NewNATSDataPublisherRepository,
			fx.As(new(bufferedDataSender.SendSensorDataRepository)),
		),
	),
)
