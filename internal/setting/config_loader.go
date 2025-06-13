package setting

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// 단일 태그 설정 정보
type Setting struct {
	Register byte
	Address  uint16
	Value    string
}

// 장치 전체 설정 구조체
type DeviceConfig struct {
	Device  string
	Address string
	Setting []Setting
}

// JSON 파일 포맷 구조
type JsonSettingItem struct {
	Address uint16 `json:"address"`
	Name    string `json:"name"`
}

type JsonConfig struct {
	Register string            `json:"Register"`
	Settings []JsonSettingItem `json:"Settings"`
}

// 문자열로부터 장치 설정 정보 파싱
func ParseDeviceConfig(input string) (DeviceConfig, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return DeviceConfig{}, fmt.Errorf("IP 또는 IP:PORT 값이 입력되지 않았습니다.")
	}

	ipPort := strings.Split(input, ":")
	ip := ipPort[0]

	// 1. IP 유효성 검사
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return DeviceConfig{}, fmt.Errorf("IP 형식이 잘못되었습니다: %s", ip)
	}
	for _, part := range octets {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 || n > 255 {
			return DeviceConfig{}, fmt.Errorf("IP 숫자 형식 오류 또는 범위 초과: %s", part)
		}
	}

	var device string
	if len(ipPort) == 2 {
		port, err := strconv.Atoi(ipPort[1])
		if err != nil {
			fmt.Errorf("PORT 숫자 형식 오류: %s", ipPort[1])
		}
		switch port {
		case 2004:
			device = "LS"
		case 8193:
			device = "FANAC"
		case 0:
			device = "CNC"
		default:
			device = "MELSEC"
		}
	} else if len(ipPort) == 1 {
		device = "CNC"
	} else {
		return DeviceConfig{}, fmt.Errorf("잘못된 IP:PORT 형식입니다. 예: 192.168.0.1:8193")
	}

	cfg := DeviceConfig{
		Address: input,
		Device:  device,
		Setting: nil,
	}

	// FANAC과 CNC는 config.json 생략
	if cfg.Device != "FANAC" && cfg.Device != "CNC" {
		data, err := os.ReadFile("config.json")
		if err != nil {
			return cfg, fmt.Errorf("config.json 읽기 실패: %w", err)
		}

		var jsonCfg JsonConfig
		if err := json.Unmarshal(data, &jsonCfg); err != nil {
			return cfg, fmt.Errorf("config.json 파싱 오류: %w", err)
		}

		if len(jsonCfg.Register) == 0 {
			return cfg, fmt.Errorf("Register 값이 비어있습니다.")
		}
		registerByte := jsonCfg.Register[0]

		for _, s := range jsonCfg.Settings {
			cfg.Setting = append(cfg.Setting, Setting{
				Register: registerByte,
				Address:  s.Address,
				Value:    s.Name,
			})
		}

		sort.Slice(cfg.Setting, func(i, j int) bool {
			return cfg.Setting[i].Address < cfg.Setting[j].Address
		})
		fmt.Println(cfg.Setting)
	}

	return cfg, nil
}
