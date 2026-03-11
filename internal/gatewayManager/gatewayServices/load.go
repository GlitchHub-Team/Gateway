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
		dataSender := buffereddatasender.NewBufferedDataSenderService(
			gateway,
			gatManager.sendSensorDataPort,
			gatManager.bufferedDataPort,
			make(chan domain.BaseCommand),
			make(chan struct{}),
			make(chan error),
			gatManager.ctx,
			gatManager.logger,
		)

		gatManager.gateways.Mu.Lock()
		gatManager.gateways.Workers[id] = dataSender
		gatManager.gateways.Mu.Unlock()

		gatManager.sensors.Mu.Lock()
		if gatManager.sensors.Workers[id] == nil {
			gatManager.sensors.Workers[id] = make(map[gatewaymanager.SensorId]sensor.SimulatedSensor)
		}
		gatManager.sensors.Mu.Unlock()

		for sensorId, sensorEntity := range gateway.Sensors {
			sensorService := sensor.NewSensorService(
				sensorEntity,
				gatManager.saveSensorDataPort,
				make(chan domain.BaseCommand),
				make(chan struct{}),
				make(chan error),
				gatManager.ctx,
				gatManager.logger,
			)

			gatManager.sensors.Mu.Lock()
			gatManager.sensors.Workers[id][sensorId] = sensorService
			gatManager.sensors.Mu.Unlock()

			go sensorService.Start()
		}

		go dataSender.Start()
	}
	return nil
}
