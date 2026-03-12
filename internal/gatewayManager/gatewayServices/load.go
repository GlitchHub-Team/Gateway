package gatewayservices

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	"Gateway/internal/domain"
	gatewaymanager "Gateway/internal/gatewayManager"
	"Gateway/internal/sensor"

	"go.uber.org/zap"
)

func (gatManager *GatewayManagerService) LoadGatewayWorkers() error {
	gateways, err := gatManager.configPort.GetAllGateways()
	if err != nil {
		return err
	}

	for id, gateway := range gateways {
		errGatChannel := make(chan error)
		cmdGatChannel := make(chan domain.BaseCommand)

		var sendSensorDataPort buffereddatasender.SendSensorDataPort
		if gateway.Token == nil {
			sendSensorDataPort, err = gatManager.sendSensorDataPortFactory.Create()
		} else {
			sendSensorDataPort, err = gatManager.sendSensorDataPortFactory.Reload(*gateway.Token, gateway.SecretKey)
		}

		if err != nil {
			gatManager.logger.Error("Errore nella creazione del SendSensorDataPort",
				zap.String("gatewayId", id.String()),
				zap.Error(err),
			)
			return err
		}

		dataSender := buffereddatasender.NewBufferedDataSenderService(
			gateway,
			sendSensorDataPort,
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
