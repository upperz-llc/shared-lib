package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/upperz-llc/shared-lib/db/cockroach"
	"github.com/upperz-llc/shared-lib/device"
)

func main() {
	ctx := context.Background()

	db, err := cockroach.NewCockroachDB(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := db.CreateDevice(ctx, uuid.New().String(), "test", "test", device.Type(1), device.MeasurementType(1)); err != nil {
		fmt.Println(err)
		return
	}

	// device, err := db.GetDevice(ctx, "95b0729d-e863-488b-8f9a-ad2813352588")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("Device: %+v\n", device)

}
