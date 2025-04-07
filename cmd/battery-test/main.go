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
		batteries, _ := battery.GetAll()

		for index, battery := range batteries {
			fmt.Printf("Battery %s: %s\n", index, spew.Sdump(battery))
		}

		time.Sleep(1 * time.Second)
	}
}
