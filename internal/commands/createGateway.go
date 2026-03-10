package commands

import (
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configmanager "Gateway/internal/configManager"
	"Gateway/internal/domain"
	gatewaymanager "Gateway/internal/gatewayManager"
	commanddata "Gateway/internal/gatewayManager/commandData"
	"Gateway/internal/sensor"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateGatewayCmd struct {
	cmdData            *commanddata.CreateGateway
	configPort         configmanager.GatewayCreatorPort
	gatewayWorkers     *gatewaymanager.GatewayWorkers
	sendSensorDataPort buffereddatasender.SendSensorDataPort
	bufferedDataPort   buffereddatasender.BufferedDataPort
	ctx                context.Context
	logger             *zap.Logger
}

func (c *CreateGatewayCmd) Execute() error {
	if err := c.configPort.CreateGateway(c.cmdData); err != nil {
		return err
	}

	gateway := &configmanager.Gateway{
		Id:       c.cmdData.GatewayId,
		TenantId: nil,
		Sensors:  make(map[uuid.UUID]*sensor.Sensor),
		Status:   configmanager.Inactive,
		Interval: c.cmdData.Interval,
	}

	dataSender := buffereddatasender.NewBufferedDataSenderService(
		gateway,
		c.sendSensorDataPort,
		c.bufferedDataPort,
		make(chan domain.BaseCommand),
		make(chan struct{}),
		c.ctx,
		c.logger,
	)

	c.gatewayWorkers.Mu.Lock()
	c.gatewayWorkers.Workers[c.cmdData.GatewayId] = dataSender
	c.gatewayWorkers.Mu.Unlock()

	go dataSender.Start()

	return nil
}

func NewCreateGatewayCmd(cmdData *commanddata.CreateGateway, configPort configmanager.GatewayCreatorPort, gatewayWorkers *gatewaymanager.GatewayWorkers, sendSensorDataPort buffereddatasender.SendSensorDataPort, bufferedDataPort buffereddatasender.BufferedDataPort, ctx context.Context, logger *zap.Logger) *CreateGatewayCmd {
	return &CreateGatewayCmd{
		cmdData:            cmdData,
		configPort:         configPort,
		gatewayWorkers:     gatewayWorkers,
		sendSensorDataPort: sendSensorDataPort,
		bufferedDataPort:   bufferedDataPort,
		ctx:                ctx,
		logger:             logger,
	}
}

func (c *CreateGatewayCmd) String() string {
	return "CreateGatewayCmd"
}
