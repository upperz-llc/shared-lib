package main

import (
	"context"
	"fmt"
	"time"

	sharedlib "github.com/upperz-llc/shared-lib"
)

func main() {
	fmt.Println("test")

	db, err := sharedlib.NewCockroachDB(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

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

	err = db.CreateDeviceConfig(context.TODO(), "18b8f73c-b9fd-4b0f-b97e-0c6914efa3e0", sharedlib.DeviceConfig{
		AlertTemperature:   35,
		TargetTemperature:  25,
		WarningTemperature: 30,
		TelemetryPeriod:    60,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	now := time.Now()
	t, err := db.GetDeviceConfig(context.TODO(), "18b8f73c-b9fd-4b0f-b97e-0c6914efa3e0")
	if err != nil {
		fmt.Println(err)
		return
	}
	if t == nil {
		fmt.Println("Not found!")
		return
	}
	fmt.Println(time.Since(now))
	fmt.Println(t)

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
