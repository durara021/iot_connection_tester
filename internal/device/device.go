package device

import (
	"fmt"
	"iot_connection_tester/internal/device/cnc/fanac"
	"iot_connection_tester/internal/device/plc/ls"
	"iot_connection_tester/internal/device/plc/melsec"
	"iot_connection_tester/internal/setting"
	"strings"
)

// 공통 장치 인터페이스
type Device interface {
	Connect() error                   // 장치와의 통신 연결 시도
	Test() (map[string]uint16, error) // 설정된 주소 레지스터 값을 읽고 파싱
	Close() error                     // 장치 연결 종료
}

// 장치 생성 팩토리 함수
// @param cfg: 장치 설정 정보
// @return Device: 설정에 맞는 장치 인스턴스 (Melsec, LS, Fanuc 등)
// @return error: 미지원 브랜드일 경우 에러 반환
var DeviceFactory = func(cfg setting.DeviceConfig) (Device, error) {
	brand := strings.ToLower(cfg.Device)

	switch brand {
	case "melsec", "mel":
		return melsec.NewMelsec(cfg), nil // Mitsubishi PLC
	case "ls":
		return ls.NewLS(cfg), nil // LS산전 PLC
	case "fanac", "cnc":
		return fanac.NewFanuc(cfg), nil // FANUC CNC
	default:
		return nil, fmt.Errorf("객체 생성 실패: 알 수 없는 brand (address: %s)", cfg.Address)
	}
}

// 장치 생성기
// @param cfg: 장치 설정 정보
// @return Device: 생성된 장치 객체
// @return error: 생성 실패 시 에러 반환
func NewDevice(cfg setting.DeviceConfig) (Device, error) {
	dev, err := DeviceFactory(cfg)
	if err != nil {
		fmt.Println(">> DeviceFactory 리턴 : ", dev, err)
	}
	return dev, err
}
