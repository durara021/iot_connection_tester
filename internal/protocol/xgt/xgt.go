package protocol

import (
	"fmt"
	"iot_connection_tester/internal/common/errs"
	"iot_connection_tester/internal/connection"
)

// XGT 프로토콜 구조체
type XGT struct {
	Conn connection.IOConnection // XGT 프로토콜 통신을 위한 IOConnection 인터페이스
}

// XGT 인스턴스 생성자
// @param conn: IOConnection 인터페이스 (TCP/UDP 등 연결 객체)
// @return *XGT: XGT 프로토콜 인스턴스
func NewXGT(conn connection.IOConnection) *XGT {
	return &XGT{Conn: conn}
}

// XGT 통신 요청 및 응답 처리
// @method Transceive: 요청 패킷 전송 및 응답 수신
// @param packet: XGT 요청 바이트 배열
// @param int(cnt) * 2: 수신하고자 하는 응답 데이터의 예상 길이 (바이트)
// @return []byte: XGT 응답 데이터 (int(cnt) * 2만큼 반환)
// @return error: 응답 에러 또는 길이 부족 시 에러 반환
func (x *XGT) Transceive(register byte, strAddr uint16, cnt uint16) ([]byte, error) {

	packet, err := buildLSReadPacket(fmt.Sprintf("%c%05d", register, strAddr), cnt)
	fmt.Printf("전송 패킷: % X\n", packet)

	// 요청 패킷 전송
	if err := x.Conn.Send(packet); err != nil {
		return nil, errs.NewErrs("", "", errs.ErrCodeWriteFailed, err)
	}

	// 응답 수신
	response, err := x.Conn.Receive()
	if err != nil {
		return nil, errs.NewErrs("", "", errs.ErrCodeReadFailed, err)
	}

	// XGT 응답 코드 검사
	if len(response) < 10 || response[8] != 0x00 || response[9] != 0x00 {
		return nil, fmt.Errorf("PLC 응답 오류: 완료 코드 = %02X%02X", response[9], response[8])
	}

	// 응답 길이 확인
	if len(response) < int(cnt)*2 {
		return nil, fmt.Errorf("응답 길이 부족: expected %d, got %d", int(cnt)*2, len(response))
	}

	// 유효한 데이터 반환
	return response[len(response)-int(cnt)*2:], nil
}
