package cnc

import "iot_connection_tester/internal/device"

// 모든 CNC 장비의 표준 동작 정의
type CNC interface {
	device.Device

	// CNC 장비로부터 데이터를 읽어오는 함수입니다.
	// @return []byte: 장비에서 수신한 원시 바이트 데이터
	// @return error: 데이터 읽기 실패 시 에러 반환
	ReadData() ([]byte, error)
}
