package parser

import (
	"iot_connection_tester/internal/setting"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseData_ValidData(t *testing.T) {
	data := []byte{
		0x00, 0x01, // 1
		0x00, 0x02, // 2
		0x00, 0x03, // 3
		0x00, 0x04, // 4
		0x00, 0x05, // 5
		0x00, 0x06, // 6
		0x00, 0x07, // 7
	}

	settings := []setting.Setting{
		{Register: 'D', Address: 100, Name: "ok"},
		{Register: 'D', Address: 104, Name: "fail"},
		{Register: 'D', Address: 106, Name: "skip"}, // offset 2 → ❌ 범위 초과
	}

	result, err := ParseData(data, settings)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 3)
	assert.Equal(t, result["ok"], uint16(1))
	assert.Equal(t, result["fail"], uint16(5))
	assert.Equal(t, result["skip"], uint16(7))
}

func TestParseData_ValidData2(t *testing.T) {
	data := []byte{
		0x00, 0x03, // 3
		0x00, 0x01, // 1
	}

	settings := []setting.Setting{
		{Register: 'D', Address: 100, Name: "ok"},
		{Register: 'D', Address: 104, Name: "fail"},
		{Register: 'D', Address: 106, Name: "skip"}, // offset 2 → ❌ 범위 초과
	}

	result, err := ParseData(data, settings)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result["ok"], uint16(3))
}

func TestParseData_Error_InvalidLength(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03} // 길이 = 3 (❌)

	settings := []setting.Setting{
		{Register: 'D', Address: 100, Name: "temp"},
	}

	_, err := ParseData(data, settings)
	assert.Error(t, err, "invalid data length should produce an error")
}

func TestParseData_Error_EmptySettings(t *testing.T) {
	data := []byte{0x01, 0x02}

	var settings []setting.Setting

	_, err := ParseData(data, settings)
	assert.Error(t, err, "empty settings should produce an error")
}
func TestParseData_Error_EmptyData(t *testing.T) {
	data := []byte{}

	var settings []setting.Setting

	_, err := ParseData(data, settings)
	assert.Error(t, err, "empty data should produce an error")
}
