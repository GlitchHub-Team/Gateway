package commandstests

import (
	"context"
	"errors"
	"testing"
	"time"

	commands "Gateway/internal/commands"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestRebootGatewayCmdExecute(t *testing.T) {
	//verifica che RebootGatewayCmd completi il riavvio dopo la durata richiesta
	cmd := commands.NewRebootGatewayCmd(
		commanddata.RebootGateway{GatewayId: uuid.New()},
		10*time.Millisecond,
		context.Background(),
		zap.NewNop(),
	)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cmd.String() != "RebootGatewayCmd" {
		t.Fatalf("expected String() to return RebootGatewayCmd, got %s", cmd.String())
	}
}

func TestRebootGatewayCmdExecuteReturnsContextError(t *testing.T) {
	//verifica che RebootGatewayCmd termini con errore se il context viene cancellato
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cmd := commands.NewRebootGatewayCmd(
		commanddata.RebootGateway{GatewayId: uuid.New()},
		time.Second,
		ctx,
		zap.NewNop(),
	)

	err := cmd.Execute()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected error %v, got %v", context.Canceled, err)
	}
}
