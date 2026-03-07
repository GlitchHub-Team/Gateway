package main

import (
	"go.uber.org/fx"

	bufferdatabase "Gateway/cmd/external/bufferDatabase"
	gatewaydatabase "Gateway/cmd/external/gatewayDatabase"
	natsserver "Gateway/cmd/external/natsServer"
	"Gateway/cmd/logger"
	modules "Gateway/cmd/modules"
	commandcontrollers "Gateway/internal/gateway/commandControllers"
	sensorprofiles "Gateway/internal/sensor/sensorProfiles"
)

func main() {
	fx.New(
		fx.Provide(logger.NewLogger),
		fx.WithLogger(logger.GetFxLogger),

		fx.Supply(commandcontrollers.AddSensorSubject("commands.addsensor")),
		fx.Supply(commandcontrollers.ChangeSensorFrequencySubject("commands.changesensorfrequency")),
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
		fx.Supply(natsserver.NatsAddress("localhost")),
		fx.Supply(natsserver.NatsPort(8888)),
		fx.Provide(natsserver.NewNATSConnection),
		fx.Provide(gatewaydatabase.NewGatewayDatabase),
		fx.Provide(bufferdatabase.NewBufferDatabase),
		fx.Provide(sensorprofiles.NewRand),

		modules.BufferedDataSenderModule,
		modules.SensorModule,
		modules.ConfigManagerModule,
		modules.GatewayModule,
		modules.CommandControllersModule,

		fx.Invoke(Init),
	).Run()
}

func Init(lc fx.Lifecycle) {
	lc.Append(fx.Hook{})
}
