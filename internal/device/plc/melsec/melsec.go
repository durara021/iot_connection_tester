package melsec

import (
	"fmt"
	"net"
	"strconv"
)

type Melsec struct {
	address string
	conn    net.Conn
}

func NewMelsec(address string) *Melsec {
	return &Melsec{
		address: address,
	}
}

func (m *Melsec) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s", m.address))
	if err != nil {
		return err
	}
	m.conn = conn
	return nil
}

func (m *Melsec) ReadMemory(strAddr string, endAddr string) ([]byte, error) {
	strAddrInt, err := strconv.Atoi(strAddr[1:])
	if err != nil {
		return nil, fmt.Errorf("Invalid address: %s", strAddr)
	}

	endAddrInt, err := strconv.Atoi(endAddr[1:])
	if err != nil {
		return nil, fmt.Errorf("Invalid address: %s", endAddr)
	}

	cnt := endAddrInt - strAddrInt + 1
	pocket := buildMelsecReadPacket(uint16(strAddrInt), uint16(cnt))

	_, err = m.conn.Write(pocket)
	if err != nil {
		return nil, fmt.Errorf("Failed to write melsec packet: %s", err)
	}

	buf := make([]byte, 1024)
	n, err := m.conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("Failed to read melsec packet: %s", err)
	}

	if buf[8] != 0x00 || buf[9] != 0x00 {
		return nil, fmt.Errorf("PLC 응답 오류: 완료 코드 = %02X%02X", buf[9], buf[8])
	}

	return buf[n-2*cnt : n], nil
}

func (m *Melsec) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}
