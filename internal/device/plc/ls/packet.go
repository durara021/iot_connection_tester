package ls

import (
	"bytes"
	"encoding/binary"
)

// LS 레지스터 읽기 패킷 생성
// @param strAddr string: 시작 주소 (예: "%D100")
// @param cnt uint16: 읽을 레지스터 수
// @return []byte: 생성된 패킷 바이트 슬라이스
func buildLSReadPacket(strAddr string, cnt uint16) []byte {
	var buf bytes.Buffer

	buf.Write([]byte{0x4C, 0x53, 0x49, 0x53}) // "LSIS" 헤더

	header := []byte{
		0x00, 0x00, // Reserved
		0x00, 0x00, // Header length
		0x00, 0x00, // Data length
		0x00, // CPU info
		0x00, // Reserved
		0x00, // Source of request
		0x54, // 명령어 종류: 0x54 = Read
	}
	buf.Write(header)

	buf.WriteString(strAddr)                     // 시작 주소
	binary.Write(&buf, binary.LittleEndian, cnt) // 레지스터 수

	return buf.Bytes()
}
