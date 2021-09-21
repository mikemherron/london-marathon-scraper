package files

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func VisitSubFolderCsvs(dataPath string, cb func(folder, file string, records [][]string)) {
	dirEntries, err := ioutil.ReadDir(dataPath)
	if err != nil {
		log.Fatalln("couldn't read directory", err)
	}

	for _, e := range dirEntries {
		if !e.IsDir() {
			continue
		}

		files, err := ioutil.ReadDir(filepath.Join(dataPath, e.Name()))
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range files {
			cb(e.Name(), c.Name(), ReadCsv(filepath.Join(dataPath, e.Name(), c.Name())))
		}
	}
}

func ReadCsv(filePath string) [][]string {
	csvFile, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make([][]string, 0)
		}
		log.Fatalln("couldn't open the csv file", err)
	}

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln("couldn't read the csv file", err)
	}

	return records
}

func WriteCsv(filePath string, records [][]string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("could not create file: %s", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		log.Fatalf("could not write file rows: %s", err)
	}
}
