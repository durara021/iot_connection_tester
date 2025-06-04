package cnc

import (
	"iot_connection_tester/internal/device"
)

type CNC interface {
	device.Device
	GetOperatingState() (string, error) // 예: "Idle", "Running", "Alarm"
	GetSpindleSpeed() (int, error)
}
