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

	fmt.Println(time.Now().Format(time.RFC3339))
	// i := 1

	// for {

	// 	if i == 10 {
	// 		break
	// 	}
	// now := time.Now()
	devices, err := db.GetDevicesByOwner(context.TODO(), "U1uhP6ekkuMYa4VDiTrrGZOtQ3G3")
	if err != nil {
		fmt.Println(err)
		return
	}

	now := time.Now()
	devices, err = db.GetDevicesByOwner(context.TODO(), "U1uhP6ekkuMYa4VDiTrrGZOtQ3G3")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Since(now))

	for _, v := range devices {
		fmt.Println(v)

	}
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
