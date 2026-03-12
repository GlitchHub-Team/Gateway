package main

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	bufferdatabase "Gateway/cmd/external/bufferDatabase"
	gatewaydatabase "Gateway/cmd/external/gatewayDatabase"
	natsserver "Gateway/cmd/external/natsServer"
	"Gateway/cmd/logger"
	modules "Gateway/cmd/modules"
	buffereddatasender "Gateway/internal/bufferedDataSender"
	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"
	gatewayservices "Gateway/internal/gatewayManager/gatewayServices"
	sensorprofiles "Gateway/internal/sensor/sensorProfiles"
)

func main() {
	fx.New(
		fx.Provide(logger.NewLogger),
		fx.WithLogger(logger.GetFxLogger),

		fx.Supply(commandcontrollers.AddSensorSubject("commands.addsensor")),
		fx.Supply(commandcontrollers.CreateGatewaySubject("commands.creategateway")),
		fx.Supply(commandcontrollers.CommissionGatewaySubject("commands.commissiongateway")),
		fx.Supply(commandcontrollers.DecommissionGatewaySubject("commands.decommissiongateway")),
		fx.Supply(commandcontrollers.DeleteGatewaySubject("commands.deletegateway")),
		fx.Supply(commandcontrollers.DeleteSensorSubject("commands.deletesensor")),
		fx.Supply(commandcontrollers.InterruptGatewaySubject("commands.interruptgateway")),
		fx.Supply(commandcontrollers.InterruptSensorSubject("commands.interruptsensor")),
		fx.Supply(commandcontrollers.RebootGatewaySubject("commands.rebootgateway")),
		fx.Supply(commandcontrollers.ResetGatewaySubject("commands.resetgateway")),
		fx.Supply(commandcontrollers.ResumeGatewaySubject("commands.resumegateway")),
		fx.Supply(commandcontrollers.ResumeSensorSubject("commands.resumesensor")),
		fx.Supply(buffereddatasender.NatsAddress("localhost")),
		fx.Supply(buffereddatasender.NatsPort(4222)),
		fx.Supply(buffereddatasender.BaseToken("")),
		fx.Supply(buffereddatasender.BaseSeed("")),
		fx.Provide(natsserver.NewNATSConnection),
		fx.Provide(gatewaydatabase.NewGatewayDatabase),
		fx.Provide(bufferdatabase.NewBufferDatabase),
		fx.Provide(sensorprofiles.NewRand),
		fx.Provide(func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		}),

		modules.CredGenerator,
		modules.BufferedDataSenderModule,
		modules.SensorModule,
		modules.ConfigManagerModule,
		modules.GatewayModule,
		modules.CommandControllersModule,

		fx.Invoke(Init),
	).Run()
}

type InitParams struct {
	fx.In

	Lc          fx.Lifecycle
	Controllers []commandcontrollers.NATSCommandController `group:"nats_controllers"`

	Loader gatewayservices.GatewaysLoader
	Nc     *nats.Conn
	Cancel context.CancelFunc
	Logger *zap.Logger
}

func Init(p InitParams) {
	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := p.Loader.LoadGatewayWorkers(); err != nil {
				return err
			}
			p.Logger.Info("Gateway e sensori salvati avviati correttamente")

			for _, controller := range p.Controllers {
				controller.Listen()
			}
			p.Logger.Info("I NATS Controller sono tutti pronti a ricevere comandi")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := p.Nc.Drain()
			p.Cancel()
			return err
		},
	})
}
