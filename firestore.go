package sharedlib

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
