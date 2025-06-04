package setting

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Setting struct {
	Key   string
	Value string
}
type DeviceConfig struct {
	Address string
	Setting []Setting
}

func ParseDeviceConfig(input string) (DeviceConfig, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return DeviceConfig{}, fmt.Errorf("input is empty")
	}

	// "," 를 기준으로 분리
	parts := strings.Split(input, ",")

	cfg := DeviceConfig{}
	if len(parts) == 0 {
		return cfg, fmt.Errorf("input is empty")
	}

	// 1. 첫 번째 항목은 IP:Port 형식으로 처리
	cfg.Address = parts[0]

	// 2. 나머지 항목은 key=value로 처리
	for _, p := range parts[1:] {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			return cfg, fmt.Errorf("invalid key=value format: %s", p)
		}
		cfg.Setting = append(cfg.Setting, Setting{
			Key:   kv[0],
			Value: kv[1],
		})
	}

	// 숫자에 따른 정렬
	sort.Slice(cfg.Setting, func(i, j int) bool {
		// "D100" -> 100 으로 변환
		ival, _ := strconv.Atoi(cfg.Setting[i].Key[1:])
		jval, _ := strconv.Atoi(cfg.Setting[j].Key[1:])
		return ival < jval
	})

	return cfg, nil
}
