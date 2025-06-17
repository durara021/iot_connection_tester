package main

import (
	"flag"
	"fmt"
	"iot_connection_tester/internal/execute"
	"log"
	"os"
)

func main() {
	help := flag.Bool("help", false, "도움말 출력")
	h := flag.Bool("h", false, "도움말 출력")
	q := flag.Bool("?", false, "도움말 출력")

	flag.Usage = func() {
		fmt.Println("사용법: test <IP 또는 IP:PORT>")
		fmt.Println("예시: test 192.168.0.1:5000")
		fmt.Println("\n옵션:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help || *h || *q {
		flag.Usage()
		return
	}

	if flag.NArg() < 1 {
		fmt.Println("오류: IP 또는 IP:PORT가 필요합니다.")
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Arg(0)

	err := execute.RunTest(input)
	if err != nil {
		log.Fatalf("테스트 실패: %v", err)
	}
}
