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

// ***********************************************

// TODO : Add telemetry_timeline
// TODO : Add alarm_timeline
type Device struct {
	ClientID         string                 `json:"client_id" firestore:"client_id"`
	DeviceType       int                    `json:"device_type,omitempty" firestore:"device_type"`
	LastSeen         time.Time              `json:"last_seen" firestore:"last_seen"`
	ConnectionStatus DeviceConnectionStatus `json:"connection_status" firestore:"connection_status"`
	MonitoringStatus DeviceMonitoringStatus `json:"monitoring_status" firestore:"monitoring_status"`
	FirmwareVersion  string                 `json:"firmware_version" firestore:"firmware_version"`
	Nickname         string                 `json:"nickname" firestore:"nickname"`
	Owner            string                 `json:"owner" firestore:"owner"`
	BoundDevices     []string               `json:"bound_devices,omitempty" firestore:"bound_devices,omitempty"`
	BoundTo          string                 `json:"bound_to" firestore:"bound_to"`
	Config           DeviceConfig           `json:"config" firestore:"config"`
}

type DeviceConfig struct {
	AlertTemperature   int `json:"temperature_alert"`
	TargetTemperature  int `json:"temperature_target"`
	WarningTemperature int `json:"temperature_warning"`
	TelemetryPeriod    int `json:"telemetry_period"`
}
