package sharedlib

type User struct {
	Email                string `json:"email"`
	PhoneNumber          string `json:"phone_number"`
	SendSMS              bool   `json:"sms_notification"`
	SendPushNotification bool   `json:"push_notification"`
}
