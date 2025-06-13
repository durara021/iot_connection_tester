package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDeviceConfig_ValidInput(t *testing.T) {
	input := "cnc,192.168.0.1:5000,D110=temp,D102=press"
	expectedAddress := "192.168.0.1:5000"
	expectedSettings := []Setting{
		{Register: 'D', Address: 102, Value: "press"},
		{Register: 'D', Address: 110, Value: "temp"},
	}

	cfg, err := ParseDeviceConfig(input)
	t.Logf("cfg: %+v", cfg)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, cfg.Address)
	assert.Len(t, cfg.Setting, len(expectedSettings))

	for i, setting := range expectedSettings {
		assert.Equal(t, setting, cfg.Setting[i])
	}
}

func TestParseDeviceConfig_InvalidFormat(t *testing.T) {
	input := "mel,192.168.0.1:5000,D100" // = 빠짐
	_, err := ParseDeviceConfig(input)
	assert.Error(t, err)
}

func TestParseDeviceConfig_EmptyInput(t *testing.T) {
	input := ""
	_, err := ParseDeviceConfig(input)
	assert.Error(t, err)
}

func TestParseDeviceConfig_OnlyAddress(t *testing.T) {
	input := "ls,192.168.0.1:5000"
	cfg, err := ParseDeviceConfig(input)

	assert.NoError(t, err)
	assert.Equal(t, "192.168.0.1:5000", cfg.Address)
	assert.Empty(t, cfg.Setting)
}
