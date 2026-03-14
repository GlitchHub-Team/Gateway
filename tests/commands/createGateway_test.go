package commandstests

import (
	"errors"
	"testing"
	"time"

	commands "Gateway/internal/commands"
	credentialsgenerator "Gateway/internal/credentialsGenerator"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"

	"github.com/google/uuid"
)

func TestCreateGatewayCmdExecute(t *testing.T) {
	//verifica che CreateGatewayCmd crei il gateway, saluti e avvii il sender
	cmdData := &commanddata.CreateGateway{GatewayId: uuid.New(), Interval: time.Second}
	credentials := &credentialsgenerator.Credentials{
		PublicIdentifier: "public",
		SecretKey:        "secret",
	}
	creator := &mockGatewayCreator{}
	greeter := &mockGatewayGreeter{}
	starter := &mockGatewayStarter{started: make(chan struct{}, 1)}

	cmd := commands.NewCreateGatewayCmd(cmdData, creator, starter, greeter, credentials, domain.Decommissioned)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !creator.called {
		t.Fatal("expected CreateGateway to be called")
	}

	if creator.receivedCmd != cmdData {
		t.Fatalf("expected cmdData %p, got %p", cmdData, creator.receivedCmd)
	}

	if creator.receivedCredentials != credentials {
		t.Fatalf("expected credentials %p, got %p", credentials, creator.receivedCredentials)
	}

	if creator.receivedStat != domain.Decommissioned {
		t.Fatalf("expected status %q, got %q", domain.Decommissioned, creator.receivedStat)
	}

	if greeter.helloCalls != 1 {
		t.Fatalf("expected Hello to be called once, got %d", greeter.helloCalls)
	}

	waitForSignal(t, starter.started, "gateway start")

	if cmd.String() != "CreateGatewayCmd" {
		t.Fatalf("expected String() to return CreateGatewayCmd, got %s", cmd.String())
	}
}

func TestCreateGatewayCmdExecuteReturnsCreateError(t *testing.T) {
	//verifica che CreateGatewayCmd non saluti ne avvii il sender se la create fallisce
	expectedErr := errors.New("create failed")
	creator := &mockGatewayCreator{err: expectedErr}
	greeter := &mockGatewayGreeter{}
	starter := &mockGatewayStarter{started: make(chan struct{}, 1)}

	cmd := commands.NewCreateGatewayCmd(&commanddata.CreateGateway{}, creator, starter, greeter, &credentialsgenerator.Credentials{}, domain.Decommissioned)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if greeter.helloCalls != 0 {
		t.Fatalf("expected Hello not to be called, got %d", greeter.helloCalls)
	}

	select {
	case <-starter.started:
		t.Fatal("expected Start not to be called on create error")
	default:
	}
}

func TestCreateGatewayCmdExecuteReturnsHelloError(t *testing.T) {
	//verifica che CreateGatewayCmd non avvii il sender se Hello fallisce
	expectedErr := errors.New("hello failed")
	creator := &mockGatewayCreator{}
	greeter := &mockGatewayGreeter{err: expectedErr}
	starter := &mockGatewayStarter{started: make(chan struct{}, 1)}

	cmd := commands.NewCreateGatewayCmd(&commanddata.CreateGateway{}, creator, starter, greeter, &credentialsgenerator.Credentials{}, domain.Decommissioned)

	err := cmd.Execute()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	select {
	case <-starter.started:
		t.Fatal("expected Start not to be called on hello error")
	default:
	}
}
