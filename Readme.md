# London Marathon Scraper

A quick and dirty collection (no real error handling, no tests) of Go programs to scrape, sanitise and consolidate anonymous london marathon results 
and output in CSV format for analysis. The project only pulls down the mass start results, for men-only, and only back 
to 2014. Although other result pages back to 2010 are accessible by directly manipulating the URL, the results format 
does not contain the clubs - I was only really interested in using the results to find out top scottish club times 
over the years, so didn't include these.

## Data

After running, `data/processed/results.csv` and `data/processed/clubs.csv` contain the 2014-2019 mass start results. 

### Results - data/processed/results.csv and data/raw/*race*/*race*-*n*.csv

| Col Number | Description |
| --- | ----------- |
| 1 | Race name |
| 2 | Club name |
| 3 | Race number |
| 4 | Age category string |
| 4 | Halfway time in seconds |
| 5 | Finish time in seconds |

### Raw Clubs - data/raw/clubs.csv 

__See step 3 in Usage below__

| Col Number | Description |
| --- | ----------- |
| 1 | Club - club name as it was in the results  |
| 2 | Country - always set to "other" for new entries |
| 3 | CanonicalClub - what the name of the club is *meant* to be  |
| 4 | Ignore - if true, will not be present in final results |
| 5 | Processed - Set to true for new entries when updating to identify records that must be changed. Should be false to updated records.  |

### Canonical Clubs - data/raw/clubs.csv 

| Col Number | Description |
| --- | ----------- |
| 1 | Club - canonical club name |
| 2 | Country - club country |

## Usage

To generate the data:

1. Run `/cmd/scraper/main.go` to pull down data for the year or years you want.
2. Run `/cmd/clubs/main.go` to create `/data/raw/clubs.csv`
3. Manually edit `/data/raw/clubs.csv` - this is the big thing! Club entry is free form, and no two people seem 
to spell their club in the same way or in the same form. This file must be manually edited to first exclude any
clubs that are nonsense (like "not attached") by setting `Ignore` to true, then to update the `CanonicalClub`
column for every misspelled club to point to the actual club name in the same file. You can also update the club country at the same
time. Once the club is updated, the `Processed` flag should set to true.
4. Run `/cmd/clubs/process.go` to generate `data/processed/results.csv` and `data/processed/clubs.csv`. The consolidated 
results will have their clubs "fixed" according to  `/data/raw/clubs.csv`.
5. Do whatever you want to do with the CSVs.
