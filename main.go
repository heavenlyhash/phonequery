package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	. "./stuff"
)

func main() {
	// Fetch ALL THE THINGS
	allDevices := UseCache("cache", CMWikiScraper)()

	fmt.Println()
	fmt.Printf("%d devices in total.\n", len(allDevices))
	fmt.Println("====================================")
	fmt.Println("valid options....")
	fmt.Println()

	// Declare filters.
	// **EDIT HERE** if you have different opinions.
	// You can see what's in the vital stats struct in "./stuff/vitals.go".
	filters := []filter{
		{"discard phablet", func(device Vitalstats) bool {
			return device.Type != "phablet"
		}},
		{"discard tablet", func(device Vitalstats) bool {
			return device.Type != "tablet"
		}},
		{"discard known non-removable batt", func(device Vitalstats) bool {
			return !strings.Contains(device.Power, "on-removable") && // skip leading 'n' because case
				!strings.Contains(device.Power, "un-removable") // lol
		}},
		{"require definitely removable batt", func(device Vitalstats) bool {
			return strings.Contains(device.Power, "removable") ||
				strings.Contains(device.Power, "removeable") // lol
		}},
		{"require latest CM", func(device Vitalstats) bool {
			return strings.Contains(device.CMSupport, "12")
		}},
	}

	// Apply the filters.
	for _, filter := range filters {
		fmt.Printf("filtering %q... ", filter.name)
		survivingDevices := make([]Vitalstats, 0)
		for _, dev := range allDevices {
			if filter.fn(dev) {
				survivingDevices = append(survivingDevices, dev)
			}
		}
		fmt.Printf("%d valid options remaining\n", len(survivingDevices))
		allDevices = survivingDevices
	}

	// Spit em all out, whatever's left.
	for _, device := range allDevices {
		json.NewEncoder(os.Stdout).Encode(device)
	}
}

type filter struct {
	name string
	fn   func(Vitalstats) bool
}
