package setting

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// 단일 태그 설정 정보
// @field Register: 레지스터 종류 (예: 'D', 'M', 'X' 등)
// @field Address: 시작 주소 (예: D100 → 100)
// @field Value: 태그명 (예: "Speed", "Temp" 등)
type Setting struct {
	Register byte
	Address  uint16
	Value    string
}

// 장치 전체 설정 구조체
// @field Device: 장치 브랜드 (예: "melsec", "ls", "fanac")
// @field Address: PLC IP:Port 주소
// @field Setting: 읽을 주소 및 태그 설정 목록
type DeviceConfig struct {
	Device  string
	Address string
	Setting []Setting
}

// 문자열로부터 장치 설정 정보 파싱
// @param input: 예) "melsec,192.168.0.1:5000,D100=Speed,D101=Temp"
// @return DeviceConfig: 파싱된 설정 구조체
// @return error: 형식 오류 등 발생 시 에러 반환
func ParseDeviceConfig(input string) (DeviceConfig, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return DeviceConfig{}, fmt.Errorf("input is empty")
	}

	parts := strings.Split(input, ",")
	cfg := DeviceConfig{}
	if len(parts) == 0 {
		return cfg, fmt.Errorf("input is empty")
	}

	// 첫 번째: 장치 종류 (브랜드)
	cfg.Device = parts[0]

	// 두 번째: IP:Port 주소
	cfg.Address = parts[1]

	// 세 번째 이후: 레지스터=태그명
	for _, p := range parts[2:] {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			return cfg, fmt.Errorf("invalid key=value format: %s", p)
		}

		// 주소 파싱: 예 "D100" → 100
		Address, err := strconv.Atoi(kv[0][1:])
		if err != nil {
			return cfg, fmt.Errorf("invalid address format: %s", kv[0])
		}

		cfg.Setting = append(cfg.Setting, Setting{
			Register: kv[0][0],        // 'D' 등
			Address:  uint16(Address), // 주소값
			Value:    kv[1],           // 태그명
		})
	}

	// 주소 기준 정렬
	if len(cfg.Setting) != 0 {
		sort.Slice(cfg.Setting, func(i, j int) bool {
			return int(cfg.Setting[i].Address) < int(cfg.Setting[j].Address)
		})
	}
	return cfg, nil
}
