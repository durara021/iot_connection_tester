package melsec

import (
	"iot_connection_tester/internal/connection"
	"iot_connection_tester/internal/device/plc/parser"
	"iot_connection_tester/internal/protocol"
	"iot_connection_tester/internal/setting"
)

// Melsec 구조체
type Melsec struct {
	Conn     connection.Connection // TCP 연결 객체
	Protocol *protocol.SLMP        // SLMP 프로토콜 처리 객체
	Config   setting.DeviceConfig  // 장비 설정 정보
}

// Melsec 인스턴스 생성
// @param cfg: 장비 설정 정보
// @return *Melsec: Melsec 인스턴스
func NewMelsec(cfg setting.DeviceConfig) *Melsec {
	conn := connection.NewTCPConnection(cfg.Address)
	return &Melsec{
		Conn:     conn,
		Config:   cfg,
		Protocol: protocol.NewSLMP(conn),
	}
}

// TCP 연결 수행
// @return error: 연결 실패 시 에러 반환
func (m *Melsec) Connect() error {
	return m.Conn.Connect()
}

// PLC 레지스터 데이터 읽기
// @return []byte: 수신한 레지스터 데이터
// @return error: 통신 실패 시 에러 반환
func (m *Melsec) ReadRegister() ([]byte, error) {
	cfg := m.Config.Setting
	cnt := cfg[0].Address - cfg[len(cfg)-1].Address + 1
	packet := buildMelsecReadPacket(cfg[0].Register, cfg[0].Address, cnt)
	return m.Protocol.Transceive(packet, int(cnt)*2)
}

// 연결 테스트 및 메모리 데이터 파싱
// @return map[string]uint16: 파싱된 결과 맵
// @return error: 실패 시 에러 반환
func (m *Melsec) Test() (map[string]uint16, error) {
	memoryData, err := m.ReadRegister()
	if err != nil {
		panic(err)
	}
	result, err := parser.ParseData(memoryData, m.Config.Setting)
	if err != nil {
		panic(err)
	}
	return result, nil
}

// 연결 종료 처리
// @return error: 닫기 실패 시 에러 반환
func (m *Melsec) Close() error {
	return m.Conn.Close()
}
