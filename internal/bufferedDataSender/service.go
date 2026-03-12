package buffereddatasender

import (
	"context"
	"time"

	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BufferedDataSenderService struct {
	gateway                   *configmanager.Gateway
	sendDataRepo              SendSensorDataPort
	sendSensorDataPortFactory SendSensorDataPortFactory
	bufferedDataPort          BufferedDataPort
	cmdChannel                chan domain.BaseCommand
	errChanel                 chan error
	ctx                       context.Context
	logger                    *zap.Logger
	ticker                    *time.Ticker
}

func NewBufferedDataSenderService(gateway *configmanager.Gateway, sendDataRepo SendSensorDataPort, bufferedDataPort BufferedDataPort, sendSensorDataPortFactory SendSensorDataPortFactory, cmdChannel chan domain.BaseCommand, errChannel chan error, ctx context.Context, logger *zap.Logger) *BufferedDataSenderService {
	return &BufferedDataSenderService{
		gateway:                   gateway,
		sendDataRepo:              sendDataRepo,
		sendSensorDataPortFactory: sendSensorDataPortFactory,
		bufferedDataPort:          bufferedDataPort,
		cmdChannel:                cmdChannel,
		errChanel:                 errChannel,
		ctx:                       ctx,
		logger:                    logger,
		ticker:                    time.NewTicker(gateway.Interval),
	}
}

func (b *BufferedDataSenderService) Start() {
	defer b.ticker.Stop()

	for b.gateway.Status != domain.Stopped {
		select {
		case cmd := <-b.cmdChannel:
			err := cmd.Execute()
			if err != nil {
				b.logger.Error("Errore nell'esecuzione del comando",
					zap.String("command", cmd.String()),
					zap.String("gatewayId", b.gateway.Id.String()),
					zap.Error(err),
				)
			}
			b.errChanel <- err
		case <-b.ticker.C:
			if b.gateway.Status == domain.Active {
				err := b.sendBufferedData()
				if err != nil {
					b.logger.Error("Errore nell'invio dei dati bufferizzati",
						zap.String("gatewayId", b.gateway.Id.String()),
						zap.Error(err),
					)
				}
			}
		case <-b.ctx.Done():
			b.logger.Warn("Gateway interrotto",
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
		if err := b.sendDataRepo.Send(d, *b.gateway.TenantId); err != nil {
			b.logger.Error("Errore nell'invio dei dati del gateway",
				zap.String("gatewayId", b.gateway.Id.String()),
				zap.String("tenantId", b.gateway.TenantId.String()),
				zap.Error(err),
			)
			continue
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

func (b *BufferedDataSenderService) Hello() error {
	if err := b.sendDataRepo.Hello(b.gateway.Id, b.gateway.PublicIdentifier); err != nil {
		return err
	}

	return nil
}

func (b *BufferedDataSenderService) Decommission() error {
	if err := b.bufferedDataPort.CleanWholeBuffer(b.gateway.Id); err != nil {
		return err
	}

	b.sendDataRepo = b.sendSensorDataPortFactory.Create()
	b.gateway.Status = domain.Decommissioned
	b.gateway.TenantId = nil
	b.gateway.Token = nil
	return nil
}

func (b *BufferedDataSenderService) Commission(tenantId uuid.UUID, commissionedToken string) error {
	sendDataPort, err := b.sendSensorDataPortFactory.Reload(commissionedToken, b.gateway.SecretKey)
	if err != nil {
		return err
	}

	b.sendDataRepo = sendDataPort
	b.gateway.Status = domain.Active
	b.gateway.TenantId = &tenantId
	b.gateway.Token = &commissionedToken

	return nil
}

func (b *BufferedDataSenderService) Reset(defaultInterval time.Duration) error {
	b.gateway.Interval = defaultInterval
	b.ticker.Reset(defaultInterval)
	if err := b.bufferedDataPort.CleanWholeBuffer(b.gateway.Id); err != nil {
		return err
	}
	return nil
}

// Alert: in caso si vogliano chiamare dall'esterno della goroutine bisogna istanziare un mutex
func (b *BufferedDataSenderService) Stop() {
	b.gateway.Status = domain.Stopped
}

func (b *BufferedDataSenderService) Interrupt() {
	b.gateway.Status = domain.Inactive
}

func (b *BufferedDataSenderService) Resume() {
	b.gateway.Status = domain.Active
}
