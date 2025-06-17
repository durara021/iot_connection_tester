package fanac

import (
	"fmt"
	"iot_connection_tester/internal/connection"
	"iot_connection_tester/internal/protocol"
	"iot_connection_tester/internal/setting"
)

// Fanuc는 FANUC CNC 장비와의 통신 구조체
type Fanuc struct {
	handle   *connection.CNCHandle // CNC 핸들
	protocol *protocol.FOCAS       // FOCAS 통신 핸들
	Config   setting.DeviceConfig  // 장비 설정 정보
}

// NewFanuc은 주어진 설정 정보를 기반으로 Fanuc 인스턴스를 생성
// FANUC 장비에 연결하기 위한 핸들을 초기화하고 FOCAS 프로토콜을 설정
func NewFanuc(cfg setting.DeviceConfig) *Fanuc {
	handle, err := connection.NewCNCHandle(cfg.Address)
	if err != nil {
		fmt.Println("FOCAS 연결 실패:", err)
		return nil
	}

	cHandle := uint16(handle.Handle())
	// fmt.Println("handle : ", cHandle)
	proto := protocol.NewFOCAS(cHandle)

	return &Fanuc{
		handle:   handle,
		protocol: proto,
		Config:   cfg,
	}
}

// Fanuc 장비 연결 초기화
// 현재 이미 연결된 상태이므로 별도의 동작 없이 성공 반환
func (f *Fanuc) Connect() error {
	return nil
}

// Test는 FANUC 장비로부터 CNC 모드, 상태, 프로그램 번호, M/T 코드 정보 추출
// @return map[string]uint16: 수집된 주요 CNC 상태 정보
// @return error: 통신 실패 시 에러 반환
func (f *Fanuc) Test() (map[string]uint16, error) {
	result := make(map[string]uint16)

	// CNC 모드 및 상태 읽기
	mode, state, err := f.protocol.ReadCNCModeAndStatus()
	if err != nil {
		return nil, fmt.Errorf("CNC 모드/상태 읽기 실패: %w", err)
	}
	result["cnc_mode"] = uint16(mode)
	result["cnc_state"] = uint16(state)

	// 프로그램 번호 읽기
	progNum, err := f.protocol.ReadProgramNumber()
	if err != nil {
		return nil, fmt.Errorf("프로그램 번호 읽기 실패: %w", err)
	} else {
		result["program_number"] = uint16(progNum)
	}

	// M 코드 읽기
	mcode, err := f.protocol.ReadMCode()
	if err != nil {
		return nil, fmt.Errorf("M코드 읽기 실패: %w", err)
	} else {
		result["M_code"] = uint16(mcode)
	}

	// T 코드 읽기
	tcode, err := f.protocol.ReadTCode()
	if err != nil {
		return nil, fmt.Errorf("T코드 읽기 실패: %w", err)
	} else {
		result["T_code"] = uint16(tcode)
	}

	return result, nil
}

// FOCAS핸들해제 / FANUC 장비연결 종료
// @return error: 해제 실패 시 에러 반환
func (f *Fanuc) Close() error {
	if f == nil {
		return nil
	}
	return f.handle.Close()
}
