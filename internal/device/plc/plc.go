package plc

import (
	"iot_connection_tester/internal/device"
)

type PLC interface {
	device.Device
	ReadMemory(startAddress string, endAddress string) ([]byte, error)
}
