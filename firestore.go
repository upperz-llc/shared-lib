package sharedlib

import "time"

type FirestoreDevice struct {
	ClientID         string       `firestore:"client_id"`
	DeviceType       int          `firestore:"device_type"`
	LastSeen         time.Time    `firestore:"last_seen"`
	ConnectionStatus string       `firestore:"connection_status"`
	MonitoringStatus string       `firestore:"monitoring_status"`
	FirmwareVersion  string       `firestore:"firmware_version"`
	Nickname         string       `firestore:"nickname"`
	Owner            string       `firestore:"owner"`
	BoundDevices     []string     `firestore:"bound_devices"`
	BoundTo          string       `firestore:"bound_to"`
	Config           DeviceConfig `firestore:"config"`
}

type FirestoreDeviceConfig struct {
	AlertTemperature   int `firestore:"alert_temperature"`
	TargetTemperature  int `firestore:"target_temperature"`
	WarningTemperature int `firestore:"warning_temperature"`
	TelemetryPeriod    int `firestore:"telemetry_period"`
}

func (fa *FirestoreDevice) ToDevice() Device {
	return Device{
		ClientID:         fa.ClientID,
		DeviceType:       fa.DeviceType,
		LastSeen:         fa.LastSeen,
		ConnectionStatus: DeviceConnectionStatus(fa.ConnectionStatus),
		MonitoringStatus: DeviceMonitoringStatus(fa.MonitoringStatus),
		FirmwareVersion:  fa.FirmwareVersion,
		Nickname:         fa.Nickname,
		Owner:            fa.Owner,
		BoundDevices:     fa.BoundDevices,
		BoundTo:          fa.BoundTo,
		Config:           fa.Config,
	}
}

type FirestoreUser struct {
	Email                string `firestore:"email"`
	PhoneNumber          string `firestore:"phone_number"`
	SendSMS              bool   `firestore:"sms_notification"`
	SendPushNotification bool   `firestore:"push_notification"`
}

func (fa *FirestoreUser) ToUser() User {
	return User{
		Email:                fa.Email,
		PhoneNumber:          fa.PhoneNumber,
		SendSMS:              fa.SendSMS,
		SendPushNotification: fa.SendPushNotification,
	}
}
