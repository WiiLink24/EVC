package main

import (
	"strconv"
	"strings"
	"time"
)

const (
	// QueryQuestionsJapan queries the questions table for questions for Japan
	QueryQuestionsJapan = `SELECT question_id, content_japanese, choice1_japanese, choice2_japanese,
						start_date, end_date
						FROM questions
						WHERE region_code = $1
						AND end_date > $2`

	// QueryQuestionsNTSC queries the questions table for questions for NTSC countries.
	QueryQuestionsNTSC = `SELECT question_id, 
						content_english, content_spanish, content_french_canada, content_portuguese,  
						choice1_english, choice1_spanish, choice1_french_canada, choice1_portuguese,
						choice2_english, choice2_spanish, choice2_french_canada, choice2_portuguese,
						start_date, end_date
						FROM questions
						WHERE region_code = $1
						AND end_date > $2`

	// QueryQuestionsPAL queries the questions table for questions for PAL countries.
	QueryQuestionsPAL = `SELECT question_id, 
						content_english, content_german, content_french, content_spanish, content_italian, content_dutch, content_portuguese,  
						choice1_english, choice1_german, choice1_french, choice1_spanish, choice1_italian, choice1_dutch, choice1_portuguese,
						choice2_english, choice2_german, choice2_french, choice2_spanish, choice2_italian, choice2_dutch, choice2_portuguese,
						start_date, end_date
						FROM questions
						WHERE region_code = $1
						AND end_date > $2`

	// QueryQuestionsWorldwide queries the questions table for worldwide questions.
	QueryQuestionsWorldwide = `SELECT question_id, 
						content_japanese, content_english, content_german, content_french, content_spanish, content_italian, content_dutch, content_portuguese, content_french_canada, 
						choice1_japanese, choice1_english, choice1_german, choice1_french, choice1_spanish, choice1_italian, choice1_dutch, choice1_portuguese, choice1_french_canada,
						choice2_japanese, choice2_english, choice2_german, choice2_french, choice2_spanish, choice2_italian, choice2_dutch, choice2_portuguese, choice2_french_canada,
						start_date, end_date
						FROM questions
						WHERE region_code = $1
						AND end_date > $2`

	// QueryResults queries the votes table for the results of a specified question.
	QueryResults = `SELECT type_cd, country_id, region_id, ans_cnt FROM votes WHERE question_id = $1`

	// QueryQuestionsForResults queries the questions table for a poll that has votes.
	QueryQuestionsForResults = `SELECT question_id FROM questions WHERE region_code = $1 AND end_date < $2 ORDER BY end_date DESC LIMIT 1`
)

// For creating the question text
var questions = make(map[LanguageCode][]string)
var response1s = make(map[LanguageCode][]string)
var response2s = make(map[LanguageCode][]string)

// National
var startDateSlice []int
var endDateSlice []int
var pollIDs []int
var nationalDetailedResults []DetailedNationalResult

// Worldwide
var worldWideStartDateSlice []int
var worldWideEndDateSlice []int
var worldWidePollIDs []int
var worldWideDetailedResults []DetailedWorldwideResult


// PrepareWorldWideResults returns the WorldWideResults for the WorldWide vote,
// as well as create a DetailedWorldwideResult slice.
func PrepareWorldWideResults() *WorldWideResult {
	var questionID int

	rows, err := pool.Query(ctx, QueryQuestionsForResults, WorldWide, time.Now().Unix())
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&questionID)
		checkError(err)
	}

	if questionID == 0 {
		// No applicable question was found
		return nil
	}

	results := WorldWideResult{
		PollID: 						 uint32(questionID),
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
	worldWideDetailedResults = make([]DetailedWorldwideResult, 34)

	// Now we query votes table
	rows, err = pool.Query(ctx, QueryResults, questionID)
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
			results.MaleVotersResponse1 += ansCNT[0]
			results.MaleVotersResponse2 += ansCNT[2]
			results.FemaleVotersResponse1 += ansCNT[1]
			results.FemaleVotersResponse2 += ansCNT[3]

			// Detailed Results
			for i, code := range countryCodes {
				if code == countryID {
					worldWideDetailedResults[i].MaleVotersResponse1 += ansCNT[0]
					worldWideDetailedResults[i].MaleVotersResponse2 += ansCNT[2]
					worldWideDetailedResults[i].FemaleVotersResponse1 += ansCNT[1]
					worldWideDetailedResults[i].FemaleVotersResponse2 += ansCNT[3]
					worldWideDetailedResults[i].CountryTableCount = 7
				}
			}
		} else if typeCD == Prediction {
			results.PredictorsResponse1 += ansCNT[0] + ansCNT[1]
			results.PredictorsResponse2 += ansCNT[2] + ansCNT[3]
		}
	}

	// Fix the country table offsets
	countryTablePos := 231
	for i := 33; i != -1; i-- {
		if worldWideDetailedResults[i].CountryTableCount == 7 {
			worldWideDetailedResults[i].CountryTableNumber = uint32(countryTablePos)
		} else {
			// Remove the current country results from the array as it is null.
			worldWideDetailedResults = append(worldWideDetailedResults[:i], worldWideDetailedResults[i + 1:]...)
		}

		countryTablePos -= 7
	}

	results.NumberOfWorldWideDetailedTables = uint8(len(worldWideDetailedResults))
	return &results
}

func PrepareNationalResults() *NationalResult  {
	var questionID int

	switch regionCodes[currentCountryCode] {
	case Japan:
		rows, err := pool.Query(ctx, QueryQuestionsForResults, Japan, time.Now().Unix())
		checkError(err)

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&questionID)
			checkError(err)
		}
	case NTSC:
		rows, err := pool.Query(ctx, QueryQuestionsForResults, NTSC, time.Now().Unix())
		checkError(err)

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&questionID)
			checkError(err)
		}
	case PAL:
		rows, err := pool.Query(ctx, QueryQuestionsForResults, PAL, time.Now().Unix())
		checkError(err)

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&questionID)
			checkError(err)
		}
	}

	if questionID == 0 {
		return nil
	}

	// Now that we know there is a question, init the nationalDetailedResults array
	nationalDetailedResults = make([]DetailedNationalResult, numberOfRegions[currentCountryCode])

	results := NationalResult{
		PollID: 							  uint32(questionID),
		MaleVotersResponse1:                  0,
		MaleVotersResponse2:                  0,
		FemaleVotersResponse1:                0,
		FemaleVotersResponse2:                0,
		PredictorsResponse1:                  0,
		PredictorsResponse2:                  0,
		ShowVoterNumberFlag:                  1,
		ShowDetailedResultsFlag:              1,
		NationalResultDetailedNumber:         numberOfRegions[currentCountryCode],
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

		ansCNT := FormatAnsCnt(strconv.FormatInt(int64(ansCNTInt), 10))
		if typeCD == Vote {
			// Main results
			results.MaleVotersResponse1 += ansCNT[0]
			results.MaleVotersResponse2 += ansCNT[2]
			results.FemaleVotersResponse1 += ansCNT[1]
			results.FemaleVotersResponse2 += ansCNT[3]

			for i := 0; i < int(numberOfRegions[currentCountryCode]); i++ {
				entryNumber := 0
				if i == regionID {
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

	return &results
}

func PrepareQuestions() {
	switch regionCodes[currentCountryCode] {
	case Japan:
		rows, err := pool.Query(ctx, QueryQuestionsJapan, Japan, time.Now().Unix())
		checkError(err)

		defer rows.Close()
		for rows.Next() {
			var questionID int

			var japaneseQuestion string
			var japaneseChoice1 string
			var japaneseChoice2 string

			var startDate	int
			var endDate		int

			err = rows.Scan(&questionID,
				&japaneseQuestion, &japaneseChoice1, &japaneseChoice2,
				&startDate, &endDate)
			checkError(err)

			// Append to the maps
			questions[Japanese] = append(questions[Japanese], japaneseQuestion)

			response1s[Japanese] = append(response1s[Japanese], japaneseChoice1)

			response2s[Japanese] = append(response2s[Japanese], japaneseChoice2)

			pollIDs = append(pollIDs, questionID)
			startDateSlice = append(startDateSlice, startDate)
			endDateSlice = append(endDateSlice, endDate)
		}
	case NTSC:
		rows, err := pool.Query(ctx, QueryQuestionsNTSC, NTSC, time.Now().Unix())
		checkError(err)

		defer rows.Close()
		for rows.Next() {
			var questionID int

			var englishQuestion string
			var spanishQuestion string
			var portugueseQuestion	string
			var frenchCanadaQuestion string

			var englishChoice1 string
			var spanishChoice1 string
			var portugueseChoice1	string
			var frenchCanadaChoice1 string

			var englishChoice2 string
			var spanishChoice2 string
			var portugueseChoice2	string
			var frenchCanadaChoice2 string

			var startDate	int
			var endDate		int

			err = rows.Scan(&questionID,
				&englishQuestion, &spanishQuestion, &frenchCanadaQuestion, &portugueseQuestion,
				&englishChoice1, &spanishChoice1, &frenchCanadaChoice1, &portugueseChoice1,
				&englishChoice2, &spanishChoice2, &frenchCanadaChoice2, &portugueseChoice2,
				&startDate, &endDate)
			checkError(err)

			// Append to the maps
			questions[English] = append(questions[English], englishQuestion)
			questions[Spanish] = append(questions[Spanish], spanishQuestion)
			questions[FrenchCanadian] = append(questions[FrenchCanadian], frenchCanadaQuestion)
			questions[Portuguese] = append(questions[Portuguese], portugueseQuestion)

			response1s[English] = append(response1s[English], englishChoice1)
			response1s[Spanish] = append(response1s[Spanish], spanishChoice1)
			response1s[FrenchCanadian] = append(response1s[FrenchCanadian], frenchCanadaChoice1)
			response1s[Portuguese] = append(response1s[Portuguese], portugueseChoice1)

			response2s[English] = append(response2s[English], englishChoice2)
			response2s[Spanish] = append(response2s[Spanish], spanishChoice2)
			response2s[FrenchCanadian] = append(response2s[FrenchCanadian], frenchCanadaChoice2)
			response2s[Portuguese] = append(response2s[Portuguese], portugueseChoice2)

			pollIDs = append(pollIDs, questionID)
			startDateSlice = append(startDateSlice, startDate)
			endDateSlice = append(endDateSlice, endDate)
		}
	case PAL:
		rows, err := pool.Query(ctx, QueryQuestionsPAL, PAL, time.Now().Unix())
		checkError(err)

		defer rows.Close()
		for rows.Next() {
			var questionID int

			var englishQuestion string
			var germanQuestion	string
			var frenchQuestion	string
			var spanishQuestion string
			var italianQuestion	string
			var dutchQuestion 	string
			var portugueseQuestion	string

			var englishChoice1 string
			var germanChoice1	string
			var frenchChoice1	string
			var spanishChoice1 string
			var italianChoice1	string
			var dutchChoice1 	string
			var portugueseChoice1	string

			var englishChoice2 string
			var germanChoice2	string
			var frenchChoice2	string
			var spanishChoice2 string
			var italianChoice2	string
			var dutchChoice2 	string
			var portugueseChoice2	string

			var startDate	int
			var endDate		int

			err = rows.Scan(&questionID,
				&englishQuestion, &germanQuestion, &frenchQuestion, &spanishQuestion, &italianQuestion, &dutchQuestion, &portugueseQuestion,
				&englishChoice1, &germanChoice1, &frenchChoice1, &spanishChoice1, &italianChoice1, &dutchChoice1, &portugueseChoice1,
				&englishChoice2, &germanChoice2, &frenchChoice2, &spanishChoice2, &italianChoice2, &dutchChoice2, &portugueseChoice2,
				&startDate, &endDate)
			checkError(err)

			// Append to the maps
			questions[English] = append(questions[English], englishQuestion)
			questions[German] = append(questions[German], germanQuestion)
			questions[French] = append(questions[French], frenchQuestion)
			questions[Spanish] = append(questions[Spanish], spanishQuestion)
			questions[Italian] = append(questions[Italian], italianQuestion)
			questions[Dutch] = append(questions[Dutch], dutchQuestion)
			questions[Portuguese] = append(questions[Portuguese], portugueseQuestion)

			response1s[English] = append(response1s[English], englishChoice1)
			response1s[German] = append(response1s[German], germanChoice1)
			response1s[French] = append(response1s[French], frenchChoice1)
			response1s[Spanish] = append(response1s[Spanish], spanishChoice1)
			response1s[Italian] = append(response1s[Italian], italianChoice1)
			response1s[Dutch] = append(response1s[Dutch], dutchChoice1)
			response1s[Portuguese] = append(response1s[Portuguese], portugueseChoice1)

			response2s[English] = append(response2s[English], englishChoice2)
			response2s[German] = append(response2s[German], germanChoice2)
			response2s[French] = append(response2s[French], frenchChoice2)
			response2s[Spanish] = append(response2s[Spanish], spanishChoice2)
			response2s[Italian] = append(response2s[Italian], italianChoice2)
			response2s[Dutch] = append(response2s[Dutch], dutchChoice2)
			response2s[Portuguese] = append(response2s[Portuguese], portugueseChoice2)

			pollIDs = append(pollIDs, questionID)
			startDateSlice = append(startDateSlice, startDate)
			endDateSlice = append(endDateSlice, endDate)
		}
	}

	// After getting all the National Questions, we can now query for worldwide
	rows, err := pool.Query(ctx, QueryQuestionsWorldwide, WorldWide, time.Now().Unix())
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		var questionID int

		var japaneseQuestion string
		var englishQuestion string
		var germanQuestion	string
		var frenchQuestion	string
		var spanishQuestion string
		var italianQuestion	string
		var dutchQuestion 	string
		var portugueseQuestion	string
		var frenchCanadaQuestion string

		var japaneseChoice1 string
		var englishChoice1 string
		var germanChoice1	string
		var frenchChoice1	string
		var spanishChoice1 string
		var italianChoice1	string
		var dutchChoice1 	string
		var portugueseChoice1	string
		var frenchCanadaChoice1 string

		var japaneseChoice2 string
		var englishChoice2 string
		var germanChoice2	string
		var frenchChoice2	string
		var spanishChoice2 string
		var italianChoice2	string
		var dutchChoice2 	string
		var portugueseChoice2	string
		var frenchCanadaChoice2 string

		var startDate int
		var endDate int

		err = rows.Scan(&questionID,
			&japaneseQuestion, &englishQuestion, &germanQuestion, &frenchQuestion, &spanishQuestion, &italianQuestion, &dutchQuestion, &portugueseQuestion, &frenchCanadaQuestion,
			&japaneseChoice1, &englishChoice1, &germanChoice1, &frenchChoice1, &spanishChoice1, &italianChoice1, &dutchChoice1, &portugueseChoice1, &frenchCanadaChoice1,
			&japaneseChoice2, &englishChoice2, &germanChoice2, &frenchChoice2, &spanishChoice2, &italianChoice2, &dutchChoice2, &portugueseChoice2, &frenchCanadaChoice2,
			&startDate, &endDate)
		checkError(err)

		// Append to slices
		questions[Japanese] = append(questions[Japanese], japaneseQuestion)
		questions[English] = append(questions[English], englishQuestion)
		questions[German] = append(questions[German], germanQuestion)
		questions[French] = append(questions[French], frenchQuestion)
		questions[Spanish] = append(questions[Spanish], spanishQuestion)
		questions[Italian] = append(questions[Italian], italianQuestion)
		questions[Dutch] = append(questions[Dutch], dutchQuestion)
		questions[Portuguese] = append(questions[Portuguese], portugueseQuestion)
		questions[FrenchCanadian] = append(questions[FrenchCanadian], frenchCanadaQuestion)

		response1s[Japanese] = append(response1s[Japanese], japaneseChoice1)
		response1s[English] = append(response1s[English], englishChoice1)
		response1s[German] = append(response1s[German], germanChoice1)
		response1s[French] = append(response1s[French], frenchChoice1)
		response1s[Spanish] = append(response1s[Spanish], spanishChoice1)
		response1s[Italian] = append(response1s[Italian], italianChoice1)
		response1s[Dutch] = append(response1s[Dutch], dutchChoice1)
		response1s[Portuguese] = append(response1s[Portuguese], portugueseChoice1)
		response1s[FrenchCanadian] = append(response1s[FrenchCanadian], frenchCanadaChoice1)

		response2s[Japanese] = append(response2s[Japanese], japaneseChoice2)
		response2s[English] = append(response2s[English], englishChoice2)
		response2s[German] = append(response2s[German], germanChoice2)
		response2s[French] = append(response2s[French], frenchChoice2)
		response2s[Spanish] = append(response2s[Spanish], spanishChoice2)
		response2s[Italian] = append(response2s[Italian], italianChoice2)
		response2s[Dutch] = append(response2s[Dutch], dutchChoice2)
		response2s[Portuguese] = append(response2s[Portuguese], portugueseChoice2)
		response2s[FrenchCanadian] = append(response2s[FrenchCanadian], frenchCanadaChoice2)

		worldWidePollIDs = append(worldWidePollIDs, questionID)
		worldWideStartDateSlice = append(worldWideStartDateSlice, startDate)
		worldWideEndDateSlice = append(worldWideEndDateSlice, endDate)
	}
}

func FormatAnsCnt(content string) []uint32 {
	uintArray := make([]uint32, 4)
	temp := make([]uint32, len(content))

	cntArray := strings.Split(content, "")

	for i, value := range cntArray {
		cnt, _ := strconv.ParseInt(value, 10, 32)
		temp[i] = uint32(cnt)
	}

	copy(uintArray[4-len(content):], temp)
	return uintArray
}

func CleanVariables() {
	questions = make(map[LanguageCode][]string)
	response1s = make(map[LanguageCode][]string)
	response2s = make(map[LanguageCode][]string)

	startDateSlice = []int{}
	endDateSlice = []int{}
	pollIDs = []int{}
	nationalDetailedResults = []DetailedNationalResult{}

	worldWideStartDateSlice = []int{}
	worldWideEndDateSlice = []int{}
	worldWidePollIDs = []int{}
	worldWideDetailedResults = []DetailedWorldwideResult{}
}
