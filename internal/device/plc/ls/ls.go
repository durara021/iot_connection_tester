package ls

import (
	"fmt"
	"net"
	"strconv"
)

type LS struct {
	address string
	conn    net.Conn
}

func NewLS(address string) *LS {
	return &LS{
		address: address,
	}
}

func (l *LS) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s", l.address))
	if err != nil {
		return err
	}
	l.conn = conn
	return nil
}

func (l *LS) ReadMemory(strAddr string, endAddr string) ([]byte, error) {
	packet := buildLSREadPacket(strAddr, uint16(strconv.Atoi(strAddr[1:])-strconv.Atoi(endAddr[1:])+1))

	_, err := l.conn.Write(packet)
	if err != nil {
		return nil, fmt.Errorf("패킷 전송 실패 : %w", err)
	}

	buf := make([]byte, 1024)
	n, err := l.conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("응답 수신 실패 : %w", err)
	}

	return buf[n-2*count : n], nil
}

func (l *LS) Close() error {
	return l.conn.Close()
}
