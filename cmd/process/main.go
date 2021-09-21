package main

import (
	"flag"
	"fmt"
	"london-results/internal/files"
	"london-results/internal/results"
	"london-results/internal/util"
	"os"
	"path/filepath"
)

func main() {

	dataPath := util.DefineDataPathFlag()

	flag.Parse()
	if *dataPath == "" {
		fmt.Print("For all available raw results in {dataPath}/raw/{race}, copies to \n" +
			"{dataPath}/processed/{race} folder, replaces club with canonical name where \n " +
			"available (as defined in {dataPath}/raw/clubs.csv), then consolidates all results \n" +
			" to {dataPath}/processed/results.csv and writes canonical clubs to \n" +
			"{dataPath}/processed/clubs.csv.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	rawResultsPath := filepath.Join(*dataPath, results.RawResultsDir)
	processedResultsPath := filepath.Join(*dataPath, results.ProcessedResultsDir)

	normaliseClubsAndCopy(rawResultsPath, processedResultsPath)
	joinResults(processedResultsPath)
}

func normaliseClubsAndCopy(rawResultsPath, processedResultsPath string) {
	clubs := results.LoadClubs(rawResultsPath)

	files.VisitSubFolderCsvs(rawResultsPath, func(folder, file string, records [][]string) {

		processedFolder := filepath.Join(processedResultsPath, folder)
		_ = os.MkdirAll(processedFolder, os.ModePerm)

		var processedResults results.Results = make([]*results.Result, 0)
		for _, record := range records {
			result := results.NewResultFromSlice(record)
			club, foundClub := clubs[result.Club]
			if !foundClub || club.Ignore || !club.Processed {
				result.Club = ""
			} else {
				result.Club = club.CanonicalName
			}

			processedResults = append(processedResults, result)
		}

		processedResults.Save(filepath.Join(processedFolder, file))
	})

	clubs.SaveCanonical(processedResultsPath)
}

func joinResults(processedResultsPath string) {
	all := make([][]string, 0)
	files.VisitSubFolderCsvs(processedResultsPath, func(folder, file string, records [][]string) {
		all = append(all, records...)
	})

	files.WriteCsv(filepath.Join(processedResultsPath, "results.csv"), all)
}
