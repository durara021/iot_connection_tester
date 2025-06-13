package ls

import (
	"iot_connection_tester/internal/connection"
	"iot_connection_tester/internal/device/plc/parser"
	"iot_connection_tester/internal/protocol"
	"iot_connection_tester/internal/setting"
	"strconv"
)

// LS PLC 장비 구조체
// @field Conn: TCP 연결 객체
// @field Protocol: XGT 프로토콜 처리 객체
// @field Config: 장비 설정 정보
type LS struct {
	Conn     connection.Connection // TCP 연결 객체
	Protocol *protocol.XGT         // XGT 프로토콜 처리 객체
	Config   setting.DeviceConfig  // 장비 설정 정보
}

// LS 인스턴스 생성
// @param cfg: 장비 설정 정보
// @return *LS: LS 인스턴스
func NewLS(cfg setting.DeviceConfig) *LS {
	conn := connection.NewTCPConnection(cfg.Address)
	return &LS{
		Conn:     conn,
		Protocol: protocol.NewXGT(conn),
	}
}

// TCP 연결 수행
// @return error: 연결 실패 시 에러 반환
func (l *LS) Connect() error {
	return l.Conn.Connect()
}

// PLC 레지스터 데이터 읽기
// @return []byte: 수신한 레지스터 데이터
// @return error: 통신 실패 시 에러 반환
func (l *LS) ReadRegister() ([]byte, error) {
	cfg := l.Config.Setting
	cnt := cfg[0].Address - cfg[len(cfg)-1].Address + 1
	packet := buildLSReadPacket(string(cfg[0].Register)+strconv.Itoa(int(cfg[0].Address)), cnt)
	return l.Protocol.Transceive(packet, int(cnt)*2)
}

// 연결 테스트 및 메모리 데이터 파싱
// @return map[string]uint16: 파싱된 결과 맵
// @return error: 실패 시 에러 반환
func (l *LS) Test() (map[string]uint16, error) {
	memoryData, err := l.ReadRegister()
	if err != nil {
		panic(err)
	}
	result, err := parser.ParseData(memoryData, l.Config.Setting)
	if err != nil {
		panic(err)
	}
	return result, nil
}

// 연결 종료 처리
// @return error: 닫기 실패 시 에러 반환
func (l *LS) Close() error {
	return l.Conn.Close()
}
