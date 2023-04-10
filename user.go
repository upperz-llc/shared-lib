package sharedlib

type User struct {
	UID                  string `json:"uid"` // document id
	Email                string `json:"email"`
	PhoneNumber          string `json:"phone_number"`
	SendSMS              bool   `json:"sms_notification"`
	SendPushNotification bool   `json:"push_notification"`
}
