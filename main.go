package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	// "strings"

	. "./stuff"
)

func main() {
	// Fetch ALL THE THINGS
	//allDevices := UseCache("cache/cmwiki", CMWikiScraper)()
	//allDevices := UseCache("cache/phonescoop", PhoneScoopScraper)()
	allDevices := UseCache("cache/lineagewiki", LineageWikiScraper)()

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
		{"require recent release year", func(device Vitalstats) bool {
			return strings.Contains(device.ReleaseDate, "2017") ||
				strings.Contains(device.ReleaseDate, "2016")
		}},
		{"discard known non-removable batt", func(device Vitalstats) bool {
			return device.BatteryRem != "no"
		}},
		{"require definitely removable batt", func(device Vitalstats) bool {
			return device.BatteryRem == "yes"
		}},
		{"require latest lineageOS", func(device Vitalstats) bool {
			// 13~15 are what's on https://www.lineageoslog.com/ today.
			return strings.Contains(device.CMSupport, "13") ||
				strings.Contains(device.CMSupport, "14") ||
				strings.Contains(device.CMSupport, "15")
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
