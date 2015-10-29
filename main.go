package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Vitalstats struct {
	Name        string
	Link        string
	Power       string
	ReleaseDate string
	Type        string // so far, ["phone"|"tablet"|"phablet"]
	CMSupport   string
}

func main() {
	var allDevices []Vitalstats
	cacheFile, err := os.Open("cache")
	if err != nil {
		cacheFile, err = os.OpenFile("cache", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		allDevices = scrape()
		enc := json.NewEncoder(cacheFile)
		for _, dev := range allDevices {
			enc.Encode(dev)
		}
	} else {
		dec := json.NewDecoder(cacheFile)
		for {
			var vitals Vitalstats
			err := dec.Decode(&vitals)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			allDevices = append(allDevices, vitals)
		}
	}

	fmt.Println()
	fmt.Printf("%d devices in total.\n", len(allDevices))
	fmt.Println("====================================")
	fmt.Println("valid options....")
	fmt.Println()

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

	//var survivingDevices []Vitalstats
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

	for _, device := range allDevices {
		json.NewEncoder(os.Stdout).Encode(device)
	}
}

func scrape() []Vitalstats {
	allDevices := make([]Vitalstats, 0, 400)

	doc, err := goquery.NewDocument("http://wiki.cyanogenmod.org/w/Devices")
	if err != nil {
		log.Fatal(err)
	}

	devicesSelection := doc.Find("span.device")
	devicesCount := len(devicesSelection.Nodes)
	devicesSelection.Each(func(i int, s *goquery.Selection) {
		vitals := Vitalstats{}

		vitals.Name = s.Find("span.name").Text()
		vitals.Link = s.Find("a").AttrOr("href", "")

		// lol they implode if you strcat and get a "//" after the tld
		doc2, err := goquery.NewDocument("http://wiki.cyanogenmod.org" + vitals.Link)
		if err != nil {
			log.Fatal(err)
		}
		doc2.Find("div#mw-content-text table tr").Each(func(_ int, s *goquery.Selection) {
			val := strings.TrimSpace(s.Find("td").Text())
			switch strings.TrimSpace(s.Find("th").Text()) {
			case "Power:":
				vitals.Power = val
			case "Release Date:":
				vitals.ReleaseDate = val
			case "Type:":
				vitals.Type = val
			case "CM Support:":
				vitals.CMSupport = val
			}
		})

		allDevices = append(allDevices, vitals)
		fmt.Fprintf(os.Stderr, "scanned device %d/%d\n", i, devicesCount)
	})

	return allDevices
}

type filter struct {
	name string
	fn   func(Vitalstats) bool
}
