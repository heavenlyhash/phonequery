package stuff

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var _ Scraper = CMWikiScraper

var CMWikiScraper = func() []Vitalstats {
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
