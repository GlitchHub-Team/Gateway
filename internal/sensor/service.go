package sensor

import (
	"context"
	"time"

	"Gateway/internal/domain"

	"go.uber.org/zap"
)

const (
	DefaultFrequency = 5 * time.Second
)

type SensorService struct {
	sensor      *Sensor
	bufferPort  SaveSensorDataPort
	cmdChannel  chan domain.BaseCommand
	stopChannel chan struct{}
	ctx         context.Context
	logger      *zap.Logger
}

func NewSensorService(sensor *Sensor, bufferPort SaveSensorDataPort, cmdChannel chan domain.BaseCommand, stopChannel chan struct{}, ctx context.Context, logger *zap.Logger) *SensorService {
	return &SensorService{
		sensor:      sensor,
		bufferPort:  bufferPort,
		cmdChannel:  cmdChannel,
		stopChannel: stopChannel,
		ctx:         ctx,
		logger:      logger,
	}
}

func (s *SensorService) Start() {
	ticker := time.NewTicker(DefaultFrequency)

	defer ticker.Stop()
	for {
		select {
		case cmd := <-s.cmdChannel:
			if err := cmd.Execute(); err != nil {
				s.logger.Error("Errore nell'esecuzione del comando",
					zap.String("command", cmd.String()),
					zap.String("sensorId", s.sensor.Id.String()),
					zap.Error(err),
				)
			}
		case <-ticker.C:
			err := s.generateData()
			if err != nil {
				s.logger.Error("Errore nella generazione dei dati del sensore",
					zap.String("sensorId", s.sensor.Id.String()),
					zap.Error(err),
				)
			}
		case <-s.stopChannel:
			return
		case <-s.ctx.Done():
			s.logger.Error("Sensore interrotto",
				zap.String("sensorId", s.sensor.Id.String()),
				zap.Error(s.ctx.Err()),
			)
			return
		}
	}
}

func (s *SensorService) generateData() error {
	sensorData := s.sensor.Profile.Generate()
	if err := s.bufferPort.Save(sensorData, s.sensor.GatewayId); err != nil {
		return err
	}
	return nil
}

func (s *SensorService) Stop() {
	s.stopChannel <- struct{}{}
}
