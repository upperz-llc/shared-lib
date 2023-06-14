package sharedlib

type User struct {
	UID                  string `json:"uid"` // document id
	Email                string `json:"email,omitempty"`
	PhoneNumber          string `json:"phone_number,omitempty"`
	SendSMS              bool   `json:"sms_notification,omitempty"`
	SendPushNotification bool   `json:"push_notification,omitempty"`
}
