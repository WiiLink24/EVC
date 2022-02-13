package main

import (
	"time"
)

func (v *Votes) MakeWorldWideQuestionsTable() {
	v.Header.WorldWideQuestionTableOffset = v.GetCurrentSize()

	worldwide := QuestionInfo{
		PollID:                     30000,
		PollCategory1:              0,
		PollCategory2:              0,
		StartingTimestamp:          uint32(time.Now().Unix()),
		EndingTimestamp:            uint32(time.Now().Unix()) + 1000,
		NumberOfSupportedLanguages: 3,
		QuestionTableEntryNumber:   0,
	}

	v.WorldWideQuestionTable = append(v.WorldWideQuestionTable, worldwide)
	v.Header.NumberOfWorldWideQuestions = uint8(len(v.WorldWideQuestionTable))
}
