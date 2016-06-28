package stuff

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var _ Scraper = PhoneScoopScraper

var PhoneScoopScraper = func() []Vitalstats {
	allDevices := make([]Vitalstats, 0, 2500)

	doc, err := goquery.NewDocument("http://www.phonescoop.com/phones/index_all.php")
	if err != nil {
		log.Fatal(err)
	}

	devicesSelection := doc.Find("p.phone")
	devicesCount := len(devicesSelection.Nodes)
	devicesSelection.Each(func(i int, s *goquery.Selection) {
		vitals := Vitalstats{}

		vitals.Name = s.Find("a").Text()
		vitals.Link = s.Find("a").AttrOr("href", "")

		doc2, err := goquery.NewDocument("http://www.phonescoop.com/phones/" + vitals.Link)
		if err != nil {
			log.Fatal(err)
		}
		doc2.Find("div#content table.hgrid tr").Each(func(_ int, s *goquery.Selection) {
			// assume two columns; first is labels, second is values.
			rowCells := s.Find("td")
			key := strings.TrimSpace(rowCells.First().Text())
			val := strings.TrimSpace(rowCells.Next().Text())
			// fmt.Fprintf(os.Stderr, ">>> k: %q  ;; v: %q\n", key, val)
			switch key {
			case "Battery":
				vitals.Power = val
			}
			// NOT DETECTABLE HERE:
			// - release date
			// - simple "phablet" labels (though as we've discovered, that's pretty crapshoot anyway)
			// - CM version
		})
		// attempt to normalize battery type reports.  phonescoop is pretty consistent, but also freetexts it.
		switch {
		case strings.Contains(vitals.Power, "Non-removable"):
			vitals.BatteryRem = "no"
		case strings.Contains(vitals.Power, "Removable"):
			vitals.BatteryRem = "yes"
		default:
			fmt.Fprintf(os.Stderr, "> couldn't detect battery removability from %q\n", vitals.Power)
			vitals.BatteryRem = "unk"
		}

		allDevices = append(allDevices, vitals)
		fmt.Fprintf(os.Stderr, "scanned device %d/%d\n", i, devicesCount)
	})

	return allDevices
}
