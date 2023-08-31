package main

import (
	"github.com/jackc/pgx/v4"
	"strconv"
	"time"
)

const (
	// QueryQuestions queries the questions table for regular questions.
	QueryQuestions = `SELECT question_id, 
						content_japanese, content_english, content_german, content_french, content_spanish, content_italian, content_dutch, content_portuguese, content_french_canada, content_catalan, content_russian,
						choice1_japanese, choice1_english, choice1_german, choice1_french, choice1_spanish, choice1_italian, choice1_dutch, choice1_portuguese, choice1_french_canada, choice1_catalan, choice1_russian,
						choice2_japanese, choice2_english, choice2_german, choice2_french, choice2_spanish, choice2_italian, choice2_dutch, choice2_portuguese, choice2_french_canada, choice2_catalan, choice2_russian,
						start_date, end_date
						FROM questions
						WHERE end_date > $1
						AND start_date <= $1  
						AND worldwide = false
						ORDER BY end_date ASC LIMIT 3`

	// QueryQuestionsWorldwide queries the questions table for worldwide questions.
	QueryQuestionsWorldwide = `SELECT question_id, 
						content_japanese, content_english, content_german, content_french, content_spanish, content_italian, content_dutch, content_portuguese, content_french_canada, content_catalan, content_russian,
						choice1_japanese, choice1_english, choice1_german, choice1_french, choice1_spanish, choice1_italian, choice1_dutch, choice1_portuguese, choice1_french_canada, choice1_catalan, choice1_russian,
						choice2_japanese, choice2_english, choice2_german, choice2_french, choice2_spanish, choice2_italian, choice2_dutch, choice2_portuguese, choice2_french_canada, choice2_catalan, choice2_russian,
						start_date, end_date
						FROM questions
						WHERE end_date > $1
						AND start_date <= $1  
						AND worldwide = true`

	// QueryResults queries the votes table for the results of a specified question.
	QueryResults = `SELECT type_cd, country_id, region_id, ans_cnt FROM votes WHERE question_id = $1`

	BaseQueryWorldwide = `SELECT question_id FROM questions WHERE end_date < $1 AND worldwide = true ORDER BY end_date DESC LIMIT 1`

	BaseQueryNational = `SELECT question_id FROM questions WHERE end_date < $1 AND worldwide = false ORDER BY end_date DESC LIMIT 1`
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
	Catalan        string
	Russian        string
}

type Question struct {
	ID           int
	QuestionText LocalizedText
	Response1    LocalizedText
	Response2    LocalizedText
	StartTime    int
	EndTime      int
}

var questions []Question
var worldwideQuestions []Question
var worldWideDetailedResults []DetailedWorldwideResult
var worldWideResult WorldWideResult

// PrepareWorldWideResults returns the WorldWideResult for the WorldWide vote,
// as well as create a DetailedWorldwideResult slice.
func PrepareWorldWideResults() {
	var questionID int

	row := pool.QueryRow(ctx, BaseQueryWorldwide, time.Now().Unix())
	err := row.Scan(&questionID)
	if err == pgx.ErrNoRows {
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
	rows, err := pool.Query(ctx, QueryResults, questionID)
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

func (v *Votes) PrepareNationalResults() (*NationalResult, []DetailedNationalResult) {
	var questionID int

	row := pool.QueryRow(ctx, BaseQueryNational, time.Now().Unix())
	err := row.Scan(&questionID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	checkError(err)

	// Now that we know there is a question, init the nationalDetailedResults array
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
		ShowDetailedResultsFlag:              1,
		NationalResultDetailedNumber:         numberOfRegions[v.currentCountryCode],
		StartingNationalResultDetailedNumber: 0,
	}

	// Now query the votes table
	rows, err := pool.Query(ctx, QueryResults, questionID)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		var typeCD VoteType
		var countryID int
		var regionID int
		var ansCNTInt int

		err = rows.Scan(&typeCD, &countryID, &regionID, &ansCNTInt)
		checkError(err)

		if countryID != int(v.currentCountryCode) {
			continue
		}

		ansCNT := FormatAnsCnt(strconv.FormatInt(int64(ansCNTInt), 10))
		if typeCD == Vote {
			// Main results
			results.MaleVotersResponse1 += ansCNT[0]
			results.MaleVotersResponse2 += ansCNT[2]
			results.FemaleVotersResponse1 += ansCNT[1]
			results.FemaleVotersResponse2 += ansCNT[3]

			for i := 0; i < int(numberOfRegions[v.currentCountryCode]); i++ {
				entryNumber := 0
				// Nintendo made the region ID start at index 1, with that being the country.
				if i == regionID-2 {
					nationalDetailedResults[i].VotersResponse1Number = ansCNT[0] + ansCNT[1]
					nationalDetailedResults[i].VotersResponse2Number = ansCNT[2] + ansCNT[3]
					nationalDetailedResults[i].PositionTableEntryNumber = uint32(entryNumber)
					nationalDetailedResults[i].PositionEntryTableCount = 1
				}

				entryNumber += 1
			}
		} else if typeCD == Prediction {
			results.PredictorsResponse1 += ansCNT[0] + ansCNT[1]
			results.PredictorsResponse2 += ansCNT[2] + ansCNT[3]
		}
	}

	return &results, nationalDetailedResults
}

func PrepareQuestions() {
	// Query all questions except for Worldwide.
	rows, err := pool.Query(ctx, QueryQuestions, time.Now().Unix())
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		var questionID int

		var japaneseQuestion string
		var englishQuestion string
		var germanQuestion string
		var frenchQuestion string
		var spanishQuestion string
		var italianQuestion string
		var dutchQuestion string
		var portugueseQuestion string
		var frenchCanadaQuestion string
		var catalanQuestion string
		var russianQuestion string

		var japaneseChoice1 string
		var englishChoice1 string
		var germanChoice1 string
		var frenchChoice1 string
		var spanishChoice1 string
		var italianChoice1 string
		var dutchChoice1 string
		var portugueseChoice1 string
		var frenchCanadaChoice1 string
		var catalanChoice1 string
		var russianChoice1 string

		var japaneseChoice2 string
		var englishChoice2 string
		var germanChoice2 string
		var frenchChoice2 string
		var spanishChoice2 string
		var italianChoice2 string
		var dutchChoice2 string
		var portugueseChoice2 string
		var frenchCanadaChoice2 string
		var catalanChoice2 string
		var russianChoice2 string

		var startDate int
		var endDate int

		err = rows.Scan(&questionID,
			&japaneseQuestion, &englishQuestion, &germanQuestion, &frenchQuestion, &spanishQuestion, &italianQuestion, &dutchQuestion, &portugueseQuestion, &frenchCanadaQuestion, &catalanQuestion, &russianQuestion,
			&japaneseChoice1, &englishChoice1, &germanChoice1, &frenchChoice1, &spanishChoice1, &italianChoice1, &dutchChoice1, &portugueseChoice1, &frenchCanadaChoice1, &catalanChoice1, &russianChoice1,
			&japaneseChoice2, &englishChoice2, &germanChoice2, &frenchChoice2, &spanishChoice2, &italianChoice2, &dutchChoice2, &portugueseChoice2, &frenchCanadaChoice2, &catalanChoice2, &russianChoice2,
			&startDate, &endDate)
		checkError(err)

		questions = append(questions, Question{
			ID: questionID,
			QuestionText: LocalizedText{
				Japanese:       japaneseQuestion,
				English:        englishQuestion,
				German:         germanQuestion,
				French:         frenchQuestion,
				Spanish:        spanishQuestion,
				Italian:        italianQuestion,
				Dutch:          dutchQuestion,
				Portuguese:     portugueseQuestion,
				FrenchCanadian: frenchCanadaQuestion,
				Catalan:        catalanQuestion,
				Russian:        russianQuestion,
			},
			Response1: LocalizedText{
				Japanese:       japaneseChoice1,
				English:        englishChoice1,
				German:         germanChoice1,
				French:         frenchChoice1,
				Spanish:        spanishChoice1,
				Italian:        italianChoice1,
				Dutch:          dutchChoice1,
				Portuguese:     portugueseChoice1,
				FrenchCanadian: frenchCanadaChoice1,
				Catalan:        catalanChoice1,
				Russian:        russianChoice1,
			},
			Response2: LocalizedText{
				Japanese:       japaneseChoice2,
				English:        englishChoice2,
				German:         germanChoice2,
				French:         frenchChoice2,
				Spanish:        spanishChoice2,
				Italian:        italianChoice2,
				Dutch:          dutchChoice2,
				Portuguese:     portugueseChoice2,
				FrenchCanadian: frenchCanadaChoice2,
				Catalan:        catalanChoice2,
				Russian:        russianChoice2,
			},
			StartTime: startDate,
			EndTime:   endDate,
		})
	}

	// After getting all the National Questions, we can now query for worldwide
	row := pool.QueryRow(ctx, QueryQuestionsWorldwide, time.Now().Unix())
	var questionID int

	var japaneseQuestion string
	var englishQuestion string
	var germanQuestion string
	var frenchQuestion string
	var spanishQuestion string
	var italianQuestion string
	var dutchQuestion string
	var portugueseQuestion string
	var frenchCanadaQuestion string
	var catalanQuestion string
	var russianQuestion string

	var japaneseChoice1 string
	var englishChoice1 string
	var germanChoice1 string
	var frenchChoice1 string
	var spanishChoice1 string
	var italianChoice1 string
	var dutchChoice1 string
	var portugueseChoice1 string
	var frenchCanadaChoice1 string
	var catalanChoice1 string
	var russianChoice1 string

	var japaneseChoice2 string
	var englishChoice2 string
	var germanChoice2 string
	var frenchChoice2 string
	var spanishChoice2 string
	var italianChoice2 string
	var dutchChoice2 string
	var portugueseChoice2 string
	var frenchCanadaChoice2 string
	var catalanChoice2 string
	var russianChoice2 string

	var startDate int
	var endDate int

	err = row.Scan(&questionID,
		&japaneseQuestion, &englishQuestion, &germanQuestion, &frenchQuestion, &spanishQuestion, &italianQuestion, &dutchQuestion, &portugueseQuestion, &frenchCanadaQuestion, &catalanQuestion, &russianQuestion,
		&japaneseChoice1, &englishChoice1, &germanChoice1, &frenchChoice1, &spanishChoice1, &italianChoice1, &dutchChoice1, &portugueseChoice1, &frenchCanadaChoice1, &catalanChoice1, &russianChoice1,
		&japaneseChoice2, &englishChoice2, &germanChoice2, &frenchChoice2, &spanishChoice2, &italianChoice2, &dutchChoice2, &portugueseChoice2, &frenchCanadaChoice2, &catalanChoice2, &russianChoice2,
		&startDate, &endDate)
	if err == pgx.ErrNoRows {
		// No question, return and let that be it.
		return
	} else {
		checkError(err)
	}

	worldwideQuestions = append(worldwideQuestions, Question{
		ID: questionID,
		QuestionText: LocalizedText{
			Japanese:       japaneseQuestion,
			English:        englishQuestion,
			German:         germanQuestion,
			French:         frenchQuestion,
			Spanish:        spanishQuestion,
			Italian:        italianQuestion,
			Dutch:          dutchQuestion,
			Portuguese:     portugueseQuestion,
			FrenchCanadian: frenchCanadaQuestion,
			Catalan:        catalanQuestion,
			Russian:        russianQuestion,
		},
		Response1: LocalizedText{
			Japanese:       japaneseChoice1,
			English:        englishChoice1,
			German:         germanChoice1,
			French:         frenchChoice1,
			Spanish:        spanishChoice1,
			Italian:        italianChoice1,
			Dutch:          dutchChoice1,
			Portuguese:     portugueseChoice1,
			FrenchCanadian: frenchCanadaChoice1,
			Catalan:        catalanChoice1,
			Russian:        russianChoice1,
		},
		Response2: LocalizedText{
			Japanese:       japaneseChoice2,
			English:        englishChoice2,
			German:         germanChoice2,
			French:         frenchChoice2,
			Spanish:        spanishChoice2,
			Italian:        italianChoice2,
			Dutch:          dutchChoice2,
			Portuguese:     portugueseChoice2,
			FrenchCanadian: frenchCanadaChoice2,
			Catalan:        catalanChoice2,
			Russian:        russianChoice2,
		},
		StartTime: startDate,
		EndTime:   endDate,
	})
}
