package commandstests

import (
	"errors"
	"testing"

	commands "Gateway/internal/commands"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
)

func TestResumeGatewayCmdExecute(t *testing.T) {
	//verifica che ResumeGatewayCmd aggiorni lo stato e riprenda il sender
	cmdData := &commanddata.ResumeGateway{GatewayId: uuid.New()}
	port := &mockGatewayResumerPort{}
	resumer := &mockGatewayResumer{}

	cmd := commands.NewResumeGatewayCmd(cmdData, resumer, port, domain.Active)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !port.called {
		t.Fatal("expected ResumeGateway to be called")
	}

	if port.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, port.receivedCmd)
	}

	if port.receivedStat != domain.Active {
		t.Fatalf("expected status %q, got %q", domain.Active, port.receivedStat)
	}

	if resumer.resumeCalls != 1 {
		t.Fatalf("expected Resume to be called once, got %d", resumer.resumeCalls)
	}

	if cmd.String() != "ResumeGatewayCmd" {
		t.Fatalf("expected String() to return ResumeGatewayCmd, got %s", cmd.String())
	}
}

func TestResumeGatewayCmdExecuteReturnsResumeError(t *testing.T) {
	//verifica che ResumeGatewayCmd non riprenda il sender se il port fallisce
	expectedErr := errors.New("resume failed")
	port := &mockGatewayResumerPort{err: expectedErr}
	resumer := &mockGatewayResumer{}

	cmd := commands.NewResumeGatewayCmd(&commanddata.ResumeGateway{}, resumer, port, domain.Active)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if resumer.resumeCalls != 0 {
		t.Fatalf("expected Resume not to be called, got %d", resumer.resumeCalls)
	}
}
