package sharedlib

import "time"

type DeviceType int

const (
	WIFI DeviceType = iota
	GATEWAY
	GATEWAYDEVICE
)

type MeasurementType int

const (
	INTERNAL MeasurementType = iota
	MAX31855
	MAX31856
)

// ManufacturingCreateDeviceRequest placeholder
type ManufacturingData struct {
	DeviceID        string          `json:"device_id"`
	Username        string          `json:"username"`
	ManufacturedAt  time.Time       `json:"manufactured_at"`
	DeviceType      DeviceType      `json:"device_type"`
	MeasurementType MeasurementType `json:"measurement_type"`
}

type NewAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
