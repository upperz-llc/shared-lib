package main

import (
	"context"
	"fmt"

	sharedlib "github.com/upperz-llc/shared-lib"
)

func main() {
	fmt.Println("test")

	db, err := sharedlib.NewCockroachDB(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	alarms, err := db.QueryAlarmsByUser(context.TODO(), "pLV0ujmdmWh81YdLIWeYrk2q5Qk2")
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range alarms {
		fmt.Printf("%+v\n", v)
	}

	// if devices, err := db.GetInactiveGatewayDevices(context.TODO(), time.Now().Add(-5*time.Minute)); err != nil {
	// 	fmt.Println(err)
	// 	return
	// } else {
	// 	fmt.Println(len(devices))
	// }

	// a, _ := db.CreateAlarm(context.TODO(), "18b8f73c-b9fd-4b0f-b97e-0c6914efa3e0", alarm.Connection)

	// ra, err := db.CloseAlarm(context.TODO(), a.ID)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(ra.ClosedAt)

	// now := time.Now()
	// t, err := db.GetDeviceTelemetry(context.TODO(), "18b8f73c-b9fd-4b0f-b97e-0c6914efa3e0", sharedlib.Hour)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(time.Since(now))
	// fmt.Println(len(t))

	// now = time.Now()
	// t, err = db.GetDeviceTelemetry(context.TODO(), "18b8f73c-b9fd-4b0f-b97e-0c6914efa3e0", sharedlib.SixHour)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(time.Since(now))
	// fmt.Println(len(t))

	// now = time.Now()
	// t, err = db.GetDeviceTelemetry(context.TODO(), "18b8f73c-b9fd-4b0f-b97e-0c6914efa3e0", sharedlib.Day)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(time.Since(now))
	// fmt.Println(len(t))

	// err = db.CreateDeviceConfig(context.TODO(), "18b8f73c-b9fd-4b0f-b97e-0c6914efa3e0", sharedlib.DeviceConfig{
	// 	AlertTemperature:   35,
	// 	TargetTemperature:  25,
	// 	WarningTemperature: 30,
	// 	TelemetryPeriod:    60,
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// now := time.Now()
	// t, err := db.GetDeviceConfig(context.TODO(), "18b8f73c-b9fd-4b0f-b97e-0c6914efa3e0")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// if t == nil {
	// 	fmt.Println("Not found!")
	// 	return
	// }
	// fmt.Println(time.Since(now))
	// fmt.Println(t)

	// device, err := db.GetACL(context.TODO(), "9e307763-2683-494d-8234-3e01896d8874", "CONFIG/9e307763-2683-494d-8234-3e01896d8874s")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if device == nil {
	// 	fmt.Println("not found")
	// 	return
	// }
	// fmt.Println(time.Since(now))
	// fmt.Println(device.Allowed)
	// // 	i++
	// }

	// uid, err := device.ID.Value()
	// fmt.Println(uid.(string))
}
