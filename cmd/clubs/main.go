package main

import (
	"flag"
	"fmt"
	"log"
	"london-results/internal/files"
	"london-results/internal/results"
	"london-results/internal/util"
	"os"
	"path/filepath"
)

func main() {

	dataDir := util.DefineDataPathFlag()

	flag.Parse()
	if *dataDir == "" {
		fmt.Print("Parses all scraped results in {dataDir}/raw/{race} and either creates or appends newly \n" +
			"found clubs to {dataDir}/raw/clubs.csv\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	rawResultsDir := filepath.Join(*dataDir, results.RawResultsDir)
	clubs := results.LoadClubs(rawResultsDir)

	files.VisitSubFolderCsvs(rawResultsDir, func(folder, file string, records [][]string) {
		for _, record := range records {
			result := results.NewResultFromSlice(record)
			if result.Club == "" {
				continue
			}

			_, exists := clubs[result.Club]
			if !exists {
				log.Printf("Added new club %s\n", result.Club)
				newClub := &results.Club{
					Name:          result.Club,
					Country:       "other",
					CanonicalName: result.Club,
					Ignore:        false,
					Processed:     false,
				}
				clubs[newClub.Name] = newClub
			}
		}
	})

	clubs.Save(rawResultsDir)
}
