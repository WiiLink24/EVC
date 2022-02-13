package main

import "time"

type QuestionInfo struct {
	PollID						uint32
	PollCategory1				uint8
	PollCategory2				uint8
	StartingTimestamp			uint32
	EndingTimestamp           	uint32
	// NumberOfSupportedLanguages is the number of languages
	// the current country supports. This will be the amount
	// of translations that will be present in the QuestionEntryTable
	// for the current question.
	NumberOfSupportedLanguages 	uint8
	// QuestionTableEntryNumber is the location of the question
	// in the QuestionEntryTable.
	QuestionTableEntryNumber	uint32
}

type NationalResult struct {
	PollID									uint32
	MaleVotersResponse1						uint32
	MaleVotersResponse2						uint32
	FemaleVotersResponse1					uint32
	FemaleVotersResponse2					uint32
	PredictorsResponse1						uint32
	PredictorsResponse2						uint32
	ShowVoterNumberFlag						uint8
	ShowDetailedResultsFlag					uint8
	NationalResultDetailedNumber			uint8
	StartingNationalResultDetailedNumber	uint32
}

func (v *Votes) MakeNationalQuestionsTable() {
	v.Header.NationalQuestionTableOffset = v.GetCurrentSize()

	nationalQuestion := QuestionInfo{
		PollID:                     30001,
		PollCategory1:              0,
		PollCategory2:              0,
		StartingTimestamp:          uint32(time.Now().Unix()),
		EndingTimestamp:            uint32(time.Now().Unix()) + 1000,
		NumberOfSupportedLanguages: 3,
		QuestionTableEntryNumber:   0,
	}

	v.NationalQuestionTable = append(v.NationalQuestionTable, nationalQuestion)
	v.Header.NumberOfNationalQuestions = uint8(len(v.NationalQuestionTable))
}