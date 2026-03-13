package sensor

import (
	"context"
	"time"

	"Gateway/internal/domain"

	"go.uber.org/zap"
)

type SensorService struct {
	sensor     *Sensor
	bufferPort SaveSensorDataPort
	cmdChannel chan domain.BaseCommand
	errChanel  chan error
	ctx        context.Context
	logger     *zap.Logger
}

func NewSensorService(sensor *Sensor, bufferPort SaveSensorDataPort, cmdChannel chan domain.BaseCommand, errChannel chan error, ctx context.Context, logger *zap.Logger) *SensorService {
	return &SensorService{
		sensor:     sensor,
		bufferPort: bufferPort,
		cmdChannel: cmdChannel,
		errChanel:  errChannel,
		ctx:        ctx,
		logger:     logger,
	}
}

func (s *SensorService) Start() {
	ticker := time.NewTicker(s.sensor.Interval)

	defer ticker.Stop()
	for s.sensor.Status != Stopped {
		select {
		case cmd := <-s.cmdChannel:
			err := cmd.Execute()
			if err != nil {
				s.logger.Error("Errore nell'esecuzione del comando",
					zap.String("command", cmd.String()),
					zap.String("sensorId", s.sensor.Id.String()),
					zap.Error(err),
				)
			}
			s.errChanel <- err
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
		case <-s.ctx.Done():
			s.logger.Warn("Sensore interrotto",
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

// Alert: in caso si vogliano chiamare dall'esterno della goroutine bisogna istanziare un mutex
func (s *SensorService) Stop() {
	s.sensor.Status = Stopped
}

func (s *SensorService) Interrupt() {
	s.sensor.Status = Inactive
}

func (s *SensorService) Resume() {
	s.sensor.Status = Active
}
