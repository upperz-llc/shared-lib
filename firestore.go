package sharedlib

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ****************** ALARM ***********************

type FirestoreAlarm struct {
	ID              string    `firestore:"id"`
	Type            int       `firestore:"type"`
	ClientID        string    `firestore:"client_id"`
	Acked           bool      `firestore:"acked"`
	Active          bool      `firestore:"active"`
	CreatedAt       time.Time `firestore:"created_at,omitempty"`
	ClosedAt        time.Time `firestore:"closed_at,omitempty"`
	AckedAt         time.Time `firestore:"acked_at,omitempty"`
	AckedBy         string    `firestore:"acked_by,omitempty"`
	AckedCheckCount int       `firestore:"acked_check_count,omitempty"`
}

func (fa *FirestoreAlarm) toAlarm() Alarm {
	return Alarm{
		ID:              fa.ID,
		Type:            fa.Type,
		ClientID:        fa.ClientID,
		CreatedAt:       fa.CreatedAt,
		ClosedAt:        fa.ClosedAt,
		AckedAt:         fa.AckedAt,
		AckedBy:         fa.AckedBy,
		Active:          fa.Active,
		AckedCheckCount: fa.AckedCheckCount,
		Acked:           fa.Acked,
	}
}

type FirestoreAlarmTimeline struct {
	ID        string    `firestore:"id"`
	Type      int       `firestore:"type"`
	CreatedAt time.Time `firestore:"created_at"`
	ClosedAt  time.Time `firestore:"closed_at,omitempty"`
}

func (fdb *FirebaseDB) CreateAlarmConnection(ctx context.Context, clientID string) (*Alarm, error) {
	alarmID := generateRandomString(20)

	alarm := FirestoreAlarm{
		ID:        alarmID,
		Type:      Connection,
		ClientID:  clientID,
		CreatedAt: time.Now(),
		Active:    true,
	}

	storedAlarm, err := fdb.QueryAlarm(ctx, clientID, Connection)
	if err != nil {
		return nil, err
	}

	// alarm already active
	if storedAlarm != nil {
		return nil, errors.New("alarm already active")
	}

	_, err = fdb.DB.Collection("alarms").Doc(alarmID).Set(ctx, alarm)
	if err != nil {
		return nil, err
	}

	returnAlarm := alarm.toAlarm()
	return &returnAlarm, nil
}

func (fdb *FirebaseDB) AddNewAlarmToAlarmTimeline(ctx context.Context, alarm Alarm) error {
	docref := fdb.DB.Collection("devices").Doc(alarm.ClientID).Collection("alarm_timeline").Doc(alarm.ID)

	_, err := docref.Create(ctx, FirestoreAlarmTimeline{
		ID:        alarm.ID,
		Type:      alarm.Type,
		CreatedAt: alarm.CreatedAt,
	})

	return err
}

func (fdb *FirebaseDB) DeleteAlarm(ctx context.Context, alarmID string) error {
	_, err := fdb.DB.Collection("alarms").Doc(alarmID).Delete(ctx)

	return err
}

func (fdb *FirebaseDB) CloseAlarm(ctx context.Context, alarm *Alarm) error {
	alarm.ClosedAt = time.Now()
	alarm.Active = false
	_, err := fdb.DB.Collection("alarms").Doc(alarm.ID).Update(ctx, []firestore.Update{
		{
			Path:  "closed_at",
			Value: alarm.ClosedAt,
		},
		{
			Path:  "active",
			Value: alarm.Active,
		},
	})
	return err
}

// GetAlarm placeholder
func (fdb *FirebaseDB) GetAlarm(ctx context.Context, alarmID string) (*Alarm, error) {
	firestoreAlarm, err := fdb.getAlarm(ctx, alarmID)
	if err != nil || firestoreAlarm == nil {
		return nil, err
	}

	alarm := firestoreAlarm.toAlarm()
	return &alarm, nil
}

func (fdb *FirebaseDB) IncrementAlarmAckCheckCount(ctx context.Context, alarmID string) error {
	_, err := fdb.DB.Collection("alarms").Doc(alarmID).Update(ctx, []firestore.Update{
		{
			Path:  "acked_check_count",
			Value: firestore.Increment(1),
		},
	})

	return err
}

func (fdb *FirebaseDB) QueryAlarm(ctx context.Context, clientID string, alarmType AlarmType) (*Alarm, error) {
	alarms := make([]FirestoreAlarm, 0)
	iter := fdb.DB.Collection("alarms").Where("client_id", "==", clientID).Where("type", "==", alarmType).Where("active", "==", true).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var alarm FirestoreAlarm
		err = doc.DataTo(&alarm)
		if err != nil {
			return nil, err
		}

		alarms = append(alarms, alarm)
	}

	// no alarms returned
	if len(alarms) == 0 {
		return nil, nil
	}
	// too many alarms returned
	if len(alarms) > 1 {
		return nil, errors.New("too many alarms returned")
	}

	returnAlarm := alarms[0].toAlarm()
	return &returnAlarm, nil
}

func (fdb *FirebaseDB) UpdateAlarmTimelineWithClosedAt(ctx context.Context, alarm Alarm) error {
	docref := fdb.DB.Collection("devices").Doc(alarm.ClientID).Collection("alarm_timeline").Doc(alarm.ID)

	_, err := docref.Update(ctx, []firestore.Update{
		{
			Path:  "closed_at",
			Value: alarm.ClosedAt,
		},
	})

	return err
}

func (fdb *FirebaseDB) getAlarm(ctx context.Context, alarmID string) (*FirestoreAlarm, error) {
	snapshot, err := fdb.DB.Collection("alarms").Doc(alarmID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	var alarm FirestoreAlarm
	if err := snapshot.DataTo(&alarm); err != nil {
		return nil, err
	}

	return &alarm, nil
}

// ****************** Device *******************

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

// ******************* Users *********************

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

// ****************** Device ******************

func (fdb *FirebaseDB) CreateDevice(ctx context.Context, device Device) error {
	firestoredevice := FirestoreDevice{
		ClientID:         device.ClientID,
		DeviceType:       device.DeviceType,
		LastSeen:         device.LastSeen,
		ConnectionStatus: string(device.ConnectionStatus),
		MonitoringStatus: string(device.MonitoringStatus),
		FirmwareVersion:  device.FirmwareVersion,
		Nickname:         device.Nickname,
		Owner:            device.Owner,
		BoundDevices:     device.BoundDevices,
		BoundTo:          device.BoundTo,
		Config:           device.Config,
	}

	_, err := fdb.DB.Collection("devices").Doc(firestoredevice.ClientID).Set(ctx, firestoredevice)
	return err
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

// ****************************************

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
