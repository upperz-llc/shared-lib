package device

type TelemetryType int

const (
	TelemetryTemperature TelemetryType = iota
	TelemetryFirmware
	TelemetryState
	TelemetryHealthCheck
	TelemetryOTA
)

type Telemetry struct {
	Type      TelemetryType `json:"type"`
	Value     int           `json:"value"`
	Timestamp int64         `json:"timestamp"`
}

type Firmware struct {
	FirmwareVersion string `json:"firmware_version"`
	Timestamp       int64  `json:"timestamp"`
}

// type Telemetry struct {
// 	ID          string    `json:"id"`
// 	DeviceID    string    `json:"device_id"`
// 	Temperature int64     `json:"temperature"`
// 	Timestamp   time.Time `json:"timestamp"`
// }

// type OTA struct {
// 	Status    string `json:"status"`
// 	Timestamp int64  `json:"timestamp"`
// }
