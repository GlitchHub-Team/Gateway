package gatewayservices

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	"Gateway/internal/domain"
	gatewaymanager "Gateway/internal/gatewayManager"
	"Gateway/internal/sensor"
)

func (gatManager *GatewayManagerService) LoadGatewayWorkers() error {
	gateways, err := gatManager.configPort.GetAllGateways()
	if err != nil {
		return err
	}

	for id, gateway := range gateways {
		errGatChannel := make(chan error)
		cmdGatChannel := make(chan domain.BaseCommand)

		dataSender := buffereddatasender.NewBufferedDataSenderService(
			gateway,
			gatManager.sendSensorDataPort,
			gatManager.bufferedDataPort,
			cmdGatChannel,
			errGatChannel,
			gatManager.ctx,
			gatManager.logger,
		)

		gatManager.gateways.Mu.Lock()
		gatManager.gateways.Workers[id] = gatewaymanager.GatewayWorker{
			Sender:     dataSender,
			ErrChannel: errGatChannel,
			CmdChannel: cmdGatChannel,
		}
		gatManager.gateways.Mu.Unlock()

		gatManager.sensors.Mu.Lock()
		if gatManager.sensors.Workers[id] == nil {
			gatManager.sensors.Workers[id] = make(map[gatewaymanager.SensorId]gatewaymanager.SensorWorker)
		}
		gatManager.sensors.Mu.Unlock()

		for sensorId, sensorEntity := range gateway.Sensors {
			errSensorChannel := make(chan error)
			cmdSensorChannel := make(chan domain.BaseCommand)

			sensorService := sensor.NewSensorService(
				sensorEntity,
				gatManager.saveSensorDataPort,
				cmdSensorChannel,
				errSensorChannel,
				gatManager.ctx,
				gatManager.logger,
			)

			gatManager.sensors.Mu.Lock()
			gatManager.sensors.Workers[id][sensorId] = gatewaymanager.SensorWorker{
				SimulatedSensor: sensorService,
				ErrChannel:      errSensorChannel,
				CmdChannel:      cmdSensorChannel,
			}
			gatManager.sensors.Mu.Unlock()

			go sensorService.Start()
		}

		go dataSender.Start()
	}
	return nil
}
