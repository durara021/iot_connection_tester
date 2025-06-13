package protocol

import (
	"fmt"
	"iot_connection_tester/internal/common/errs"
	"iot_connection_tester/internal/connection"
)

type XGT struct {
	Conn connection.IOConnection
}

func NewXGT(conn connection.IOConnection) *XGT {
	return &XGT{Conn: conn}
}

func (x *XGT) Transceive(packet []byte, expectSize int) ([]byte, error) {
	// fmt.Println("Test Log ------- Transceiving XGT", packet, " -------")
	response, err := x.Conn.Transceive(packet)
	if err != nil {
		return nil, errs.NewErrs("", "", errs.ErrCodeInvalidResponse, err)
	}

	if len(response) < 10 || response[8] != 0x00 || response[9] != 0x00 {
		return nil, fmt.Errorf("PLC 응답 오류: 완료 코드 = %02X%02X", response[9], response[8])
	}

	if len(response) < expectSize {
		return nil, fmt.Errorf("응답 길이 부족: expected %d, got %d", expectSize, len(response))
	}

	return response[len(response)-expectSize:], nil
}
