package telemetry

type TelemetryType int

const (
	DeviceTelemetryTemperature TelemetryType = iota
	DeviceTelemetryFirmware
	DeviceTelemetryOTA
)

type Telemetry struct {
	Type TelemetryType `json:"type"`
}

type Temperature struct {
	Telemetry
	Temperature int   `json:"temperature"`
	Timestamp   int64 `json:"timestamp"`
}

type Firmware struct {
	FirmwareVersion string `json:"firmware_version"`
	Timestamp       int64  `json:"timestamp"`
}

type OTA struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}
