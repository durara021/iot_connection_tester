package plc

import (
	"iot_connection_tester/internal/device"
)

// PLC 공통 인터페이스
type PLC interface {
	device.Device // 공통 장치 기능 포함 (Connect(), Close(), Test() 등)
	// PLC 메모리(레지스터) 데이터를 읽는 함수
	// @return []byte: 읽은 데이터 (2바이트 단위의 바이트 배열)
	// @return error: 통신 실패 등 에러 발생 시 반환
	ReadRegister() ([]byte, error) // 레지스터 데이터 읽기 함수
}
