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

	for _, device := range allDevices {
		if device.Type == "phablet" {
			continue
		}
		if device.Type == "tablet" {
			continue
		}
		if strings.Contains(device.Power, "non-removable") {
			continue
		}
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
			switch strings.TrimSpace(s.Find("th").Text()) {
			case "Power:":
				vitals.Power = s.Find("td").Text()
			case "Release Date:":
				vitals.ReleaseDate = s.Find("td").Text()
			case "Type:":
				vitals.Type = s.Find("td").Text()
			case "CM Support:":
				vitals.CMSupport = s.Find("td").Text()
			}
		})

		allDevices = append(allDevices, vitals)
		fmt.Fprintf(os.Stderr, "scanned device %d/%d\n", i, devicesCount)
	})

	return allDevices
}
