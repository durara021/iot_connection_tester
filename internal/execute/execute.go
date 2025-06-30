package execute

import (
	"encoding/json"
	"fmt"
	"iot_connection_tester/internal/common/errs"
	"iot_connection_tester/internal/device"
	"iot_connection_tester/internal/setting"
	"reflect"
	"strings"
	"time"
)

// 문자열 기반 장비 테스트수행
// @param input: 장치 설정 문자열 (예: "melsec,192.168.0.1:5000")
// @return error: 테스트 실패 시 상세 에러, 성공 시 nil
func RunTest(input string) error {
	start := time.Now()

	// 설정 파싱
	cfg, err := setting.ParseDeviceConfig(input)
	if err != nil {
		return errs.NewErrs("", "", errs.ErrCodeConfigParseFailed, err)
	}

	// 장비 생성
	dev, err := device.NewDevice(cfg)
	if err != nil {
		return errs.NewErrs(cfg.Device, "", errs.ErrCodeConfigParseFailed, err)
	}

	// dev가 nil인지 확인 필요!
	if dev == nil || reflect.ValueOf(dev).IsNil() {
		return fmt.Errorf("장비 생성 실패: dev is nil")
	}

	defer func() {
		if dev != nil {
			if err := dev.Close(); err != nil {
				fmt.Printf("⚠️ 장비 닫기 실패: %v\n", err)
			}
		}
	}()

	// 최대 3회 테스트 시도
	var finalErr error
	for attempt := 1; attempt <= 3; attempt++ {
		if attempt != 1 {
			fmt.Printf("▶️ 테스트 시도 %d회...\n", attempt)
		}

		finalErr = runDeviceTest(dev)

		if finalErr == nil {
			duration := time.Since(start)
			fmt.Printf("✅ 테스트 성공 (%d회 시도): %s\n", attempt, duration)
			return nil
		}

		fmt.Printf("⚠️ 테스트 실패 (%d회 시도): %v\n", attempt, finalErr)
		time.Sleep(10 * time.Millisecond)
	}

	duration := time.Since(start)
	fmt.Printf("❌ 테스트 최종 실패: %s\n", duration)
	return finalErr
}

// 데이터 수집 및 유효성 검사
// @param dev: device.Device 인터페이스 (Connect, Test, Close 함수 제공)
// @return error: 실패 시 에러, 성공 시 nil
func runDeviceTest(dev device.Device) error {
	if err := dev.Connect(); err != nil {
		return errs.NewErrs("", "", errs.ErrCodeConnectionFailed, err)
	}

	result, err := dev.Test()
	if err != nil {
		return errs.NewErrs("", "", errs.ErrCodeReadFailed, err)
	}

	if len(result) == 0 {
		return errs.NewErrs("", "", errs.ErrCodeEmptyResult, nil)
	}

	printJsonResult(result)
	// printResult(result)
	return nil
}

// 테스트 결과를 표 형태로 출력
// @param result: key=태그명(string), value=수집된 데이터(uint16)
func printJsonResult(result map[string]uint16) {
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("❌ JSON 변환 실패:", err)
		return
	}
	fmt.Println(string(jsonBytes))
}

// 테스트 결과를 표 형태로 출력
// @param result: key=태그명(string), value=수집된 데이터(uint16)
func printResult(result map[string]uint16) {
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}

	colWidth := 22

	for _, key := range keys {
		fmt.Printf("%s | ", center(key, colWidth))
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", (colWidth+3)*len(result)-1))

	for _, key := range keys {
		fmt.Printf("%s | ", center(fmt.Sprintf("%d", result[key]), colWidth))
	}
	fmt.Println()
}

// 텍스트를 지정된 너비로 가운데 정렬한 문자열 반환
// @param text: 출력할 문자열
// @param width: 전체 칼럼 너비
// @return string: 가운데 정렬된 결과 문자열
func center(text string, width int) string {
	pad := width - len(text)
	if pad <= 0 {
		return text
	}
	left := pad / 2
	right := pad - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}
