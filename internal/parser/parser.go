package parser

import (
	"encoding/binary"
	"fmt"
)

func ParseWord(data []byte) (uint16, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("적합하지 않은 길이의 데이터입니다!")
	}
	return binary.LittleEndian.Uint16(data), nil
}
