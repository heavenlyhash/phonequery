package stuff

type Vitalstats struct {
	Name        string // freetext
	Link        string // full url
	Power       string // more or less freetext
	BatteryRem  string // ["yes"|"no"|"unk"] (normalized by scraper making best guess)
	ReleaseDate string
	Type        string // so far, ["phone"|"tablet"|"phablet"]
	CMSupport   string // freetext (more or less csv)
}

type Scraper func() []Vitalstats
