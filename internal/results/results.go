package results

import (
	"fmt"
	"london-results/internal/files"
	"london-results/internal/util"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	RawResultsDir       = "raw"
	ProcessedResultsDir = "processed"
	clubsFile           = "clubs.csv"
)

type ResultCollector struct {
	folderPath string
	filePrefix string
	flushed    int
	results    Results
}

func NewResultCollector(folderPath, filePrefix string) *ResultCollector {
	return &ResultCollector{flushed: 0,
		results:    make([]*Result, 0),
		folderPath: folderPath,
		filePrefix: filePrefix}
}

func (r *ResultCollector) Collect(result *Result) {
	r.results = append(r.results, result)
	if len(r.results) > 1000 {
		r.Flush()
	}
}

func (r *ResultCollector) Flush() {
	if len(r.results) == 0 {
		return
	}

	_ = os.MkdirAll(r.folderPath, os.ModePerm)

	fileName := filepath.Join(r.folderPath, fmt.Sprintf("%s-%d.csv", r.filePrefix, r.flushed+1))
	r.results.Save(fileName)
	r.flushed++
	r.results = make([]*Result, 0)
}

type Results []*Result

func (r Results) asSlice() [][]string {
	slice := make([][]string, 0, len(r))
	for _, result := range r {
		slice = append(slice, result.asSlice())
	}

	return slice
}

func (r Results) Save(filePath string) {
	files.WriteCsv(filePath, r.asSlice())
}

type Result struct {
	Race              string
	Club              string
	Number            int
	Category          string
	HalfTimeSeconds   int
	FinishTimeSeconds int
}

func (r *Result) asSlice() []string {
	return []string{
		r.Race,
		r.Club,
		fmt.Sprintf("%d", r.Number),
		r.Category,
		fmt.Sprintf("%d", r.HalfTimeSeconds),
		fmt.Sprintf("%d", r.FinishTimeSeconds),
	}
}

func NewResultFromSlice(s []string) *Result {
	return &Result{
		Race:              s[0],
		Club:              strings.ToLower(s[1]),
		Number:            util.TryParseInt(s[2]),
		Category:          s[3],
		HalfTimeSeconds:   util.TryParseInt(s[4]),
		FinishTimeSeconds: util.TryParseInt(s[5]),
	}
}

type Clubs map[string]*Club

type Club struct {
	Name          string
	Country       string
	CanonicalName string
	Ignore        bool
	Processed     bool
}

func (c *Club) toSlice() []string {
	return []string{
		c.Name,
		c.Country,
		c.CanonicalName,
		strconv.FormatBool(c.Ignore),
		strconv.FormatBool(c.Processed),
	}
}

func (c Clubs) Save(path string) {
	s := c.asSlice()
	sort.SliceStable(s, func(i, j int) bool {
		return s[i][0] < s[j][0]
	})

	files.WriteCsv(filepath.Join(path, clubsFile), s)
}

func (c Clubs) SaveCanonical(path string) {
	savedClubs := make(map[string]bool)
	s := make([][]string, 0)
	for _, club := range c {
		if _, haveSavedClub := savedClubs[club.CanonicalName]; !haveSavedClub {
			savedClubs[club.CanonicalName] = true
			s = append(s, []string{club.CanonicalName, club.Country})
		}
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i][0] < s[j][0]
	})

	files.WriteCsv(filepath.Join(path, clubsFile), s)
}

func (c Clubs) asSlice() [][]string {
	slice := make([][]string, 0, len(c))
	for _, club := range c {
		slice = append(slice, club.toSlice())
	}

	return slice
}

func LoadClubs(path string) Clubs {
	rows := files.ReadCsv(filepath.Join(path, clubsFile))
	clubs := make(map[string]*Club, 0)
	for _, record := range rows {
		club := &Club{
			Name:          record[0],
			Country:       record[1],
			CanonicalName: record[2],
			Ignore:        util.TryParseBool(record[3], false),
			Processed:     util.TryParseBool(record[4], false),
		}
		clubs[club.Name] = club
	}

	return clubs
}
