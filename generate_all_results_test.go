package main

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"runtime"
	"sync"
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

	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < 21315-20000; i++ {
		currentTime = time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
		currentTime = currentTime.AddDate(0, 0, -2*i)

		wg := sync.WaitGroup{}
		semaphore := make(chan any, 5)

		wg.Add(len(countryCodes))
		for _, countryCode := range countryCodes {
			go func(countryCode uint8) {
				defer wg.Done()
				semaphore <- struct{}{}
				Generate(countryCode)
				<-semaphore
			}(countryCode)
		}

		wg.Wait()
	}
}
