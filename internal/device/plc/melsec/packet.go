package melsec

func buildMelsecReadPacket(startAddr uint16, wordCount uint16) []byte {
	return []byte{
		// 서브 헤더 및 라우팅 정보
		0x50, 0x00, // Subheader
		0x00, 0x00, // Network / PC No
		0xFF, 0x03, // 요청 대상 모듈
		0x00, 0x00, // 멀티드롭
		0x0C, 0x00, 0x00, 0x00, // 이후 데이터 길이 = 12Byte

		// 커맨드 영역
		0x01, 0x04, // 명령 (0x0401 = Batch Read)
		0x00, 0x00, // Subcommand

		// 디바이스 주소
		byte(startAddr), byte(startAddr >> 8), 0x00,
		0xA8, // D 레지스터 (A8)

		// 요청 갯수
		byte(wordCount), byte(wordCount >> 8),
	}
}
