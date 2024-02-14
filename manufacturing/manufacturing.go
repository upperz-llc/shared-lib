package manufacturing

import (
	"time"

	"github.com/upperz-llc/shared-lib/device"
)

// ManufacturingCreateDeviceRequest placeholder
type ManufacturingData struct {
	DeviceID        string                 `json:"device_id"`
	Username        string                 `json:"username"`
	Password        string                 `json:"password"`
	ManufacturedAt  time.Time              `json:"manufactured_at"`
	DeviceType      device.Type            `json:"device_type"`
	MeasurementType device.MeasurementType `json:"measurement_type"`
}

type NewAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
