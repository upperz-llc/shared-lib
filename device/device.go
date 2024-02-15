package device

import (
	"time"
)

// ********** ENUMS AND CONSTANTS ************

type MonitoringStatus int

const (
	Errored MonitoringStatus = iota - 1
	WaitingForConfiguration
	Monitoring
	Alerted
)

func (dms MonitoringStatus) String() string {
	switch dms {
	case Monitoring:
		return "monitoring"
	case WaitingForConfiguration:
		return "waiting_for_configuration"
	case Alerted:
		return "alerted"
	case Errored:
		return "errored"
	}
	return "unknown"
}

type ConnectionStatus int

const (
	Connected ConnectionStatus = iota
	Disconnected
)

func (dms ConnectionStatus) String() string {
	switch dms {
	case Connected:
		return "connected"
	case Disconnected:
		return "disconnected"
	}
	return "unknown"
}

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

type Type int

const (
	WIFI Type = iota
	GATEWAY
	GATEWAYDEVICE
)

type MeasurementType int

const (
	NONE MeasurementType = iota
	INTERNAL
	MAX31855
	MAX31856
)

// ***********************************************

// TODO : Add telemetry_timeline
// TODO : Add alarm_timeline
type Device struct {
	ID               string           `json:"id"`
	Type             Type             `json:"type"`
	ConnectionStatus ConnectionStatus `json:"connection_status"`
	MonitoringStatus MonitoringStatus `json:"monitoring_status"`
	FirmwareVersion  string           `json:"firmware_version"`
	Owner            string           `json:"owner"`
}

type Config struct {
	AlertTemperature   int `json:"temperature_alert"`
	TargetTemperature  int `json:"temperature_target"`
	WarningTemperature int `json:"temperature_warning"`
	TelemetryPeriod    int `json:"telemetry_period"`
	Version            int `json:"version"`
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
