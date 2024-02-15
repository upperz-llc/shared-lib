package messages

type MessageType int

const (
	Telemetry MessageType = iota
	Firmware
	State
	HealthCheck
)
