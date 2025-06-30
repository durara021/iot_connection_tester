package connection

import (
	"iot_connection_tester/internal/common/errs"
	"net"
)

// TCPConnection 구조체
type TCPConnection struct {
	Address string
	conn    net.Conn
}

// 새로운 TCPConnection 인스턴스를 생성
// @param address: 대상 IP:PORT 형식의 주소 문자열
// @return *TCPConnection: TCPConnection 포인터 객체
func NewTCPConnection(address string) *TCPConnection {
	return &TCPConnection{Address: address}
}

// TCP 연결
// @return error: 연결 실패 시 에러 반환, 성공 시 nil
func (t *TCPConnection) Connect() error {
	conn, err := net.Dial("tcp", t.Address)
	if err != nil {
		return errs.NewErrs("", "", errs.ErrCodeConnectionFailed, err)
	}
	t.conn = conn

	return nil
}

// 데이터를 전송
// @param data: 전송할 바이트 슬라이스
// @return error: 전송 실패 시 에러 반환, 성공 시 nil
func (t *TCPConnection) Send(data []byte) error {
	_, err := t.conn.Write(data)
	if err != nil {
		return errs.NewErrs("", "", errs.ErrCodeWriteFailed, err)
	}

	return nil
}

// 데이터를 수신
// @return []byte: 수신한 바이트 슬라이스
// @return error: 수신 실패 시 에러 반환
func (t *TCPConnection) Receive() ([]byte, error) {
	buf := make([]byte, 2048)
	n, err := t.conn.Read(buf)
	if err != nil {
		return nil, errs.NewErrs("", "", errs.ErrCodeReadFailed, err)
	}

	return buf[:n], nil
}

// TCP 연결을 종료
// @return error: 닫기 실패 시 에러 반환
// @note: nil인 연결에 대해서는 별도로 닫을 수 없으므로 주의
func (t *TCPConnection) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}
