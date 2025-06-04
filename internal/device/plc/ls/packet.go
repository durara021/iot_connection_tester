package ls

import (
	"bytes"
	"encoding/binary"
)

func buildLSREadPacket(address string, wordCount uint16) []byte {
	var buf bytes.Buffer

	buf.Write([]byte{0x4C, 0x53, 0x49, 0x53})
	// 명령 종류: 0x54 = Read
	header := []byte{
		0x00, 0x00, // Reserved
		0x00, 0x00, // Header length
		0x00, 0x00, // Data length
		0x00, // CPU info
		0x00, // Reserved
		0x00, // Source of request
		0x54, // Read command
	}
	buf.Write(header)

	// 주소 포맷 → ASCII 문자열 ("%D100")
	buf.WriteString(address)

	// 개수
	binary.Write(&buf, binary.LittleEndian, wordCount)

	return buf.Bytes()
}
