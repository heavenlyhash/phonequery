package main

import (
	"encoding/json"
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
	Type        string // e.g. "tablet"
	CMSupport   string
}

func main() {
	doc, err := goquery.NewDocument("http://wiki.cyanogenmod.org/w/Devices")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("span.device").Each(func(i int, s *goquery.Selection) {
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
		json.NewEncoder(os.Stdout).Encode(vitals)
	})

	// remove:
	//  - 'non-removable'
}
