package telemetry

type TelemetryType int

const (
	Temperature TelemetryType = iota
	Firmware
	OTA
)

type Telemetry struct {
	Type TelemetryType `json:"type"`
}

type TemperatureTelemetry struct {
	Telemetry
	Temperature int   `json:"temperature"`
	Timestamp   int64 `json:"timestamp"`
}
