package parser

import (
	"iot_connection_tester/internal/setting"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseData_Error_InvalidLength(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03} // 길이 = 3 (❌)

	settings := []setting.Setting{
		{Register: 'D', Address: 100, Value: "temp"},
	}

	_, err := ParseData(data, settings)
	assert.Error(t, err, "invalid data length should produce an error")
}

func TestParseData_Error_DataTooShort(t *testing.T) {
	data := []byte{0x01, 0x02}

	settings := []setting.Setting{
		{Register: 'D', Address: 100, Value: "ok"},
		{Register: 'D', Address: 101, Value: "fail"},
		{Register: 'D', Address: 102, Value: "skip"}, // offset 2 → ❌ 범위 초과
	}

	_, err := ParseData(data, settings)
	assert.Error(t, err, "data too short for settings should produce an error")
}

func TestParseData_Error_EmptySettings(t *testing.T) {
	data := []byte{0x01, 0x02}

	var settings []setting.Setting

	_, err := ParseData(data, settings)
	assert.Error(t, err, "empty settings should produce an error")
}
