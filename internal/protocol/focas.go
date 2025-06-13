package protocol

/*
#cgo windows CFLAGS: -I${SRCDIR}/../../external/fwlib64/30i
#cgo windows LDFLAGS: -L${SRCDIR}/../../external/fwlib64 -lFwlib64
#include <stdint.h>
#include <Fwlib64.h>

// M/T 코드 추출용 보조 C 함수
short get_go_code(ODBMDL *mdl) {
    if (mdl == 0) return -1;
    if (mdl->datano == 106 || mdl->datano == 108) {
        return (short)mdl->modal.aux.aux_data;
    }
    return -1;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// FANUC FOCAS 라이브러리 연동 구조체
type FOCAS struct {
	handle C.ushort // 세션 핸들 (C 라이브러리의 ushort형)
}

// FOCAS 객체 생성자
// @param handle: uint16형 FANUC 연결 핸들
// @return *FOCAS: FOCAS 프로토콜 인스턴스
func NewFOCAS(handle uint16) *FOCAS {
	return &FOCAS{handle: C.ushort(handle)}
}

// CNC 자동 운전 모드 및 상태 읽기
// @return mode: 자동 모드 상태 (예: MDI, AUTO 등)
// @return state: 동작 상태 (예: RUN, STOP 등)
// @return error: 읽기 실패 시 에러
func (f *FOCAS) ReadCNCModeAndStatus() (mode int, state int, err error) {
	var stat C.ODBST
	// fmt.Println("handle : ", f.handle)
	ret := C.cnc_statinfo(f.handle, &stat)
	if ret != 0 {
		return 0, 0, fmt.Errorf("CNC 상태 읽기 실패 (code=%d)", ret)
	}
	return int(stat.aut), int(stat.run), nil
}

// 현재 실행 중인 M 코드 읽기
// @return int: M 코드 번호
// @return error: 읽기 실패 시 에러
func (f *FOCAS) ReadMCode() (int, error) {
	var mdl C.ODBMDL
	// fmt.Println("handle : ", f.handle)
	ret := C.cnc_modal(f.handle, 106, 0, &mdl)
	if ret != 0 {
		return 0, fmt.Errorf("M코드 읽기 실패 (code=%d)", ret)
	}
	return int(C.get_go_code(&mdl)), nil
}

// 현재 실행 중인 T 코드 읽기
// @return int: T 코드 번호
// @return error: 읽기 실패 시 에러
func (f *FOCAS) ReadTCode() (int, error) {
	var mdl C.ODBMDL
	// fmt.Println("handle : ", f.handle)
	ret := C.cnc_modal(f.handle, 108, 0, &mdl)
	if ret != 0 {
		return 0, fmt.Errorf("T코드 읽기 실패 (code=%d)", ret)
	}
	return int(C.get_go_code(&mdl)), nil
}

// 현재 실행 중인 프로그램 번호 읽기
// @return int: 프로그램 번호
// @return error: 읽기 실패 시 에러
func (f *FOCAS) ReadProgramNumber() (int, error) {
	var prog C.ODBPRO
	// fmt.Println("handle : ", f.handle)
	ret := C.cnc_rdprgnum(f.handle, &prog)
	if ret != 0 {
		return 0, fmt.Errorf("프로그램 번호 읽기 실패 (code=%d)", ret)
	}
	return int(prog.mdata), nil
}

// FANUC 동적 데이터 구조체 (Go 버전)
// CNC 절대위치, 상대위치 등 실시간 축 상태 포함
type ODBDYGo struct {
	Dummy int16
	Axis  int16
	Alarm int16
	_     int16 // padding

	PrgNum  int32
	PrgmNum int32
	SeqNum  int32
	ActF    int32
	ActS    int32

	Absolute [8]int32 // 절대 위치
	Machine  [8]int32 // 머신 위치
	Relative [8]int32 // 상대 위치
	Distance [8]int32 // 이동 거리
}

// CNC 동적 정보 읽기 (좌표 등)
// @return *ODBDYGo: 동적 데이터 포인터
// @return error: 읽기 실패 시 에러
func (f *FOCAS) ReadDynamicODBDY() (*ODBDYGo, error) {
	var buf C.ODBDY
	// fmt.Println("handle : ", f.handle)
	ret := C.cnc_rddynamic(f.handle, -1, C.short(unsafe.Sizeof(buf)), &buf)
	if ret != 0 {
		return nil, fmt.Errorf("동적 데이터 읽기 실패 (code=%d)", ret)
	}
	return ConvertODBDY(&buf), nil
}

// C 구조체를 Go 구조체로 변환
// @param cbuf: C.ODBDY 포인터
// @return *ODBDYGo: Go 형식 동적 데이터 구조체 포인터
func ConvertODBDY(cbuf *C.ODBDY) *ODBDYGo {
	var result ODBDYGo
	ptr := unsafe.Pointer(cbuf)
	result = *(*ODBDYGo)(ptr)
	return &result
}
