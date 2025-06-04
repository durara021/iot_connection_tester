package usecase

import (
	"flag"
	"fmt"
	"iot_connection_tester/internal/device"
	"iot_connection_tester/internal/device/cnc"
	"iot_connection_tester/internal/device/plc"
	"iot_connection_tester/internal/parser"
	"iot_connection_tester/internal/setting"
)

func RunTest(input string) error {

	inf := flag.Int("test", 100, "desc")
	flag.Parse()

	cfg, err := setting.ParseDeviceConfig(input)
	if err != nil {
		return fmt.Errorf("Config 파싱 실패: %w", err)
	}

	dev, err := device.New(cfg)
	if err != nil {
		return fmt.Errorf("장비 생성 실패: %w", err)
	}
	defer dev.Close()

	if err := dev.Connect(); err != nil {
		return fmt.Errorf("장비 연결 실패: %w", err)
	}

	var parsingData []uint16
	switch d := dev.(type) {
	case cnc.CNC:
		// TODO: CNC 처리 추가
	case plc.PLC:
		memoryData, err := d.ReadMemory(cfg.Setting[0].Key, cfg.Setting[len(cfg.Setting)-1].Key)
		if err != nil {
			return fmt.Errorf("PLC 메모리 읽기 실패: %w", err)
		}
		parsingData, err = parser.ParseWord(memoryData)
		if err != nil {
			return fmt.Errorf("메모리 파싱 실패: %w", err)
		}
	default:
		return fmt.Errorf("지원하지 않는 장비 타입입니다: %T", d)
	}

	if len(parsingData) == 0 {
		return fmt.Errorf("읽어온 메모리 데이터가 비어있습니다")
	}

	fmt.Println("📦 파싱된 데이터:")
	for i, val := range parsingData {
		fmt.Printf("[%d] = %d\n", i, val)
	}

	return nil
}
