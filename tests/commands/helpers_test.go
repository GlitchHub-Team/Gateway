package commandstests

import (
	"testing"
	"time"

	credentialsgenerator "Gateway/internal/credentialsGenerator"
	"Gateway/internal/domain"
	commanddata "Gateway/internal/gatewayManager/commandData"
	sensor "Gateway/internal/sensor"
	profiles "Gateway/internal/sensor/sensorProfiles"

	"github.com/google/uuid"
)

type mockSensorStarter struct {
	started chan struct{}
}

func (m *mockSensorStarter) Start() {
	if m.started != nil {
		select {
		case m.started <- struct{}{}:
		default:
		}
	}
}

type mockSensorStopper struct {
	stopCalls int
}

func (m *mockSensorStopper) Stop() {
	m.stopCalls++
}

type mockSensorInterrupter struct {
	interruptCalls int
}

func (m *mockSensorInterrupter) Interrupt() {
	m.interruptCalls++
}

type mockSensorResumer struct {
	resumeCalls int
}

func (m *mockSensorResumer) Resume() {
	m.resumeCalls++
}

type mockSensorAdder struct {
	err          error
	called       bool
	receivedCmd  *commanddata.AddSensor
	receivedStat sensor.SensorStatus
}

func (m *mockSensorAdder) AddSensor(cmd *commanddata.AddSensor, status sensor.SensorStatus) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedStat = status
	return m.err
}

type mockSensorDeleter struct {
	err         error
	called      bool
	receivedCmd *commanddata.DeleteSensor
}

func (m *mockSensorDeleter) DeleteSensor(cmd *commanddata.DeleteSensor) error {
	m.called = true
	m.receivedCmd = cmd
	return m.err
}

type mockSensorInterrupterPort struct {
	err          error
	called       bool
	receivedCmd  *commanddata.InterruptSensor
	receivedStat sensor.SensorStatus
}

func (m *mockSensorInterrupterPort) InterruptSensor(cmd *commanddata.InterruptSensor, status sensor.SensorStatus) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedStat = status
	return m.err
}

type mockSensorResumerPort struct {
	err          error
	called       bool
	receivedCmd  *commanddata.ResumeSensor
	receivedStat sensor.SensorStatus
}

func (m *mockSensorResumerPort) ResumeSensor(cmd *commanddata.ResumeSensor, status sensor.SensorStatus) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedStat = status
	return m.err
}

type mockGatewayStopper struct {
	stopCalls int
}

func (m *mockGatewayStopper) Stop() {
	m.stopCalls++
}

type mockGatewayInterrupter struct {
	interruptCalls int
}

func (m *mockGatewayInterrupter) Interrupt() {
	m.interruptCalls++
}

type mockGatewayResumer struct {
	resumeCalls int
}

func (m *mockGatewayResumer) Resume() {
	m.resumeCalls++
}

type mockGatewayResetter struct {
	err              error
	resetCalls       int
	receivedInterval time.Duration
}

func (m *mockGatewayResetter) Reset(interval time.Duration) error {
	m.resetCalls++
	m.receivedInterval = interval
	return m.err
}

type mockGatewayDeleter struct {
	err         error
	called      bool
	receivedCmd *commanddata.DeleteGateway
}

func (m *mockGatewayDeleter) DeleteGateway(cmd *commanddata.DeleteGateway) error {
	m.called = true
	m.receivedCmd = cmd
	return m.err
}

type mockGatewayInterrupterPort struct {
	err          error
	called       bool
	receivedCmd  *commanddata.InterruptGateway
	receivedStat domain.GatewayStatus
}

func (m *mockGatewayInterrupterPort) InterruptGateway(cmd *commanddata.InterruptGateway, status domain.GatewayStatus) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedStat = status
	return m.err
}

type mockGatewayResumerPort struct {
	err          error
	called       bool
	receivedCmd  *commanddata.ResumeGateway
	receivedStat domain.GatewayStatus
}

func (m *mockGatewayResumerPort) ResumeGateway(cmd *commanddata.ResumeGateway, status domain.GatewayStatus) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedStat = status
	return m.err
}

type mockGatewayResetterPort struct {
	err              error
	called           bool
	receivedCmd      *commanddata.ResetGateway
	receivedInterval time.Duration
}

func (m *mockGatewayResetterPort) ResetGateway(cmd *commanddata.ResetGateway, interval time.Duration) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedInterval = interval
	return m.err
}

type mockGatewayCreator struct {
	err                 error
	called              bool
	receivedCmd         *commanddata.CreateGateway
	receivedCredentials *credentialsgenerator.Credentials
	receivedStat        domain.GatewayStatus
}

func (m *mockGatewayCreator) CreateGateway(cmd *commanddata.CreateGateway, credentials *credentialsgenerator.Credentials, status domain.GatewayStatus) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedCredentials = credentials
	m.receivedStat = status
	return m.err
}

type mockGatewayGreeter struct {
	err        error
	helloCalls int
}

func (m *mockGatewayGreeter) Hello() error {
	m.helloCalls++
	return m.err
}

type mockGatewayStarter struct {
	started chan struct{}
}

func (m *mockGatewayStarter) Start() {
	if m.started != nil {
		select {
		case m.started <- struct{}{}:
		default:
		}
	}
}

type mockGatewayCommissionerPort struct {
	err          error
	called       bool
	receivedCmd  *commanddata.CommissionGateway
	receivedStat domain.GatewayStatus
}

func (m *mockGatewayCommissionerPort) CommissionGateway(cmd *commanddata.CommissionGateway, status domain.GatewayStatus) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedStat = status
	return m.err
}

type mockGatewayCommissioner struct {
	err              error
	commissionCalls  int
	receivedTenantID uuid.UUID
	receivedToken    string
}

func (m *mockGatewayCommissioner) Commission(tenantID uuid.UUID, commissionedToken string) error {
	m.commissionCalls++
	m.receivedTenantID = tenantID
	m.receivedToken = commissionedToken
	return m.err
}

type mockGatewayDecommissionerPort struct {
	err          error
	called       bool
	receivedCmd  *commanddata.DecommissionGateway
	receivedStat domain.GatewayStatus
}

func (m *mockGatewayDecommissionerPort) DecommissionGateway(cmd *commanddata.DecommissionGateway, status domain.GatewayStatus) error {
	m.called = true
	m.receivedCmd = cmd
	m.receivedStat = status
	return m.err
}

type mockGatewayDecommissioner struct {
	err               error
	decommissionCalls int
}

func (m *mockGatewayDecommissioner) Decommission() error {
	m.decommissionCalls++
	return m.err
}

type mockSerializableData struct{}

func (m *mockSerializableData) Serialize() ([]byte, error) {
	return []byte(`{"mock":true}`), nil
}

func newGeneratedSensorData() *profiles.GeneratedSensorData {
	return &profiles.GeneratedSensorData{
		SensorId:  uuid.New(),
		Timestamp: time.Now(),
		Profile:   "mock-profile",
		Data:      &mockSerializableData{},
	}
}

func waitForSignal(t *testing.T, ch <-chan struct{}, label string) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", label)
	}
}
