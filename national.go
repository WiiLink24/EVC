package main

import (
	"encoding/hex"
)

// QuestionInfo contains metadata for both national and worldwide questions.
type QuestionInfo struct {
	PollID            uint32
	PollCategory1     uint8
	PollCategory2     uint8
	StartingTimestamp uint32
	EndingTimestamp   uint32
	// NumberOfSupportedLanguages is the number of languages
	// the current country supports. This will be the amount
	// of translations that will be present in the QuestionEntryTable
	// for the current question.
	NumberOfSupportedLanguages uint8
	// QuestionTableEntryNumber is the location of the question
	// in the QuestionEntryTable.
	QuestionTableEntryNumber uint32
}

// NationalResult contains the overall results for a national question.
type NationalResult struct {
	PollID                               uint32
	MaleVotersResponse1                  uint32
	MaleVotersResponse2                  uint32
	FemaleVotersResponse1                uint32
	FemaleVotersResponse2                uint32
	PredictorsResponse1                  uint32
	PredictorsResponse2                  uint32
	ShowVoterNumberFlag                  uint8
	ShowDetailedResultsFlag              uint8
	NationalResultDetailedNumber         uint8
	StartingNationalResultDetailedNumber uint32
}

// DetailedNationalResult contains the results of a question
// for a specific region.
type DetailedNationalResult struct {
	VotersResponse1Number    uint32
	VotersResponse2Number    uint32
	PositionEntryTableCount  uint8
	PositionTableEntryNumber uint32
}

// MakeNationalQuestionsTable gets the available questions from
// the database and forms the metadata.
func (v *Votes) MakeNationalQuestionsTable() {
	v.Header.NationalQuestionTableOffset = v.GetCurrentSize()
	entryNum := 0

	for _, question := range nationalQuestions {
		v.NationalQuestionTable = append(v.NationalQuestionTable, QuestionInfo{
			PollID:                     uint32(question.ID),
			PollCategory1:              uint8(question.Category),
			PollCategory2:              categoryKV[question.Category],
			StartingTimestamp:          CreateTimestamp(int(question.Time.Unix())),
			EndingTimestamp:            CreateTimestamp(int(question.Time.Unix())) + 10080,
			NumberOfSupportedLanguages: uint8(len(countriesSupportedLanguages[v.currentCountryCode])),
			QuestionTableEntryNumber:   uint32(entryNum),
		})

		entryNum += len(countriesSupportedLanguages[v.currentCountryCode])
	}

	v.Header.NumberOfNationalQuestions = uint8(len(v.NationalQuestionTable))
}

// MakeNationalResultsTable creates the results for the past six (6) national questions.
func (v *Votes) MakeNationalResultsTable() {
	result, detailed := v.PrepareNationalResults()
	v.tempDetailedResults = detailed

	if result != nil {
		v.Header.NationalResultTableOffset = v.GetCurrentSize()
		v.NationalResults = append(v.NationalResults, result...)
	}

	v.Header.NumberOfNationalResults = uint8(len(v.NationalResults))
}

// MakeDetailedNationalResultsTable creates the detailed results for the current national question.
func (v *Votes) MakeDetailedNationalResultsTable() {
	v.Header.DetailedNationalResultTableOffset = v.GetCurrentSize()

	for _, result := range v.tempDetailedResults {
		v.DetailedNationalResults = append(v.DetailedNationalResults, result...)
	}

	v.Header.NumberOfDetailedNationalResults = uint16(len(v.DetailedNationalResults))
}

// MakePositionTable creates the position table for the current country.
func (v *Votes) MakePositionTable() {
	for i, str := range positionData {
		if uint8(i) == v.currentCountryCode {
			v.Header.PositionTableOffset = v.GetCurrentSize()
			v.Header.NumberOfPositionTables = uint16(numberOfRegions[v.currentCountryCode])

			position, err := hex.DecodeString(str)
			checkError(err)

			v.PositionEntryTable = position
		}
	}
}
