package stuff

type Vitalstats struct {
	Name        string
	Link        string
	Power       string
	ReleaseDate string
	Type        string // so far, ["phone"|"tablet"|"phablet"]
	CMSupport   string
}

type Scraper func() []Vitalstats
