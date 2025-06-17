package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// buildLSReadPacket generates a valid LSIS FEnet Read command packet
func buildLSReadPacket(strAddr string, cnt uint16) ([]byte, error) {
	if len(strAddr) != 6 {
		return nil, fmt.Errorf("주소는 6자리 ASCII 문자열이어야 합니다 (예: \"D00010\")")
	}

	var buf bytes.Buffer

	// [1] 고정 헤더 "LSIS"
	buf.Write([]byte{'L', 'S', 'I', 'S'})

	// [2] Reserved (2 bytes)
	buf.Write([]byte{0x00, 0x00})

	// [3] Header Length (2 bytes): 12
	binary.Write(&buf, binary.LittleEndian, uint16(12))

	// [4] Data Length (2 bytes): 주소(6) + 개수(2) = 8
	binary.Write(&buf, binary.LittleEndian, uint16(8))

	// [5] 제어 정보 (4 bytes)
	buf.Write([]byte{
		0x00, // CPU 정보
		0x00, // Reserved
		0x00, // Source of request
		0x54, // 명령어 코드: 0x54 = Read
	})

	// [6] 주소 (6 bytes, ASCII 고정)
	buf.WriteString(strAddr)

	// [7] 읽을 개수 (2 bytes, Little Endian)
	binary.Write(&buf, binary.LittleEndian, cnt)

	return buf.Bytes(), nil
}
