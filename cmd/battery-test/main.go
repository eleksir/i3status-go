package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/distatus/battery"
)

// Program entry point.
func main() {
	for {
		batteries, err := battery.GetAll()

		if err != nil {
			fmt.Printf("Unable to get battery info: %w", err)
		}

		spew.Dump(batteries)
		time.Sleep(1 * time.Second)
	}
}
