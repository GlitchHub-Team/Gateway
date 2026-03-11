package sensor

import (
	"context"
	"time"

	"Gateway/internal/domain"

	"go.uber.org/zap"
)

type SensorService struct {
	sensor      *Sensor
	bufferPort  SaveSensorDataPort
	cmdChannel  chan domain.BaseCommand
	stopChannel chan struct{}
	errChanel   chan error
	ctx         context.Context
	logger      *zap.Logger
}

func NewSensorService(sensor *Sensor, bufferPort SaveSensorDataPort, cmdChannel chan domain.BaseCommand, stopChannel chan struct{}, errChannel chan error, ctx context.Context, logger *zap.Logger) *SensorService {
	return &SensorService{
		sensor:      sensor,
		bufferPort:  bufferPort,
		cmdChannel:  cmdChannel,
		stopChannel: stopChannel,
		errChanel:   errChannel,
		ctx:         ctx,
		logger:      logger,
	}
}

func (s *SensorService) Start() {
	ticker := time.NewTicker(s.sensor.Interval)

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
				s.errChanel <- err
			}
		case <-ticker.C:
			if s.sensor.Status == Inactive {
				continue
			}
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
	select {
	case s.stopChannel <- struct{}{}:
	default:
	}
}

func (s *SensorService) Interrupt() {
	s.sensor.Status = Inactive
}

func (s *SensorService) Resume() {
	s.sensor.Status = Active
}
