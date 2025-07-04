package protocol

import (
	"fmt"
	"iot_connection_tester/internal/common/errs"
	"iot_connection_tester/internal/connection"
)

// SLMP 프로토콜 구조체
type SLMP struct {
	Conn connection.IOConnection //SLMP 프로토콜 통신을 위한 IOConnection 인터페이스
}

// SLMP 인스턴스 생성자
// @param conn: IOConnection 인터페이스 (TCP/UDP 등 연결 객체)
// @return *SLMP: SLMP 프로토콜 인스턴스
func NewSLMP(conn connection.IOConnection) *SLMP {
	return &SLMP{Conn: conn}
}

// SLMP 통신 요청 및 응답 처리
// @method Transceive: 요청 패킷 전송 및 응답 수신
// @param packet: SLMP 요청 바이트 배열
// @param expectSize: 수신하고자 하는 응답 데이터의 예상 길이 (바이트)
// @return []byte: SLMP 응답 데이터 (expectSize만큼 반환)
// @return error: 응답 에러 또는 길이 부족 시 에러 반환
func (s *SLMP) Transceive(register byte, strAddr uint16, cnt uint16) ([]byte, error) {

	packet := buildMelsecReadPacket(register, strAddr, cnt)
	// 요청 패킷 전송
	if err := s.Conn.Send(packet); err != nil {
		return nil, errs.NewErrs("", "", errs.ErrCodeWriteFailed, err)
	}

	// 응답 수신
	response, err := s.Conn.Receive()
	if err != nil {
		return nil, errs.NewErrs("", "", errs.ErrCodeReadFailed, err)
	}

	// SLMP 응답 코드 검사
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
