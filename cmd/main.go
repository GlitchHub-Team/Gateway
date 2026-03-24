package main

import (
	"context"
	"os"
	"strconv"

	bufferdatabase "Gateway/cmd/external/bufferDatabase"
	gatewaydatabase "Gateway/cmd/external/gatewayDatabase"
	"Gateway/cmd/logger"
	"Gateway/cmd/modules"
	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"
	gatewayservices "Gateway/internal/gatewayManager/gatewayServices"
	"Gateway/internal/natsutil"
	sensorprofiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

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
		fx.Supply(natsutil.NatsAddress(os.Getenv("NATS_HOST"))),
		fx.Supply(natsutil.NatsPort(envInt("NATS_PORT", 4222))),
		fx.Supply(natsutil.NatsCredsPath(os.Getenv("BASE_CREDS_PATH"))),
		fx.Supply(natsutil.NatsCAPemPath(os.Getenv("CA_PEM_PATH"))),
		fx.Provide(natsutil.NewNATSConnection),
		fx.Provide(natsutil.NewJetStreamContext),
		fx.Provide(gatewaydatabase.NewGatewayDatabase),
		fx.Provide(bufferdatabase.NewBufferDatabase),
		fx.Provide(sensorprofiles.NewRand),
		fx.Provide(func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		}),

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
