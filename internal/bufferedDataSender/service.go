package buffereddatasender

import (
	"context"
	"time"

	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"

	"go.uber.org/zap"
)

const (
	DefaultFrequency = 5 * time.Second
)

type BufferedDataSenderService struct {
	gateway          *configmanager.Gateway
	sendDataRepo     SendSensorDataPort
	bufferedDataPort BufferedDataPort
	cmdChannel       chan domain.BaseCommand
	stopChannel      chan struct{}
	ctx              context.Context
	logger           *zap.Logger
}

func NewBufferedDataSenderService(gateway *configmanager.Gateway, sendDataRepo SendSensorDataPort, bufferedDataPort BufferedDataPort, cmdChannel chan domain.BaseCommand, stopChannel chan struct{}, ctx context.Context, logger *zap.Logger) *BufferedDataSenderService {
	return &BufferedDataSenderService{
		gateway:          gateway,
		sendDataRepo:     sendDataRepo,
		bufferedDataPort: bufferedDataPort,
		cmdChannel:       cmdChannel,
		stopChannel:      stopChannel,
		ctx:              ctx,
		logger:           logger,
	}
}

func (b *BufferedDataSenderService) Start() {
	ticker := time.NewTicker(DefaultFrequency)

	defer ticker.Stop()
	for {
		select {
		case cmd := <-b.cmdChannel:
			if err := cmd.Execute(); err != nil {
				b.logger.Error("Errore nell'esecuzione del comando",
					zap.String("command", cmd.String()),
					zap.String("gatewayId", b.gateway.Id.String()),
					zap.Error(err),
				)
			}
		case <-ticker.C:
			err := b.sendBufferedData()
			if err != nil {
				b.logger.Error("Errore nell'invio dei dati bufferizzati",
					zap.String("gatewayId", b.gateway.Id.String()),
					zap.Error(err),
				)
			}
		case <-b.stopChannel:
			return
		case <-b.ctx.Done():
			b.logger.Error("Gateway interrotto",
				zap.String("gatewayId", b.gateway.Id.String()),
				zap.Error(b.ctx.Err()),
			)
			return
		}
	}
}

func (b *BufferedDataSenderService) sendBufferedData() error {
	confirmedData := []*sensorData{}
	data, err := b.bufferedDataPort.GetOrderedBufferedData(b.gateway.Id)
	if err != nil {
		return err
	}

	for _, d := range data {
		if err := b.sendDataRepo.Send(d); err != nil {
			b.logger.Error("Errore nell'invio dei dati del gateway",
				zap.String("gatewayId", b.gateway.Id.String()),
				zap.Error(err),
			)
			break
		}
		confirmedData = append(confirmedData, d)
	}

	if err := b.bufferedDataPort.CleanBufferedData(confirmedData); err != nil {
		b.logger.Error("RISCHIO DATI DUPLICATI: Errore nella pulitura dei dati bufferizzati",
			zap.String("gatewayId", b.gateway.Id.String()),
			zap.Error(err),
		)
	}

	return nil
}

func (b *BufferedDataSenderService) Stop() {
	select {
	case b.stopChannel <- struct{}{}:
	default:
	}
}
