package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"london-results/internal/results"
	"london-results/internal/util"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const resultsDomain = "https://results.virginmoneylondonmarathon.com"

type scraperSpec struct {
	url     string
	scraper func(raceName, dataPath, url string)
}

var specs = map[string]scraperSpec{
	"london-2019": {
		"/2019/?page=1&event=MAS&num_results=1000&pid=list&search%5Bsex%5D=M&search%5Bage_class%5D=%25",
		bootstrapScrape,
	},
	"london-2018": {
		"/2018/?page=1&event=MAS&num_results=1000&pid=list&search%5Bage_class%5D=%25&search%5Bsex%5D=M",
		tableScrape,
	},
	"london-2017": {
		"/2017/?page=1&event=MAS&num_results=1000&pid=list&search%5Bage_class%5D=%25&search%5Bsex%5D=M",
		tableScrape,
	},
	"london-2016": {
		"/2016/?page=1&event=MAS&num_results=1000&pid=list&search%5Bage_class%5D=%25&search%5Bsex%5D=M",
		tableScrape,
	},
	"london-2015": {
		"/2015/?page=1&event=MAS&num_results=1000&pid=list&search%5Bage_class%5D=%25&search%5Bsex%5D=M",
		tableScrape,
	},
	"london-2014": {
		"/2014/?page=1&event=MAS&num_results=1000&pid=search&search%5Bsex%5D=%25&search%5Bnation%5D=%25&search_sort=place_nosex",
		tableScrape,
	},
}

func main() {

	races := availableRaces()

	dataDir := util.DefineDataPathFlag()

	raceName := flag.String("race", "",
		fmt.Sprintf("Required. Name of race to scrape. One of: %s", strings.Join(races, ", ")))

	flag.Parse()

	spec, raceExists := specs[*raceName]
	if *dataDir == "" || !raceExists {
		fmt.Print("Scrapes london marathon results from the specified year and stores in a series of CSV files.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	spec.scraper(*raceName, resultsDomain+spec.url, *dataDir)
}

func availableRaces() []string {
	n := make([]string, 0, len(specs))
	for name, _ := range specs {
		n = append(n, name)
	}

	sort.Strings(n)

	return n
}

func tableScrape(name, startingUrl, dataPath string) {

	getNthCellText := func(e *colly.HTMLElement, n int) string {
		selector := fmt.Sprintf("td:nth-child(%d)", n)
		text := e.ChildText(selector)
		if strings.Contains(text, "...") {
			expandedText := e.ChildAttr(selector+" span", "title")
			if expandedText != "" {
				text = expandedText
			}
		}

		return text
	}

	getNthCellInt := func(e *colly.HTMLElement, n int) int {
		return util.TryParseInt(getNthCellText(e, n))
	}

	getNthCellTime := func(e *colly.HTMLElement, n int) int {
		return util.TryParseTimeToSeconds(getNthCellText(e, n))
	}

	c := colly.NewCollector()
	r := results.NewResultCollector(filepath.Join(dataPath, "raw", name), name)

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {

		result := &results.Result{
			Race:              name,
			Club:              getNthCellText(e, 5),
			Number:            getNthCellInt(e, 6),
			Category:          getNthCellText(e, 7),
			HalfTimeSeconds:   getNthCellTime(e, 8),
			FinishTimeSeconds: getNthCellTime(e, 9),
		}

		r.Collect(result)
	})

	c.OnHTML("a.pages-nav-button", func(e *colly.HTMLElement) {
		if e.Text == ">" {
			err := e.Request.Visit(e.Attr("href"))
			if err != nil {
				fmt.Printf("error: %s", err)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(startingUrl)
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	r.Flush()
}

func bootstrapScrape(name, startingUrl, dataPath string) {
	c := colly.NewCollector()
	r := results.NewResultCollector(filepath.Join(dataPath, "raw", name), name)

	c.OnHTML("li.list-group-item", func(e *colly.HTMLElement) {

		place := e.ChildText("div.col-md-5 > :first-child > :first-child")
		if util.TryParseInt(place) == -1 {
			return
		}

		club := e.ChildText("div.col-md-7 > div.pull-left > :first-child > :first-child")
		club = strings.TrimPrefix(club, "Club")
		if club == "–" {
			club = ""
		}

		number := e.ChildText("div.col-md-7 > div.pull-left > :first-child > :nth-child(2)")
		number = strings.TrimPrefix(number, "Runner Number")

		category := e.ChildText("div.col-md-7 > div.pull-left > :first-child > :nth-child(3)")
		category = strings.TrimPrefix(category, "Category")

		halfTime := e.ChildText("div.col-md-7 > div.pull-left > :first-child > :nth-child(4)")
		halfTime = strings.TrimPrefix(halfTime, "Half")

		fullTime := e.ChildText("div.col-md-7 > div.pull-right > :first-child > :first-child")
		fullTime = strings.TrimPrefix(fullTime, "Finish")

		result := &results.Result{
			Race:              name,
			Club:              club,
			Number:            util.TryParseInt(number),
			Category:          category,
			HalfTimeSeconds:   util.TryParseTimeToSeconds(halfTime),
			FinishTimeSeconds: util.TryParseTimeToSeconds(fullTime),
		}

		r.Collect(result)
	})

	c.OnHTML("li.pages-nav-button", func(e *colly.HTMLElement) {
		linkText := e.ChildText("a")
		if linkText == ">" {
			err := e.Request.Visit(e.ChildAttr("a", "href"))
			if err != nil {
				fmt.Printf("error: %s", err)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(startingUrl)
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	r.Flush()
}
