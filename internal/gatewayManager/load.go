package gatewaymanager

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	"Gateway/internal/domain"
	"Gateway/internal/sensor"
)

type GatewaysLoader interface {
	LoadGatewayWorkers() error
}

func (gatManager *GatewayManagerService) LoadGatewayWorkers() error {
	gateways, err := gatManager.gatFetcher.GetAllGateways()
	if err != nil {
		return err
	}

	for id, gateway := range gateways {
		gatManager.gateways[id] = buffereddatasender.NewBufferedDataSenderService(
			gateway,
			gatManager.sendSensorDataPort,
			gatManager.bufferedDataPort,
			make(chan domain.BaseCommand),
			make(chan struct{}),
			gatManager.ctx,
			gatManager.logger,
		)

		if gatManager.sensors[id] == nil {
			gatManager.sensors[id] = make(map[SensorId]sensor.SimulatedSensor)
		}

		for sensorId, sensorEntity := range gateway.Sensors {
			gatManager.sensors[id][sensorId] = sensor.NewSensorService(
				sensorEntity,
				gatManager.saveSensorDataPort,
				make(chan domain.BaseCommand),
				make(chan struct{}),
				gatManager.ctx,
				gatManager.logger,
			)
			go gatManager.sensors[id][sensorId].Start()
		}
		go gatManager.gateways[id].Start()
	}
	return nil
}
