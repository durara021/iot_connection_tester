package device

import (
	"fmt"
	"iot_connection_tester/internal/device/plc/ls"
	"iot_connection_tester/internal/device/plc/melsec"
	"iot_connection_tester/internal/setting"
	"strings"
)

func New(cfg setting.DeviceConfig) (Device, error) {
	brand := DetectBrand(cfg.Setting)

	switch brand {
	case "melsec":
		return melsec.NewMelsec(cfg.Address), nil
	case "ls":
		return ls.NewLS(cfg.Address), nil
	/**
	case "cnc":
		return cnc.NewCNC(cfg.Address), nil
	*/
	default:
		return nil, fmt.Errorf("객체 생성 실패: 알 수 없는 brand (address: %s)", cfg.Address)
	}
}

func DetectBrand(setting []setting.Setting) string {
	if len(setting) == 0 {
		return "cnc"
	}

	if strings.HasPrefix(setting[0].Key, "D") {
		return "melsec"
	}

	return "ls"
}
