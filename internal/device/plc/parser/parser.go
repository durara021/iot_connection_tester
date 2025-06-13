package parser

import (
	"encoding/binary"
	"iot_connection_tester/internal/common/errs"
	"iot_connection_tester/internal/setting"
)

// PLC 응답 데이터 파싱 함수
// @param data: PLC에서 수신한 원시 바이트 배열 (2바이트 단위)
// @param settings: 읽은 주소에 대응되는 태그 설정 정보 배열
// @return map[string]uint16: 태그 이름과 값의 매핑 결과
// @return error: 파싱 실패 또는 설정 없음 등의 에러
func ParseData(data []byte, settings []setting.Setting) (map[string]uint16, error) {
	// 유효성 검사: 바이트 수가 2의 배수가 아니면 오류
	if len(data)%2 != 0 {
		return nil, errs.NewErrs("", "", errs.ErrCodeDataParseFailed, nil)
	}
	// 설정 정보가 없을 경우 오류
	if len(settings) == 0 {
		return nil, errs.NewErrs("", "", errs.ErrCodeEmptyResult, nil)
	}

	result := make(map[string]uint16)
	baseAddr := settings[0].Address // 기준 주소 설정

	for _, s := range settings {
		offset := int(s.Address - baseAddr) // 기준 주소로부터의 오프셋
		start := offset * 2                 // 바이트 오프셋 계산
		if start+1 >= len(data) {
			continue // 범위를 벗어나면 무시
		}
		val := binary.BigEndian.Uint16(data[start : start+2]) // 빅엔디안으로 2바이트 정수 해석
		result[s.Value] = val                                 // 태그 이름 (s.Value)에 값 할당
	}

	// 결과가 비어 있으면 오류 반환
	if len(result) == 0 {
		return nil, errs.NewErrs("", "", errs.ErrCodeEmptyResult, nil)
	}

	return result, nil
}
