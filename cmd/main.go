package main

import (
	"iot_connection_tester/internal/usecase"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("사용법: iot_connection_tester <설정 문자열>")
	}
	input := os.Args[1]

	// 테스트 실행
	err := usecase.RunTest(input)
	if err != nil {
		log.Fatalf("테스트 실패: %v", err)
	}
}
