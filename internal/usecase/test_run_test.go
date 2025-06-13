package usecase_test

import (
	"iot_connection_tester/internal/device"
	"iot_connection_tester/internal/device/plc/parser"
	"iot_connection_tester/internal/setting"
	"iot_connection_tester/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockPLC 구현
type MockPLC struct {
	MockData []byte
	Config   setting.DeviceConfig
}

func (m *MockPLC) Connect() error {
	// fmt.Println("Test Log ------- Connecting to ", m.Config.Address, " -------")
	return nil
}
func (m *MockPLC) Close() error {
	// fmt.Println("Test Log ------- Closing connection to ", m.Config.Address, " -------")
	return nil
}
func (m *MockPLC) ReadRegister() ([]byte, error) {
	// fmt.Println("Test Log ------- Sending data to ", m.Config.Address, "-------")
	// fmt.Println("Test Log ------- Receiving data from ", m.Config.Address, " -------")
	return m.MockData, nil
}
func (m *MockPLC) Test() (map[string]uint16, error) {
	// fmt.Println("Test Log ------- Testing -------")
	return parser.ParseData(m.MockData, m.Config.Setting)
}

func mockFactory(cfg setting.DeviceConfig) (device.Device, error) {
	return &MockPLC{
		MockData: []byte{0x03, 0xE8, 0x07, 0xD0}, // 예: 1000, 2000
		Config:   cfg,
	}, nil
}

func TestRunTest_MockedPLC(t *testing.T) {
	original := device.DeviceFactory
	device.DeviceFactory = mockFactory
	defer func() { device.DeviceFactory = original }()

	input := "ls,127.0.0.1:5000,D100=heat,D101=pressure"

	err := usecase.RunTest(input)
	assert.NoError(t, err)
}
