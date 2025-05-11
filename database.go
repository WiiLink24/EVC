package main

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"strconv"
	"time"
)

const (
	// QueryNationalQuestions queries the questions table for regular questions.
	QueryNationalQuestions = `SELECT * FROM questions 
							WHERE date > $1
							AND date <= CURRENT_DATE
							AND type = 'n'
							ORDER BY date
							LIMIT 3`

	// QueryQuestionsWorldwide queries the questions table for worldwide questions.
	QueryQuestionsWorldwide = `SELECT * FROM questions 
         					WHERE date > $1
           					AND date <= CURRENT_DATE
           					AND type = 'w'
         					ORDER BY date`

	// QueryApplicableNationalResults queries the questions table for national questions that have results.
	QueryApplicableNationalResults = `SELECT question_id FROM questions
							WHERE date <= $1
  							AND type = 'n'
							ORDER BY date DESC LIMIT 6`

	// QueryApplicableWorldwideResult queries the questions table for worldwide questions that have results.
	QueryApplicableWorldwideResult = `SELECT question_id FROM questions
							WHERE date <= $1
  							AND type = 'w'
							ORDER BY date DESC LIMIT 1`

	QueryVoterData = `SELECT type_cd, region_id, ans_cnt FROM votes 
                    WHERE question_id = $1 
                    AND country_id = $2`

	QueryWorldwideVoterData = `SELECT type_cd, country_id, region_id, ans_cnt FROM votes 
                    WHERE question_id = $1`
)

type LocalizedText struct {
	Japanese       string
	English        string
	German         string
	French         string
	Spanish        string
	Italian        string
	Dutch          string
	Portuguese     string
	FrenchCanadian string
}

type Question struct {
	ID           int
	QuestionText LocalizedText
	Response1    LocalizedText
	Response2    LocalizedText
	Category     int
	Time         time.Time
}

var (
	// Questions
	nationalQuestions []Question
	worldwideQuestion Question

	// Results
	worldWideDetailedResults []DetailedWorldwideResult
	worldWideResult          WorldWideResult
)

// PrepareWorldWideResults returns the WorldWideResult for the WorldWide vote,
// as well as create a DetailedWorldwideResult slice.
func PrepareWorldWideResults() {
	var questionID int

	// Worldwide polls run for 15 days. At the time this code will be executed, it should be 15 days after a
	// poll has closed.
	row := pool.QueryRow(ctx, QueryApplicableWorldwideResult, currentTime.AddDate(0, 0, -15))
	err := row.Scan(&questionID)
	if errors.Is(err, pgx.ErrNoRows) {
		return
	}

	checkError(err)

	worldWideResult = WorldWideResult{
		PollID:                          uint32(questionID),
		MaleVotersResponse1:             0,
		MaleVotersResponse2:             0,
		FemaleVotersResponse1:           0,
		FemaleVotersResponse2:           0,
		PredictorsResponse1:             0,
		PredictorsResponse2:             0,
		NumberOfWorldWideDetailedTables: 0,
		WorldWideDetailedTableNumber:    0,
	}

	// Now that we know that there is a question, init the array
	worldWideDetailedResults = make([]DetailedWorldwideResult, len(countryCodes)+1)

	// Now we query votes table
	rows, err := pool.Query(ctx, QueryWorldwideVoterData, questionID)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		var typeCD VoteType
		var countryID int
		var regionID int
		var ansCNTInt int

		err = rows.Scan(&typeCD, &countryID, &regionID, &ansCNTInt)
		checkError(err)

		ansCNT := FormatAnsCnt(strconv.FormatInt(int64(ansCNTInt), 10))
		if typeCD == Vote {
			// Main results
			worldWideResult.MaleVotersResponse1 += ansCNT[0]
			worldWideResult.MaleVotersResponse2 += ansCNT[2]
			worldWideResult.FemaleVotersResponse1 += ansCNT[1]
			worldWideResult.FemaleVotersResponse2 += ansCNT[3]

			// Detailed Results
			for i, code := range countryCodes {
				if code == uint8(countryID) {
					worldWideDetailedResults[i].MaleVotersResponse1 += ansCNT[0]
					worldWideDetailedResults[i].MaleVotersResponse2 += ansCNT[2]
					worldWideDetailedResults[i].FemaleVotersResponse1 += ansCNT[1]
					worldWideDetailedResults[i].FemaleVotersResponse2 += ansCNT[3]
					worldWideDetailedResults[i].CountryTableCount = 7
				}
			}
		} else if typeCD == Prediction {
			worldWideResult.PredictorsResponse1 += ansCNT[0] + ansCNT[1]
			worldWideResult.PredictorsResponse2 += ansCNT[2] + ansCNT[3]
		}
	}

	countryTablePos := len(countryCodes) * 7
	for i := len(countryCodes); i != -1; i-- {
		if worldWideDetailedResults[i].CountryTableCount == 7 {
			worldWideDetailedResults[i].CountryTableNumber = uint32(countryTablePos)
		} else {
			// Remove the current country results from the array as it is null.
			worldWideDetailedResults = append(worldWideDetailedResults[:i], worldWideDetailedResults[i+1:]...)
		}

		countryTablePos -= 7
	}

	worldWideResult.NumberOfWorldWideDetailedTables = uint8(len(worldWideDetailedResults))
}

func (v *Votes) PrepareNationalResults() ([]NationalResult, [][]DetailedNationalResult) {
	var nationalResults []NationalResult
	var detailedNationalResultsForResults [][]DetailedNationalResult

	// First query for applicable results.
	rows, err := pool.Query(ctx, QueryApplicableNationalResults, currentTime.AddDate(0, 0, -7))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	checkError(err)

	index := 0
	defer rows.Close()
	for rows.Next() {
		// Now get voter data.
		var questionID int
		err = rows.Scan(&questionID)
		checkError(err)

		// Allocate space for the detailed results and the base result metadata
		nationalDetailedResults := make([]DetailedNationalResult, numberOfRegions[v.currentCountryCode])
		results := NationalResult{
			PollID:                               uint32(questionID),
			MaleVotersResponse1:                  0,
			MaleVotersResponse2:                  0,
			FemaleVotersResponse1:                0,
			FemaleVotersResponse2:                0,
			PredictorsResponse1:                  0,
			PredictorsResponse2:                  0,
			ShowVoterNumberFlag:                  1,
			ShowDetailedResultsFlag:              0,
			NationalResultDetailedNumber:         numberOfRegions[v.currentCountryCode],
			StartingNationalResultDetailedNumber: uint32(numberOfRegions[v.currentCountryCode] * uint8(index)),
		}

		voterRows, err := pool.Query(ctx, QueryVoterData, questionID, v.currentCountryCode)
		checkError(err)

		for voterRows.Next() {
			var typeCD VoteType
			var regionID int
			var ansCNTInt int

			err = voterRows.Scan(&typeCD, &regionID, &ansCNTInt)
			checkError(err)

			// Show the country map if we got a position table
			if _, ok := positionTable[v.currentCountryCode]; ok {
				results.ShowDetailedResultsFlag = 1
			}

			ansCNT := FormatAnsCnt(strconv.FormatInt(int64(ansCNTInt), 10))
			if typeCD == Vote {
				// Main results
				results.MaleVotersResponse1 += ansCNT[0]
				results.MaleVotersResponse2 += ansCNT[2]
				results.FemaleVotersResponse1 += ansCNT[1]
				results.FemaleVotersResponse2 += ansCNT[3]

				for i := 0; i < int(numberOfRegions[v.currentCountryCode]); i++ {
					// Nintendo made the region ID start at index 1, with that being the country.
					if i == regionID-2 {
						nationalDetailedResults[i].VotersResponse1Number += ansCNT[0] + ansCNT[1]
						nationalDetailedResults[i].VotersResponse2Number += ansCNT[2] + ansCNT[3]
						if _, ok := positionTable[v.currentCountryCode]; ok {
							nationalDetailedResults[i].PositionEntryTableCount = positionTable[v.currentCountryCode][i]
						} else {
							nationalDetailedResults[i].PositionEntryTableCount = 0
						}
					}
					if _, ok := positionTable[v.currentCountryCode]; ok {
						nationalDetailedResults[i].PositionTableEntryNumber = uint32(sum(positionTable[v.currentCountryCode][:i]))
					}
				}
			} else if typeCD == Prediction {
				results.PredictorsResponse1 += ansCNT[0] + ansCNT[1]
				results.PredictorsResponse2 += ansCNT[2] + ansCNT[3]
			}
		}

		index++
		nationalResults = append(nationalResults, results)
		detailedNationalResultsForResults = append(detailedNationalResultsForResults, nationalDetailedResults)
		voterRows.Close()

		if fileType == Results {
			// Only one result is required for this file type.
			break
		}
	}

	return nationalResults, detailedNationalResultsForResults
}

func PrepareNationalQuestions() {
	rows, err := pool.Query(ctx, QueryNationalQuestions, currentTime.AddDate(0, 0, -7))
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		question := Question{}
		err = rows.Scan(&question.ID,
			&question.QuestionText.English, &question.QuestionText.German, &question.QuestionText.French,
			&question.QuestionText.Spanish, &question.QuestionText.Italian, &question.QuestionText.Dutch,
			&question.QuestionText.Portuguese, &question.QuestionText.FrenchCanadian,
			&question.Response1.English, &question.Response1.German, &question.Response1.French,
			&question.Response1.Spanish, &question.Response1.Italian, &question.Response1.Dutch,
			&question.Response1.Portuguese, &question.Response1.FrenchCanadian,
			&question.Response2.English, &question.Response2.German, &question.Response2.French,
			&question.Response2.Spanish, &question.Response2.Italian, &question.Response2.Dutch,
			&question.Response2.Portuguese, &question.Response2.FrenchCanadian, nil, &question.Category,
			&question.Time,
		)
		checkError(err)

		// Apply wordwrap for each question
		question.SanitizeText()

		// Finally append to the list of national questions.
		nationalQuestions = append(nationalQuestions, question)
	}
}

func PrepareWorldWideQuestion() {
	row := pool.QueryRow(ctx, QueryQuestionsWorldwide, currentTime.AddDate(0, 0, -15))

	question := Question{}
	err := row.Scan(&question.ID,
		&question.QuestionText.English, &question.QuestionText.German, &question.QuestionText.French,
		&question.QuestionText.Spanish, &question.QuestionText.Italian, &question.QuestionText.Dutch,
		&question.QuestionText.Portuguese, &question.QuestionText.FrenchCanadian,
		&question.Response1.English, &question.Response1.German, &question.Response1.French,
		&question.Response1.Spanish, &question.Response1.Italian, &question.Response1.Dutch,
		&question.Response1.Portuguese, &question.Response1.FrenchCanadian,
		&question.Response2.English, &question.Response2.German, &question.Response2.French,
		&question.Response2.Spanish, &question.Response2.Italian, &question.Response2.Dutch,
		&question.Response2.Portuguese, &question.Response2.FrenchCanadian, nil, &question.Category,
		&question.Time,
	)
	checkError(err)

	// Apply wordwrap for each question
	question.SanitizeText()

	// Finally assign as our worldwide question.
	worldwideQuestion = question
}
