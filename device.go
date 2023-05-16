package sharedlib

import "time"

// ********** ENUMS AND CONSTANTS ************

type DeviceMonitoringStatus string

const (
	Monitoring DeviceMonitoringStatus = "monitoring"
	Alerted    DeviceMonitoringStatus = "alerted"
	Errored    DeviceMonitoringStatus = "errored"
)

func (dms DeviceMonitoringStatus) String() string {
	switch dms {
	case Monitoring:
		return "monitoring"
	case Alerted:
		return "alerted"
	case Errored:
		return "errored"
	}
	return "unknown"
}

type DeviceConnectionStatus string

const (
	Connected    DeviceConnectionStatus = "connected"
	Disconnected DeviceConnectionStatus = "disconnected"
)

type DeviceCommands int

const (
	PING DeviceCommands = iota
	OTA
)

type BrokerCommands int

const (
	BIND BrokerCommands = iota
	UNBIND
)

type OTAStatus int

const (
	STARTED OTAStatus = iota
	NO_NEW_VERSION
	UPDATING
	FINISHED
	FAILED
)

func (otas OTAStatus) String() string {
	switch otas {
	case STARTED:
		return "started"
	case NO_NEW_VERSION:
		return "no new version"
	case UPDATING:
		return "updating"
	case FINISHED:
		return "finished"
	case FAILED:
		return "failed"
	}
	return "unknown"
}

// ***********************************************

// TODO : Add telemetry_timeline
// TODO : Add alarm_timeline
type Device struct {
	ClientID         string                 `json:"client_id"`
	DeviceType       int                    `json:"device_type"`
	LastSeen         time.Time              `json:"last_seen"`
	ConnectionStatus DeviceConnectionStatus `json:"connection_status"`
	MonitoringStatus DeviceMonitoringStatus `json:"monitoring_status"`
	FirmwareVersion  string                 `json:"firmware_version"`
	Nickname         string                 `json:"nickname"`
	Owner            string                 `json:"owner"`
	BoundDevices     []string               `json:"bound_devices,omitempty"`
	BoundTo          string                 `json:"bound_to"`
	Temperature      float64                `json:"temperature"`
	Config           DeviceConfig           `json:"config"`
}

type DeviceConfig struct {
	AlertTemperature   int `json:"temperature_alert"`
	TargetTemperature  int `json:"temperature_target"`
	WarningTemperature int `json:"temperature_warning"`
	TelemetryPeriod    int `json:"telemetry_period"`
}

type DeviceTelemetry struct {
	CreatedAt   time.Time `json:"created_at"`
	DeleteAt    time.Time `json:"deleted_at"`
	Timestamp   int64     `json:"timestamp"`
	Temperature float64   `json:"temperature"`
}

type OTATelemetry struct {
	Timestamp time.Time `json:"timestamp"`
	Status    OTAStatus `json:"status"`
}
