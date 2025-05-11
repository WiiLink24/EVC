package main

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"testing"
	"time"
)

func TestGenerateAllNationalResults(t *testing.T) {
	config := GetConfig()

	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)
	pool, err = pgxpool.ConnectConfig(ctx, dbConf)
	checkError(err)

	defer pool.Close()

	fileType = Results
	locality = National

	currentTime = time.Date(2025, 5, 8, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 21362-20000; i++ {
		fmt.Printf("Starting %d\n", i)
		for _, countryCode := range countryCodes {
			Generate(countryCode)
		}
		fmt.Printf("Finished %d\n", i)

		// Get to the next question.
		if currentTime.Weekday() == time.Tuesday {
			currentTime = currentTime.AddDate(0, 0, -3)
		} else {
			currentTime = currentTime.AddDate(0, 0, -2)
		}
	}
}

func TestGenerateAllWorldwideResults(t *testing.T) {
	config := GetConfig()

	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)
	pool, err = pgxpool.ConnectConfig(ctx, dbConf)
	checkError(err)

	defer pool.Close()

	fileType = Results
	locality = Worldwide

	currentTime = time.Date(2025, 5, 16, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 21362-20000; i++ {
		PrepareWorldWideResults()

		fmt.Printf("Starting %d\n", i)
		for _, countryCode := range countryCodes {
			Generate(countryCode)
		}
		fmt.Printf("Finished %d\n", i)

		// Get to the next question.
		if currentTime.Day() == 1 {
			currentTime = currentTime.AddDate(0, -1, 0)
			currentTime = time.Date(currentTime.Year(), currentTime.Month(), 16, 0, 0, 0, 0, time.UTC)
		} else {
			currentTime = time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, time.UTC)
		}

		worldWideDetailedResults = nil
	}
}
