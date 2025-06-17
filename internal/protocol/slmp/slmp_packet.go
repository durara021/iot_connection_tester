package protocol

// SLMP 레지스터 읽기 패킷 생성
// @param strAddr string: 시작 주소 (예: "%D100")
// @param cnt uint16: 읽을 레지스터 수
// @return []byte: 생성된 패킷 바이트 슬라이스
func buildMelsecReadPacket(register byte, strAddr uint16, cnt uint16) []byte {
	return []byte{
		// 서브 헤더 및 라우팅 정보
		0x50, 0x00, // Subheader (Fixed value: 0x0050)
		0x00, 0x00, // Network No. / PC No.
		0xFF, 0x03, // 요청 대상 모듈 I/O 번호 / 스테이션 번호 (FF03 = 기본 CPU)
		0x00, 0x00, // 멀티드롭 설정 (Not used)

		// 데이터 길이 (이후 바이트 수, 12바이트 고정)
		0x0C, 0x00, 0x00, 0x00,

		// 커맨드 영역
		0x01, 0x04, // Command: 0x0401 (Batch Read)
		0x00, 0x00, // Subcommand: 0x0000 (디바이스 지정 단위 기본 설정)

		// 디바이스 주소
		byte(strAddr),      // 시작 주소 (하위 바이트)
		byte(strAddr >> 8), // 시작 주소 (상위 바이트)
		0x00,               // 확장 주소 (보통 사용 안 함)
		register,           // 디바이스 종류 (예: 'D' = 0xA8)

		// 요청 개수
		byte(cnt),      // 읽을 개수 (하위 바이트)
		byte(cnt >> 8), // 읽을 개수 (상위 바이트)
	}
}
