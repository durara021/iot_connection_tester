package errs

import (
	"fmt"
)

// 장비/프로토콜 에러 커스텀
type Errs struct {
	DeviceType    string // 장비 종류 (예: FANUC_CNC, MELSEC_PLC, LS_PLC)
	ProtocolType  string // 프로토콜 종류 (예: FOCAS, SLMP, XGT)
	ErrorCode     int    // 정의된 에러 코드
	Message       string // 에러에 대한 설명 메시지
	OriginalError error  // 원래 발생한 Go 에러
}

// @param deviceType: 장비 종류 문자열
// @param protocolType: 프로토콜 종류 문자열
// @param errorCode: 정의된 에러 코드
// @param originalErr: 래핑할 원본 에러
// @return *Errs: 에러 정보가 포함된 포인터 객체
func NewErrs(deviceType, protocolType string, errorCode int, originalErr error) *Errs {
	message := GetErrorMessage(errorCode)
	return &Errs{
		DeviceType:    deviceType,
		ProtocolType:  protocolType,
		ErrorCode:     errorCode,
		Message:       message,
		OriginalError: originalErr,
	}
}

// @param errorCode: 정의된 에러 코드
// @return string: 설명 메시지 (정의되지 않은 경우는 기본 메시지 반환)
func GetErrorMessage(errorCode int) string {
	if msg, ok := errorCodeMessages[errorCode]; ok {
		return msg
	}
	return fmt.Sprintf("알 수 없는 에러 코드 (%d)", errorCode)
}

// error 인터페이스 구현 메서드, 에러 메시지 반환
func (e *Errs) Error() string {
	var prefix string
	switch {
	case e.DeviceType != "" && e.ProtocolType != "":
		prefix = fmt.Sprintf("[%s/%s] ", e.DeviceType, e.ProtocolType)
	case e.DeviceType != "":
		prefix = fmt.Sprintf("[%s] ", e.DeviceType)
	case e.ProtocolType != "":
		prefix = fmt.Sprintf("[%s] ", e.ProtocolType)
	default:
		prefix = ""
	}

	if e.OriginalError != nil {
		return fmt.Sprintf("%s \n ErrorCode: %d \n Message: %s \n Original: %v",
			prefix, e.ErrorCode, e.Message, e.OriginalError)
	}
	return fmt.Sprintf("%s \n ErrorCode: %d \n Message: %s ",
		prefix, e.ErrorCode, e.Message)
}

// errors 패키지의 Unwrap 기능을 지원, 원본 에러 추출
func (e *Errs) Unwrap() error {
	return e.OriginalError
}

// ------------------------- 상수 정의 -----------------------------

// Device Types: 장비 종류 식별용 문자열 상수
const (
	DeviceTypeFanuc  = "FANUC_CNC"
	DeviceTypeMelsec = "MELSEC_PLC"
	DeviceTypeLS     = "LS_PLC"
)

// Protocol Types: 통신 프로토콜 식별용 문자열 상수
const (
	ProtocolTypeFocas = "FOCAS"
	ProtocolTypeSLMP  = "SLMP"
	ProtocolTypeXGT   = "XGT"
)

// 통신 관련 에러 코드
const (
	ErrCodeConnectionFailed = 100 // 장비 연결 실패
	ErrCodeTimeout          = 101 // 통신 응답 시간 초과
	ErrCodeInvalidResponse  = 102 // 잘못된 응답 포맷
	ErrCodeDeviceError      = 103 // 장비 내부 에러
	ErrCodeReadFailed       = 104 // 데이터 읽기 실패
	ErrCodeWriteFailed      = 105 // 데이터 쓰기 실패
	ErrCodeCloseFailed      = 106 // 연결 해제 실패
)

// 애플리케이션 내부 에러 코드
const (
	ErrCodeConfigParseFailed  = 200 // 설정 파일 파싱 실패
	ErrCodeDeviceCreateFailed = 201 // 장비 객체 생성 실패
	ErrCodeEmptyResult        = 202 // 반환값이 비어 있음
	ErrCodeDataParseFailed    = 300 // 데이터 파싱 실패
)

// errorCodeMessages는 에러 코드에 대응하는 설명 메시지를 저장한 맵입니다.
var errorCodeMessages = map[int]string{
	ErrCodeConnectionFailed: "장비 연결에 실패했습니다.",
	ErrCodeTimeout:          "통신 응답 시간이 초과되었습니다.",
	ErrCodeInvalidResponse:  "장비 응답 형식이 유효하지 않습니다.",
	ErrCodeDeviceError:      "장비 내부에서 오류가 발생했습니다.",
	ErrCodeReadFailed:       "데이터 읽기에 실패했습니다.",
	ErrCodeWriteFailed:      "데이터 쓰기에 실패했습니다.",
	ErrCodeCloseFailed:      "연결 해제에 실패했습니다.",

	ErrCodeConfigParseFailed:  "설정 파일 파싱 중 오류가 발생했습니다.",
	ErrCodeDeviceCreateFailed: "장비 객체 생성에 실패했습니다.",
	ErrCodeEmptyResult:        "요청한 데이터가 비어있습니다.",

	ErrCodeDataParseFailed: "요청한 데이터 파싱 중 오류가 발생했습니다.",
}
