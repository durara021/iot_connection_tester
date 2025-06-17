package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDeviceConfig_ValidInput(t *testing.T) {
	input := "192.168.0.1:5000"

	expectedSettings := []Setting{
		{Register: 68, Address: 100, Name: "heat"},
		{Register: 68, Address: 102, Name: "vibrate"},
	}

	cfg, err := ParseDeviceConfig(input)
	assert.NoError(t, err)
	assert.Equal(t, "MELSEC", cfg.Device)
	assert.Equal(t, "192.168.0.1:5000", cfg.Address)
	assert.Len(t, cfg.Setting, len(expectedSettings))

	for i, setting := range expectedSettings {
		assert.Equal(t, setting, cfg.Setting[i])
	}
}

func TestParseDeviceConfig_NoPort(t *testing.T) {
	input := "192.168.0.1"

	device, err := ParseDeviceConfig(input)
	assert.NoError(t, err)
	assert.Equal(t, "CNC", device.Device)
}

func TestParseDeviceConfig_2004Port(t *testing.T) {
	input := "192.168.0.1:2004"

	device, err := ParseDeviceConfig(input)
	assert.NoError(t, err)
	assert.Equal(t, "LS", device.Device)
}

func TestParseDeviceConfig_EmptyInput(t *testing.T) {
	input := "192.168.0.1"
	_, err := ParseDeviceConfig(input)
	assert.Error(t, err)
}
