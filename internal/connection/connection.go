package connection

// 기본적인 연결 생명주기 관리 인터페이스
type Connection interface {

	// 외부 장비나 서비스와의 연결을 시도
	// @return error: 연결 실패 시 에러 반환, 성공 시 nil
	Connect() error

	// 기존 연결을 종료하고 리소스를 해제
	// @return error: 종료 실패 시 에러 반환, 성공 시 nil
	Close() error
}

// Connection 인터페이스를 확장하여 IO 전송 기능을 포함
// 일반적인 송신/수신/요청-응답 기반 통신을 정의
type IOConnection interface {
	Connection

	// 외부 장비나 시스템에 바이트 스트림 데이터를 전송
	// @param data: 전송할 바이트 슬라이스
	// @return error: 전송 실패 시 에러 반환, 성공 시 nil
	Send([]byte) error

	// 외부 장비나 시스템으로부터 데이터를 수신
	// @return []byte: 수신된 데이터
	// @return error: 수신 실패 시 에러 반환
	Receive() ([]byte, error)
}
