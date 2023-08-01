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

	now := time.Now()
	device, err := db.GetDevice(context.TODO(), "9e307763-2683-494d-8234-3e01896d8874")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(device)
	}
	fmt.Println(time.Since(now))

	// uid, err := device.ID.Value()
	// fmt.Println(uid.(string))
}
