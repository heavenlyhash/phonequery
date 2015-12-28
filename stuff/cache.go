package stuff

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func UseCache(cachePath string, upstreamFn Scraper) Scraper {
	cacheFile, err := os.Open(cachePath)
	// Return a scraper that just re-emits the cache if there is one
	if err == nil {
		return func() []Vitalstats {
			defer cacheFile.Close()
			dec := json.NewDecoder(cacheFile)
			var allDevices []Vitalstats
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
			return allDevices
		}
	}
	// Return a scraper that scrapes and writes the cache
	return func() []Vitalstats {
		cacheFile, err = os.OpenFile(cachePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		var allDevices []Vitalstats
		allDevices = upstreamFn()
		enc := json.NewEncoder(cacheFile)
		for _, dev := range allDevices {
			enc.Encode(dev)
		}
		return allDevices
	}
}
