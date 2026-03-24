package integration_tests

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	bufferdatabase "Gateway/cmd/external/bufferDatabase"
	gatewaydatabase "Gateway/cmd/external/gatewayDatabase"
	buffereddatasender "Gateway/internal/bufferedDataSender"
	configrepositories "Gateway/internal/configManager/configRepositories"
	credentialsgenerator "Gateway/internal/credentialsGenerator"
	"Gateway/internal/domain"
	gatewaymanager "Gateway/internal/gatewayManager"
	commandcontrollers "Gateway/internal/gatewayManager/commandControllers"
	commanddata "Gateway/internal/gatewayManager/commandData"
	gatewayservices "Gateway/internal/gatewayManager/gatewayServices"
	"Gateway/internal/natsutil"
	sensor "Gateway/internal/sensor"
	sensorprofiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const commissionGatewaySubject = "commands.commissiongateway"

var (
	gateway1ID     = mustParseUUID("11111111-1111-1111-1111-111111111111")
	gateway2ID     = mustParseUUID("22222222-2222-2222-2222-222222222222")
	adminGatewayID = mustParseUUID("33333333-3333-3333-3333-333333333333")
	tenant1ID      = gateway1ID
	adminTenantID  = mustParseUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
)

type commissionGatewayRequest struct {
	GatewayId         string `json:"gatewayId"`
	TenantId          string `json:"tenantId"`
	CommissionedToken string `json:"commissionedToken"`
}

type commissionFixture struct {
	ctx          context.Context
	cancel       context.CancelFunc
	gatewayDb    *configrepositories.ConfigDbConnection
	bufferDb     sensor.BufferDbConnection
	configRepo   *configrepositories.SQLiteConfigRepository
	service      *gatewayservices.GatewayManagerService
	controllerNc *nats.Conn
	publisherNc  *nats.Conn
	observerNc   *nats.Conn
	repoRoot     string
	caPath       string
}

func TestNATSCommissionGatewayIntegration(t *testing.T) {
	t.Run("commissioning valido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		creds := fx.gateway1Creds(t)
		sensorID := uuid.New()
		fx.insertBufferedData(t, gateway1ID, sensorID)

		subject := sensorSubject(tenant1ID, gateway1ID, sensorID)
		sub, err := fx.observerNc.SubscribeSync(subject)
		if err != nil {
			t.Fatalf("cannot subscribe to sensor subject: %v", err)
		}
		if err := fx.observerNc.FlushTimeout(2 * time.Second); err != nil {
			t.Fatalf("cannot flush observer subscription: %v", err)
		}

		res := fx.sendCommissionCommand(t, gateway1ID, tenant1ID, creds.JWT)
		if !res.Success {
			t.Fatalf("expected success response, got %+v", res)
		}
		if res.Message != "Gateway commissionato correttamente" {
			t.Fatalf("unexpected success message: %q", res.Message)
		}

		fx.assertCommissionedState(t, gateway1ID, tenant1ID, creds.JWT)

		_, err = sub.NextMsg(3 * time.Second)
		if err != nil {
			t.Fatalf("expected data publish after commissioning: %v", err)
		}
	})

	t.Run("tenantId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		creds := fx.gateway1Creds(t)
		res := fx.sendCommissionCommandRaw(t, gateway1ID.String(), "not-a-uuid", creds.JWT)
		if res.Success {
			t.Fatalf("expected failure response, got %+v", res)
		}
		if !strings.Contains(res.Message, "Formato del comando incorretto") {
			t.Fatalf("expected format error, got: %q", res.Message)
		}
	})

	t.Run("gatewayId invalido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		creds := fx.gateway1Creds(t)
		res := fx.sendCommissionCommandRaw(t, "not-a-uuid", tenant1ID.String(), creds.JWT)
		if res.Success {
			t.Fatalf("expected failure response, got %+v", res)
		}
		if !strings.Contains(res.Message, "Formato del comando incorretto") {
			t.Fatalf("expected format error, got: %q", res.Message)
		}
	})

	t.Run("gatewayId inesistente", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		creds := fx.gateway1Creds(t)
		res := fx.sendCommissionCommand(t, uuid.New(), tenant1ID, creds.JWT)
		if res.Success {
			t.Fatalf("expected failure response, got %+v", res)
		}
		if !strings.Contains(res.Message, "nessun gateway trovato") {
			t.Fatalf("expected not-found error, got: %q", res.Message)
		}
	})

	t.Run("gateway gia commissionato", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{precommissionGateway1: true})
		defer fx.close(t)

		creds := fx.gateway1Creds(t)
		res := fx.sendCommissionCommand(t, gateway1ID, tenant1ID, creds.JWT)
		if res.Success {
			t.Fatalf("expected failure response, got %+v", res)
		}
		if !strings.Contains(res.Message, "gateway gia commissionato") {
			t.Fatalf("expected already commissioned error, got: %q", res.Message)
		}
	})

	t.Run("jwt non valido", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		res := fx.sendCommissionCommand(t, gateway1ID, tenant1ID, "invalid.jwt.token")
		if res.Success {
			t.Fatalf("expected failure response, got %+v", res)
		}
		fx.assertNotCommissionedState(t, gateway1ID)
	})

	t.Run("jwt associato ad un altro tenant", func(t *testing.T) {
		fx := newCommissionFixture(t, commissionFixtureOptions{})
		defer fx.close(t)

		creds := fx.gateway2Creds(t)
		res := fx.sendCommissionCommand(t, gateway1ID, tenant1ID, creds.JWT)
		if res.Success {
			t.Fatalf("expected failure response, got %+v", res)
		}
		fx.assertNotCommissionedState(t, gateway1ID)
	})
}

type commissionFixtureOptions struct {
	precommissionGateway1 bool
}

func newCommissionFixture(t *testing.T, opts commissionFixtureOptions) *commissionFixture {
	t.Helper()

	cleanupSQLiteFiles(t)

	repoRoot := findRepoRoot(t)
	ctx, cancel := context.WithCancel(context.Background())
	logger := zap.NewNop()

	gatewayDb, err := gatewaydatabase.NewGatewayDatabase(ctx)
	if err != nil {
		cancel()
		t.Fatalf("cannot create gateway db: %v", err)
	}

	bufferDb := bufferdatabase.NewBufferDatabase()
	configRepo := configrepositories.NewSQLiteConfigRepository(ctx, gatewayDb)
	saveRepo := sensor.NewSQLiteSaveSensorDataRepository(ctx, bufferDb)
	bufferRepo := buffereddatasender.NewBufferedDataRepository(ctx, bufferDb)

	natsAddress := natsutil.NatsAddress(getenv("NATS_HOST", "nats"))
	natsPort := natsutil.NatsPort(getenvInt("NATS_PORT", 4222))

	caPath := absolutePathInCmd(repoRoot, getenv("CA_PEM_PATH", "ca.pem"))
	baseCreds := firstNonEmpty(getenv("BASE_CREDS_PATH", ""), getenv("GATEWAY_BASE_CREDS_PATH", ""), "base.creds")
	testCreds := getenv("TEST_CREDS_PATH", "admin_test.creds")

	controllerNc := natsutil.NewNATSConnection(
		natsAddress,
		natsPort,
		natsutil.NatsCredsPath(absolutePathInCmd(repoRoot, baseCreds)),
		natsutil.NatsCAPemPath(caPath),
	)

	js, err := natsutil.NewJetStreamContext(controllerNc)
	if err != nil {
		cancel()
		controllerNc.Close()
		t.Fatalf("cannot create jetstream context: %v", err)
	}

	publisherFactory := buffereddatasender.NewNATSDataPublisherFactory(js, controllerNc, natsAddress, natsPort, ctx, natsutil.NatsCAPemPath(caPath))
	service := gatewayservices.NewGatewayManagerService(
		gatewaymanager.NewGatewayWorkers(),
		gatewaymanager.NewSensorWorkers(),
		saveRepo,
		bufferRepo,
		publisherFactory,
		configRepo,
		ctx,
		logger,
	)

	seedGatewaysForCommissioning(t, repoRoot, configRepo)
	if opts.precommissionGateway1 {
		creds := gatewayCredsFromFile(t, repoRoot, "Infrastructure/testCreds/gateway1.creds")
		cmd := &commanddata.CommissionGateway{
			GatewayId:         gateway1ID,
			TenantId:          tenant1ID,
			CommissionedToken: creds.JWT,
		}
		if err := configRepo.CommissionGateway(cmd, domain.Active); err != nil {
			cancel()
			controllerNc.Close()
			t.Fatalf("cannot precommission gateway: %v", err)
		}
	}

	if err := service.LoadGatewayWorkers(); err != nil {
		cancel()
		controllerNc.Close()
		t.Fatalf("cannot load gateway workers: %v", err)
	}

	controller := commandcontrollers.NewNATSCommissionGatewayController(
		controllerNc,
		commandcontrollers.CommissionGatewaySubject(commissionGatewaySubject),
		service,
		logger,
	)
	controller.Listen()

	if err := controllerNc.FlushTimeout(2 * time.Second); err != nil {
		cancel()
		controllerNc.Close()
		t.Fatalf("cannot flush commission subscription: %v", err)
	}

	publisherNc := natsutil.NewNATSConnection(
		natsAddress,
		natsPort,
		natsutil.NatsCredsPath(absolutePathInCmd(repoRoot, testCreds)),
		natsutil.NatsCAPemPath(caPath),
	)

	observerNc := natsutil.NewNATSConnection(
		natsAddress,
		natsPort,
		natsutil.NatsCredsPath(absolutePathInCmd(repoRoot, testCreds)),
		natsutil.NatsCAPemPath(caPath),
	)

	return &commissionFixture{
		ctx:          ctx,
		cancel:       cancel,
		gatewayDb:    gatewayDb,
		bufferDb:     bufferDb,
		configRepo:   configRepo,
		service:      service,
		controllerNc: controllerNc,
		publisherNc:  publisherNc,
		observerNc:   observerNc,
		repoRoot:     repoRoot,
		caPath:       caPath,
	}
}

func (f *commissionFixture) close(t *testing.T) {
	t.Helper()

	f.cancel()

	if f.publisherNc != nil {
		f.publisherNc.Close()
	}
	if f.observerNc != nil {
		f.observerNc.Close()
	}
	if f.controllerNc != nil {
		f.controllerNc.Close()
	}

	if f.gatewayDb != nil && f.gatewayDb.DB != nil {
		if err := f.gatewayDb.Close(); err != nil {
			t.Fatalf("cannot close gateway db: %v", err)
		}
	}
	if f.bufferDb.DB != nil {
		if err := f.bufferDb.Close(); err != nil {
			t.Fatalf("cannot close buffer db: %v", err)
		}
	}

	cleanupSQLiteFiles(t)
}

func (f *commissionFixture) gateway1Creds(t *testing.T) natsCreds {
	return gatewayCredsFromFile(t, f.repoRoot, "Infrastructure/testCreds/gateway1.creds")
}

func (f *commissionFixture) gateway2Creds(t *testing.T) natsCreds {
	return gatewayCredsFromFile(t, f.repoRoot, "Infrastructure/testCreds/gateway2.creds")
}

func (f *commissionFixture) adminCreds(t *testing.T) natsCreds {
	credsPath := absolutePathInCmd(f.repoRoot, getenv("TEST_CREDS_PATH", "admin_test.creds"))
	return readNatsCreds(t, credsPath)
}

func (f *commissionFixture) sendCommissionCommand(t *testing.T, gatewayID uuid.UUID, tenantID uuid.UUID, token string) gatewayservices.Response {
	return f.sendCommissionCommandRaw(t, gatewayID.String(), tenantID.String(), token)
}

func (f *commissionFixture) sendCommissionCommandRaw(t *testing.T, gatewayID, tenantID, token string) gatewayservices.Response {
	t.Helper()

	reqBody := commissionGatewayRequest{
		GatewayId:         gatewayID,
		TenantId:          tenantID,
		CommissionedToken: token,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("cannot marshal commission payload: %v", err)
	}

	msg, err := f.publisherNc.Request(commissionGatewaySubject, payload, 3*time.Second)
	if err != nil {
		t.Fatalf("cannot publish commission command: %v", err)
	}

	var res gatewayservices.Response
	if err := json.Unmarshal(msg.Data, &res); err != nil {
		t.Fatalf("cannot unmarshal commission response: %v", err)
	}

	return res
}

func (f *commissionFixture) insertBufferedData(t *testing.T, gatewayID, sensorID uuid.UUID) {
	t.Helper()

	profile := sensorprofiles.NewHeartRateProfile(sensorID, sensorprofiles.NewRand())
	data := profile.Generate()
	if err := sensor.NewSQLiteSaveSensorDataRepository(f.ctx, f.bufferDb).Save(data, gatewayID); err != nil {
		t.Fatalf("cannot insert buffered data: %v", err)
	}
}

func (f *commissionFixture) assertCommissionedState(t *testing.T, gatewayID, tenantID uuid.UUID, token string) {
	t.Helper()

	var gotTenant, gotToken, gotStatus string
	err := f.gatewayDb.QueryRowContext(
		f.ctx,
		`SELECT tenantId, token, status FROM gateways WHERE id = ?`,
		gatewayID.String(),
	).Scan(&gotTenant, &gotToken, &gotStatus)
	if err != nil {
		t.Fatalf("cannot read gateway row: %v", err)
	}

	if gotTenant != tenantID.String() || gotToken != token || gotStatus != string(domain.Active) {
		t.Fatalf("unexpected commissioned gateway row")
	}
}

func (f *commissionFixture) assertNotCommissionedState(t *testing.T, gatewayID uuid.UUID) {
	t.Helper()

	var gotTenant, gotToken *string
	var gotStatus string
	err := f.gatewayDb.QueryRowContext(
		f.ctx,
		`SELECT tenantId, token, status FROM gateways WHERE id = ?`,
		gatewayID.String(),
	).Scan(&gotTenant, &gotToken, &gotStatus)
	if err != nil {
		t.Fatalf("cannot read gateway row: %v", err)
	}

	if gotTenant != nil || gotToken != nil || gotStatus != string(domain.Decommissioned) {
		t.Fatalf("expected gateway to remain decommissioned")
	}
}

func seedGatewaysForCommissioning(t *testing.T, repoRoot string, repo *configrepositories.SQLiteConfigRepository) {
	creds1 := gatewayCredsFromFile(t, repoRoot, "Infrastructure/testCreds/gateway1.creds")
	creds2 := gatewayCredsFromFile(t, repoRoot, "Infrastructure/testCreds/gateway2.creds")
	adminCredsPath := absolutePathInCmd(repoRoot, getenv("TEST_CREDS_PATH", "admin_test.creds"))
	adminCreds := readNatsCreds(t, adminCredsPath)

	createGatewayWithCreds(t, repo, gateway1ID, creds1.Public, creds1.Seed)
	createGatewayWithCreds(t, repo, gateway2ID, creds2.Public, creds2.Seed)
	createGatewayWithCreds(t, repo, adminGatewayID, adminCreds.Public, adminCreds.Seed)
}

func gatewayCredsFromFile(t *testing.T, repoRoot, relPath string) natsCreds {
	credsPath := absolutePathInInfrastructure(repoRoot, strings.TrimPrefix(relPath, "Infrastructure/"))
	return readNatsCreds(t, credsPath)
}

func createGatewayWithCreds(t *testing.T, repo *configrepositories.SQLiteConfigRepository, gatewayID uuid.UUID, publicKey, seed string) {
	creds := &credentialsgenerator.Credentials{PublicIdentifier: publicKey, SecretKey: seed}
	cmd := &commanddata.CreateGateway{GatewayId: gatewayID, Interval: 25 * time.Millisecond}
	if err := repo.CreateGateway(cmd, creds, domain.Decommissioned); err != nil {
		t.Fatalf("cannot seed gateway %s: %v", gatewayID, err)
	}
}

func sensorSubject(tenantID, gatewayID, sensorID uuid.UUID) string {
	return "sensor." + tenantID.String() + "." + gatewayID.String() + "." + sensorID.String()
}

func mustParseUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}
	return id
}
