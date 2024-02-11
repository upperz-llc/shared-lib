package sharedlib

import (
	"time"
)

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
	ID               string                 `json:"id"`
	DeviceType       int64                  `json:"device_type"`
	LastSeen         time.Time              `json:"last_seen"`
	ConnectionStatus DeviceConnectionStatus `json:"connection_status"`
	MonitoringStatus DeviceMonitoringStatus `json:"monitoring_status"`
	FirmwareVersion  string                 `json:"firmware_version"`
	Nickname         string                 `json:"nickname"`
	Owner            string                 `json:"owner"`
	// BoundDevices     []string               `json:"bound_devices,omitempty"`
	// BoundTo          string                 `json:"bound_to"`
	Temperature int64 `json:"temperature"`
}

type DeviceConfig struct {
	AlertTemperature   int `json:"temperature_alert"`
	TargetTemperature  int `json:"temperature_target"`
	WarningTemperature int `json:"temperature_warning"`
	TelemetryPeriod    int `json:"telemetry_period"`
	Version            int `json:"version"`
}

type DeviceTelemetry struct {
	ID          string    `json:"id"`
	DeviceID    string    `json:"device_id"`
	Temperature int64     `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

type OTATelemetry struct {
	Timestamp time.Time `json:"timestamp"`
	Status    OTAStatus `json:"status"`
}

type TelemetryRange int

const (
	Hour TelemetryRange = iota
	SixHour
	Day
)

func (tr TelemetryRange) ToTime() time.Time {
	switch tr {
	case Hour:
		return time.Now().Add(-1 * time.Hour)
	case SixHour:
		return time.Now().Add(-6 * time.Hour)
	case Day:
		return time.Now().Add(-24 * time.Hour)

	}
	return time.Now()
}
