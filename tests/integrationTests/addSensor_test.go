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

const addSensorSubject = "commands.addsensor"

type integrationFixture struct {
	ctx          context.Context
	cancel       context.CancelFunc
	gatewayDb    *configrepositories.ConfigDbConnection
	bufferDb     sensor.BufferDbConnection
	configRepo   *configrepositories.SQLiteConfigRepository
	service      *gatewayservices.GatewayManagerService
	controllerNc *nats.Conn
	publisherNc  *nats.Conn
	gatewayID    uuid.UUID
	repoRoot     string
}

type addSensorRequest struct {
	GatewayId string `json:"gatewayId"`
	SensorId  string `json:"sensorId"`
	Profile   string `json:"profile"`
	Interval  int    `json:"interval"`
}

func TestNATSAddSensorIntegration(t *testing.T) {
	t.Run("creazione corretta", func(t *testing.T) {
		fx := newIntegrationFixture(t, fixtureOptions{createGateway: true})
		defer fx.close(t)

		sensorID := uuid.New()
		res := fx.sendAddSensorCommand(t, fx.gatewayID, sensorID, "heart_rate", 25)

		if !res.Success {
			t.Fatalf("expected success response, got %+v", res)
		}
		if res.Message != "Sensore aggiunto con successo" {
			t.Fatalf("unexpected success message: %q", res.Message)
		}

		if count := fx.countSensorsInGateway(t, fx.gatewayID, sensorID); count != 1 {
			t.Fatalf("expected one persisted sensor row, got %d", count)
		}

		fx.waitForBufferData(t, fx.gatewayID)

		duplicate := fx.sendAddSensorCommand(t, fx.gatewayID, sensorID, "heart_rate", 25)
		if duplicate.Success {
			t.Fatalf("expected duplicate creation to fail, got %+v", duplicate)
		}
		if !strings.Contains(duplicate.Message, "già presente") {
			t.Fatalf("expected service-state duplicate error, got: %q", duplicate.Message)
		}
	})

	t.Run("creazione con gatewayId non esistente", func(t *testing.T) {
		fx := newIntegrationFixture(t, fixtureOptions{createGateway: true})
		defer fx.close(t)

		res := fx.sendAddSensorCommand(t, uuid.New(), uuid.New(), "ecg_custom", 10)
		if res.Success {
			t.Fatalf("expected failure response, got %+v", res)
		}
		if !strings.Contains(res.Message, "non trovato") {
			t.Fatalf("expected not-found message, got: %q", res.Message)
		}
	})

	t.Run("creazione con sensorId esistente nel gateway", func(t *testing.T) {
		existingSensorID := uuid.New()
		fx := newIntegrationFixture(t, fixtureOptions{
			createGateway:     true,
			preloadSensor:     true,
			preloadSensorID:   existingSensorID,
			preloadProfile:    "pulse_oximeter",
			preloadIntervalMs: 40,
		})
		defer fx.close(t)

		res := fx.sendAddSensorCommand(t, fx.gatewayID, existingSensorID, "pulse_oximeter", 40)
		if res.Success {
			t.Fatalf("expected duplicate sensor failure, got %+v", res)
		}
		if !strings.Contains(res.Message, "già presente") {
			t.Fatalf("expected duplicate sensor message, got: %q", res.Message)
		}
	})

	t.Run("creazione con interval minore di 1", func(t *testing.T) {
		fx := newIntegrationFixture(t, fixtureOptions{createGateway: true})
		defer fx.close(t)

		res := fx.sendAddSensorCommand(t, fx.gatewayID, uuid.New(), "environmental_sensing", 0)
		if res.Success {
			t.Fatalf("expected parse failure, got %+v", res)
		}
		if !strings.Contains(res.Message, "Formato del comando incorretto") {
			t.Fatalf("expected command format error, got: %q", res.Message)
		}
		if !strings.Contains(res.Message, "intervallo non valido") {
			t.Fatalf("expected interval validation error, got: %q", res.Message)
		}
	})

	t.Run("creazione con profilo non valido", func(t *testing.T) {
		fx := newIntegrationFixture(t, fixtureOptions{createGateway: true})
		defer fx.close(t)

		res := fx.sendAddSensorCommand(t, fx.gatewayID, uuid.New(), "WrongProfile", 20)
		if res.Success {
			t.Fatalf("expected parse failure, got %+v", res)
		}
		if !strings.Contains(res.Message, "Formato del comando incorretto") {
			t.Fatalf("expected command format error, got: %q", res.Message)
		}
		if !strings.Contains(res.Message, "profilo sensore non valido") {
			t.Fatalf("expected profile validation error, got: %q", res.Message)
		}
	})
}

type fixtureOptions struct {
	createGateway     bool
	preloadSensor     bool
	preloadSensorID   uuid.UUID
	preloadProfile    string
	preloadIntervalMs int
}

func newIntegrationFixture(t *testing.T, opts fixtureOptions) *integrationFixture {
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

	gatewayID := uuid.New()
	if opts.createGateway {
		creds, err := credentialsgenerator.GenerateCredentials()
		if err != nil {
			cancel()
			controllerNc.Close()
			t.Fatalf("cannot generate gateway credentials: %v", err)
		}
		createCmd := &commanddata.CreateGateway{GatewayId: gatewayID, Interval: 25 * time.Millisecond}
		if err := configRepo.CreateGateway(createCmd, creds, domain.Decommissioned); err != nil {
			cancel()
			controllerNc.Close()
			t.Fatalf("cannot seed gateway row: %v", err)
		}
	}

	if opts.preloadSensor {
		profile := sensorprofiles.ParseSensorProfile(sensorId, opts.preloadProfile, sensorprofiles.NewRand())
		if profile == nil {
			cancel()
			controllerNc.Close()
			t.Fatalf("invalid preload profile: %s", opts.preloadProfile)
		}
		addCmd := &commanddata.AddSensor{
			GatewayId: gatewayID,
			SensorId:  opts.preloadSensorID,
			Profile:   profile,
			Interval:  time.Duration(opts.preloadIntervalMs) * time.Millisecond,
		}
		if err := configRepo.AddSensor(addCmd, sensor.Active); err != nil {
			cancel()
			controllerNc.Close()
			t.Fatalf("cannot seed sensor row: %v", err)
		}
	}

	if err := service.LoadGatewayWorkers(); err != nil {
		cancel()
		controllerNc.Close()
		t.Fatalf("cannot load gateway workers: %v", err)
	}

	controller := commandcontrollers.NewNATSAddSensorController(
		controllerNc,
		commandcontrollers.AddSensorSubject(addSensorSubject),
		service,
		sensorprofiles.NewRand(),
		logger,
	)
	controller.Listen()

	if err := controllerNc.FlushTimeout(2 * time.Second); err != nil {
		cancel()
		controllerNc.Close()
		t.Fatalf("cannot flush addSensor subscription: %v", err)
	}

	publisherNc := natsutil.NewNATSConnection(
		natsAddress,
		natsPort,
		natsutil.NatsCredsPath(absolutePathInCmd(repoRoot, testCreds)),
		natsutil.NatsCAPemPath(caPath),
	)

	return &integrationFixture{
		ctx:          ctx,
		cancel:       cancel,
		gatewayDb:    gatewayDb,
		bufferDb:     bufferDb,
		configRepo:   configRepo,
		service:      service,
		controllerNc: controllerNc,
		publisherNc:  publisherNc,
		gatewayID:    gatewayID,
		repoRoot:     repoRoot,
	}
}

func (f *integrationFixture) close(t *testing.T) {
	t.Helper()

	f.cancel()

	if f.publisherNc != nil {
		f.publisherNc.Close()
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

func (f *integrationFixture) sendAddSensorCommand(t *testing.T, gatewayID uuid.UUID, sensorID uuid.UUID, profile string, interval int) gatewayservices.Response {
	t.Helper()

	reqBody := addSensorRequest{
		GatewayId: gatewayID.String(),
		SensorId:  sensorID.String(),
		Profile:   profile,
		Interval:  interval,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("cannot marshal addSensor payload: %v", err)
	}

	msg, err := f.publisherNc.Request(addSensorSubject, payload, 3*time.Second)
	if err != nil {
		t.Fatalf("cannot publish addSensor command: %v", err)
	}

	var res gatewayservices.Response
	if err := json.Unmarshal(msg.Data, &res); err != nil {
		t.Fatalf("cannot unmarshal addSensor response: %v", err)
	}

	return res
}

func (f *integrationFixture) countSensorsInGateway(t *testing.T, gatewayID uuid.UUID, sensorID uuid.UUID) int {
	t.Helper()

	var count int
	err := f.gatewayDb.QueryRowContext(
		f.ctx,
		`SELECT COUNT(*) FROM sensors WHERE gatewayId = ? AND id = ?`,
		gatewayID.String(),
		sensorID.String(),
	).Scan(&count)
	if err != nil {
		t.Fatalf("cannot count sensors in gateway: %v", err)
	}
	return count
}

func (f *integrationFixture) waitForBufferData(t *testing.T, gatewayID uuid.UUID) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		var count int
		err := f.bufferDb.QueryRowContext(
			f.ctx,
			`SELECT COUNT(*) FROM buffer WHERE gatewayId = ?`,
			gatewayID.String(),
		).Scan(&count)
		if err == nil && count > 0 {
			return
		}
		time.Sleep(40 * time.Millisecond)
	}
	t.Fatalf("expected sensor goroutine to persist at least one buffered data row for gateway %s", gatewayID)
}
