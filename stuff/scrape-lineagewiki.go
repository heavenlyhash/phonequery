package stuff

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//  LineageOS has their shit WAY more together than previous generations of android/CM!
//
//  All their data is actually available in YAML in GIT!
//    https://github.com/LineageOS/lineage_wiki/blob/master/_data/devices/bacon.yml
//
//  Of course, I finished patching this scraper before seeing that... so.
//
//  Unfortunately, fields like "battery" are still no mode standardized in their
//  source.  (Maybe we should write a linter for them, since that's actually
//  conceivable when they've already got a reasonably formatted datastore.)

var _ Scraper = LineageWikiScraper

var LineageWikiScraper = func() []Vitalstats {
	allDevices := make([]Vitalstats, 0, 2500)

	doc, err := goquery.NewDocument("https://wiki.lineageos.org/devices/")
	if err != nil {
		log.Fatal(err)
	}

	devicesSelection := doc.Find("table.device tr")
	devicesCount := len(devicesSelection.Nodes)
	devicesSelection.Each(func(i int, s *goquery.Selection) {
		vitals := Vitalstats{}

		vitals.Name = s.Find("td").Text()
		vitals.Link = s.Find("td a").AttrOr("href", "")
		vitals.Type = strings.ToLower(s.Find("td").Next().Next().Text())

		doc2, err := goquery.NewDocument("https://wiki.lineageos.org/" + vitals.Link)
		if err != nil {
			log.Fatal(err)
		}
		// LOL at this selector, they've not marked their stats table in any way,
		//  so we're going entirely on the happenstance of their *grid layout* here.
		doc2.Find("div.col-md-4 table tr").Each(func(_ int, s *goquery.Selection) {
			// assume two columns; first is labels, second is values.
			rowCells := s.Find("td")
			key := strings.TrimSpace(rowCells.First().Text())
			val := strings.TrimSpace(rowCells.Next().Text())
			// fmt.Fprintf(os.Stderr, ">>> k: %q  ;; v: %q\n", key, val)
			switch key {
			case "Battery":
				vitals.Power = val
			case "Released":
				vitals.ReleaseDate = val
			case "Supported versions":
				vitals.CMSupport = val
			}
		})
		// attempt to normalize battery type reports.  phonescoop is pretty consistent, but also freetexts it.
		switch {
		case strings.Contains(strings.ToLower(vitals.Power), "non-removable"):
			vitals.BatteryRem = "no"
		case strings.Contains(strings.ToLower(vitals.Power), "removable"):
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
