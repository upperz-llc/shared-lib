package sharedlib

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func (fa *FirestoreDevice) toDevice() Device {
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

func (fa *FirestoreUser) toUser() User {
	return User{
		Email:                fa.Email,
		PhoneNumber:          fa.PhoneNumber,
		SendSMS:              fa.SendSMS,
		SendPushNotification: fa.SendPushNotification,
	}
}

// ***********************************************************

type FirebaseDB struct {
	DB *firestore.Client
}

func (fdb *FirebaseDB) GetDeviceOwner(ctx context.Context, deviceID string) (string, error) {
	docsnapshot, err := fdb.DB.Collection("devices").Doc(deviceID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", nil
		}
		return "", err
	}

	datamap := docsnapshot.Data()

	val, ok := datamap["owner"]
	if !ok {
		return "", nil
	}

	return val.(string), nil
}

func (fdb *FirebaseDB) GetUser(ctx context.Context, uid string) (*User, error) {
	firestoreUser, err := fdb.getUser(ctx, uid)
	if err != nil || firestoreUser == nil {
		return nil, err
	}

	user := firestoreUser.toUser()
	return &user, nil
}

func (fdb *FirebaseDB) GetUserNotificationSettings(ctx context.Context, uid string) (*User, error) {
	firestoreUser, err := fdb.getUser(ctx, uid)
	if err != nil || firestoreUser == nil {
		return nil, err
	}

	user := firestoreUser.toUser()
	return &user, nil
}

func (fdb *FirebaseDB) getUser(ctx context.Context, uid string) (*FirestoreUser, error) {
	docsnapshot, err := fdb.DB.Collection("users").Doc(uid).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	var user FirestoreUser
	if err := docsnapshot.DataTo(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func NewDBFirestore(ctx context.Context) (*FirebaseDB, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}

	db, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}
	return &FirebaseDB{
		DB: db,
	}, nil
}
