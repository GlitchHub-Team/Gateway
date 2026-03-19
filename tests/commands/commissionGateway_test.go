package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
)

func TestCommissionGatewayCmdExecute(t *testing.T) {
	//verifica che CommissionGatewayCmd salvi il commissioning e poi commissioni il sender
	cmdData := &commanddata.CommissionGateway{
		GatewayId:         uuid.New(),
		TenantId:          uuid.New(),
		CommissionedToken: "token",
	}
	port := &mockGatewayCommissionerPort{}
	commissioner := &mockGatewayCommissioner{}

	cmd := commands.NewCommissionGatewayCmd(cmdData, port, commissioner, domain.Active)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !port.called {
		t.Fatal("expected CommissionGateway to be called")
	}

	if port.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, port.receivedCmd)
	}

	if port.receivedStat != domain.Active {
		t.Fatalf("expected status %q, got %q", domain.Active, port.receivedStat)
	}

	if commissioner.commissionCalls != 1 {
		t.Fatalf("expected Commission to be called once, got %d", commissioner.commissionCalls)
	}

	if commissioner.receivedTenantID != cmdData.TenantId {
		t.Fatalf("expected tenant id %v, got %v", cmdData.TenantId, commissioner.receivedTenantID)
	}

	if commissioner.receivedToken != cmdData.CommissionedToken {
		t.Fatalf("expected token %q, got %q", cmdData.CommissionedToken, commissioner.receivedToken)
	}

	if cmd.String() != "CommissionGatewayCmd" {
		t.Fatalf("expected String() to return CommissionGatewayCmd, got %s", cmd.String())
	}
}

func TestCommissionGatewayCmdExecuteReturnsCommissioningError(t *testing.T) {
	//verifica che CommissionGatewayCmd non commissioni il sender se il port fallisce
	expectedErr := errors.New("commission failed")
	port := &mockGatewayCommissionerPort{err: expectedErr}
	commissioner := &mockGatewayCommissioner{}

	cmd := commands.NewCommissionGatewayCmd(&commanddata.CommissionGateway{}, port, commissioner, domain.Active)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if commissioner.commissionCalls != 0 {
		t.Fatalf("expected Commission not to be called, got %d", commissioner.commissionCalls)
	}
}

func TestCommissionGatewayCmdExecuteReturnsSenderError(t *testing.T) {
	//verifica che CommissionGatewayCmd propaghi l'errore del sender
	expectedErr := errors.New("sender commission failed")
	port := &mockGatewayCommissionerPort{}
	commissioner := &mockGatewayCommissioner{err: expectedErr}

	cmd := commands.NewCommissionGatewayCmd(&commanddata.CommissionGateway{}, port, commissioner, domain.Active)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
