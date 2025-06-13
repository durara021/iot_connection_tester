package connection

/*
#cgo windows CFLAGS: -I${SRCDIR}/../../external/fwlib64/30i
#cgo windows LDFLAGS: -L${SRCDIR}/../../external/fwlib64 -lFwlib64
#include <stdlib.h>
#include <Fwlib64.h>
*/
import "C"
import (
	"fmt"
	"iot_connection_tester/internal/common/errs"
	"unsafe"
)

// FOCAS 라이브러리 핸들 구조체
type CNCHandle struct {
	handle C.ushort
}

// cnc_allclibhndl3 API를 호출, 핸들 생성
// @param addr: IP 주소 (예: "192.168.1.28")
// @return *CNCHandle: 핸들 구조체 포인터
// @return error: 연결 실패 시 반환되는 에러
func NewCNCHandle(addr string) (*CNCHandle, error) {
	ip := C.CString(addr)
	defer C.free(unsafe.Pointer(ip))

	var handle C.ushort
	ret := C.cnc_allclibhndl3(ip, 8193, 10, &handle)
	if ret != 0 {
		return nil, errs.NewErrs(errs.DeviceTypeFanuc, "", errs.ErrCodeConnectionFailed, fmt.Errorf("핸들 할당 실패"))
	}
	return &CNCHandle{handle: handle}, nil
}

// 인터페이스를 위한 동작없는 베서드
func (c *CNCHandle) Connect() error {
	return nil
}

// cnc_freelibhndl함수를 사용 FOCAS 핸들을 해제/리소스 반환
// @return error: 핸들 해제 실패 시 에러 (현재는 nil 반환, 추후 개선 여지 있음)
func (c *CNCHandle) Close() error {
	ret := C.cnc_freelibhndl(c.handle)
	if ret != 0 {
		fmt.Errorf("클로징 실패! : %d", ret)
		// 추후: return errs.NewErrs(...)
	}
	return nil
}

// 외부에서 FOCAS API 호출 시 사용되는
// 내부 FOCAS 핸들 값을 반환하는 getter 함수
// @return C.ushort: FOCAS 핸들 값
func (c *CNCHandle) Handle() C.ushort {
	return c.handle
}
