package main

import (
	"context"
	"fmt"

	"github.com/upperz-llc/shared-lib/db/cockroach"
)

func main() {
	ctx := context.Background()

	db, err := cockroach.NewCockroachDB(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	device, err := db.GetDevice(ctx, "95b0729d-e863-488b-8f9a-ad2813352588")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Device: %+v\n", device)

}
