package xgt

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BuildXGTReadPacket generates a packet to read 1 DWORD value from given register and address
// Example: register='D', strAddr=100 → %DW100
func BuildXGTBlockReadPacket(register byte, strAddr uint16, cnt uint16) ([]byte, error) {
	var buf bytes.Buffer

	// --- HEADER ---
	buf.Write([]byte{'L', 'S', 'I', 'S', '-', 'X', 'G', 'T', 0x00, 0x00}) // "LSIS-XGT"
	buf.Write([]byte{0x00, 0x00})                                         // PLC Info
	buf.WriteByte(0x00)                                                   // CPU Info
	buf.WriteByte(0x33)                                                   // PC → PLC
	buf.Write([]byte{0x00, 0x00})                                         // Invoke ID
	// --- HEADER ---

	// --- BODY ---
	var body bytes.Buffer
	binary.Write(&body, binary.LittleEndian, uint16(0x54)) // 읽기요구
	binary.Write(&body, binary.LittleEndian, uint16(0x14)) // 연속 읽기
	binary.Write(&body, binary.LittleEndian, uint16(0x00)) // 예약영역
	binary.Write(&body, binary.LittleEndian, uint16(0x01)) // 블록 수
	binary.Write(&body, binary.LittleEndian, uint16(0x06)) // 변수명 길이
	var addr []byte
	addr = fmt.Appendf(addr, "%%%cB%d", register, strAddr*2)
	body.Write(addr)

	binary.Write(&body, binary.LittleEndian, uint16(cnt*2)) // 데이터 갯수
	// --- BODY ---

	// --- HEADER ---
	binary.Write(&buf, binary.LittleEndian, uint16(body.Len())) // Length
	buf.Write([]byte{0x00, 0x00})                               // Frame ID
	// --- HEADER ---

	buf.Write(body.Bytes())

	return buf.Bytes(), nil
}
