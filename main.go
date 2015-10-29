package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
	doc, err := goquery.NewDocument("http://wiki.cyanogenmod.org/w/Devices")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("span.device").Each(func(i int, s *goquery.Selection) {
		deviceName := s.Find("span.name").Text()
		deviceLink := s.Find("a").AttrOr("href", "")
		fmt.Printf("Device %d: %q >> %s", i, deviceName, deviceLink)

		// lol they implode if you strcat and get a "//" after the tld
		doc2, err := goquery.NewDocument("http://wiki.cyanogenmod.org" + deviceLink)
		if err != nil {
			log.Fatal(err)
		}
		doc2.Find("div#mw-content-text table tr").Each(func(_ int, s *goquery.Selection) {
			//maul, _ := s.Html()
			//fmt.Printf(">>>>     %s    ", maul)
			if strings.TrimSpace(s.Find("th").Text()) != "Power:" {
				//fmt.Printf(" :: %q", s.Find("th").Text())
				return
			}
			powerDesc := s.Find("td").Text()
			fmt.Printf(" ... %q", powerDesc)
		})
		fmt.Print("\n")
	})
}

func main() {
	ExampleScrape()
}
