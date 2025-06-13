package connection

/*
#cgo windows CFLAGS: -I../../external/fwlib64/30i
#cgo windows LDFLAGS: -L../../external/fwlib64/30i -lFwlib64
#include <stdlib.h>
#include <fwlib64.h>
*/
import "C"
import (
	"fmt"
	"iot_connection_tester/internal/setting"
	"unsafe"
)

// CNCHandle은 Fanuc 장비의 FOCAS 핸들을 관리하는 구조체입니다.
type CNCHandle struct {
	handle C.ushort
}

// NewCNCHandle은 주소 기반으로 FOCAS 라이브러리 연결을 수행하고 핸들을 반환합니다.
func NewCNCHandle(cfg setting.DeviceConfig) (*CNCHandle, error) {
	ip := C.CString(cfg.Address)
	defer C.free(unsafe.Pointer(ip))

	var handle C.ushort
	ret := C.cnc_allclibhndl3(ip, 8193, 10, &handle)
	if ret != 0 {
		return nil, fmt.Errorf("FOCAS 연결 실패 (code=%d)", ret)
	}
	return &CNCHandle{handle: handle}, nil
}

// Connect는 이미 연결된 상태이므로 구현 생략
func (c *CNCHandle) Connect() error {
	return nil
}

// Close는 핸들을 해제합니다.
func (c *CNCHandle) Close() error {
	C.cnc_freelibhndl(c.handle)
	return nil
}

// Handle은 내부 C 핸들을 반환합니다.
func (c *CNCHandle) Handle() C.ushort {
	return c.handle
}
