package util

import (
	"flag"
	"strconv"
	"time"
)

func DefineDataPathFlag() *string {
	return flag.String("dataPath", "",
		"Required. Path to directory where scraped and processed results should be stored.")
}

func TryParseInt(s string) int {
	intVal, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return -1
	}

	return int(intVal)
}

func TryParseBool(s string, d bool) bool {
	boolVal, err := strconv.ParseBool(s)
	if err != nil {
		return d
	}

	return boolVal
}

func TryParseTimeToSeconds(s string) int {
	timeVal, err := time.Parse("03:04:05", s)
	if err != nil {
		return -1
		//log.Fatalf("Unable to parse half time %s: %s", s, err)
	}

	return timeVal.Second() + (timeVal.Minute() * 60) + (timeVal.Hour() * 60 * 60)
}
