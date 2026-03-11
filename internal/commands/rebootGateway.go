package commands

import (
	"context"
	"time"

	commanddata "Gateway/internal/gatewayManager/commandData"

	"go.uber.org/zap"
)

type RebootGatewayCmd struct {
	cmdData        commanddata.RebootGateway
	rebootDuration time.Duration
	ctx            context.Context
	logger         *zap.Logger
}

func (c *RebootGatewayCmd) Execute() error {
	c.logger.Info("Riavvio in corso...", zap.String("gatewayId", c.cmdData.GatewayId.String()))

	select {
	case <-time.After(c.rebootDuration):
		c.logger.Info("Gateway riavviato con successo", zap.String("gatewayId", c.cmdData.GatewayId.String()))
	case <-c.ctx.Done():
		c.logger.Warn("Reboot interrotto dallo shutdown dell'applicazione", zap.String("gatewayId", c.cmdData.GatewayId.String()))
		return c.ctx.Err()
	}

	return nil
}

func NewRebootGatewayCmd(cmdData commanddata.RebootGateway, rebootDuration time.Duration, ctx context.Context, logger *zap.Logger) *RebootGatewayCmd {
	return &RebootGatewayCmd{
		cmdData:        cmdData,
		rebootDuration: rebootDuration,
		ctx:            ctx,
		logger:         logger,
	}
}

func (c *RebootGatewayCmd) String() string {
	return "RebootGatewayCmd"
}
