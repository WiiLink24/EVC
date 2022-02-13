package main

import (
	"unicode/utf16"
)

// QuestionTextInfo contains the offsets for the
// question, responses and language.
type QuestionTextInfo struct {
	LanguageCode	uint8
	QuestionOffset	uint32
	Response1Offset	uint32
	Response2Offset	uint32
}

type QuestionText struct {
	Question	[]uint16
	Response1	[]uint16
	Response2	[]uint16
}

func (v *Votes) MakeQuestionTextInfoTable() {
	v.Header.QuestionTextInfoTableOffset = v.GetCurrentSize()

	questionTextInfoTable := QuestionTextInfo{
		LanguageCode:    1,
		QuestionOffset:  0,
		Response1Offset: 0,
		Response2Offset: 0,
	}

	questionTextInfoTable1 := QuestionTextInfo{
		LanguageCode:    4,
		QuestionOffset:  0,
		Response1Offset: 0,
		Response2Offset: 0,
	}


	questionTextInfoTable2 := QuestionTextInfo{
		LanguageCode:    8,
		QuestionOffset:  0,
		Response1Offset: 0,
		Response2Offset: 0,
	}

	v.QuestionTextInfoTable = append(v.QuestionTextInfoTable, questionTextInfoTable)
	v.QuestionTextInfoTable = append(v.QuestionTextInfoTable, questionTextInfoTable1)
	v.QuestionTextInfoTable = append(v.QuestionTextInfoTable, questionTextInfoTable2)
	v.Header.NumberOfQuestions = uint8(len(v.QuestionTextInfoTable))
}

func (v *Votes) MakeQuestionTextTable() {
	questionText := QuestionText{
		Question:  utf16.Encode([]rune("Is Sketch among us?")),
		Response1: utf16.Encode([]rune("Yes")),
		Response2: utf16.Encode([]rune("No")),
	}

	questionText.Question = append(questionText.Question, uint16(0))
	questionText.Response1 = append(questionText.Response1, uint16(0))
	questionText.Response2 = append(questionText.Response2, uint16(0))

	v.QuestionText = append(v.QuestionText, questionText)

	v.QuestionTextInfoTable[0].QuestionOffset = v.GetCurrentSize() - uint32(len(questionText.Question) * 2) - uint32(len(questionText.Response1) * 2) - uint32(len(questionText.Response2) * 2)
	v.QuestionTextInfoTable[0].Response1Offset = v.GetCurrentSize() - uint32(len(questionText.Response1) * 2) - uint32(len(questionText.Response2) * 2)
	v.QuestionTextInfoTable[0].Response2Offset = v.GetCurrentSize() - uint32(len(questionText.Response2) * 2)
	v.QuestionTextInfoTable[1].QuestionOffset = v.GetCurrentSize() - uint32(len(questionText.Question) * 2) - uint32(len(questionText.Response1) * 2) - uint32(len(questionText.Response2) * 2)
	v.QuestionTextInfoTable[1].Response1Offset = v.GetCurrentSize() - uint32(len(questionText.Response1) * 2) - uint32(len(questionText.Response2) * 2)
	v.QuestionTextInfoTable[1].Response2Offset = v.GetCurrentSize() - uint32(len(questionText.Response2) * 2)
	v.QuestionTextInfoTable[2].QuestionOffset = v.GetCurrentSize() - uint32(len(questionText.Question) * 2) - uint32(len(questionText.Response1) * 2) - uint32(len(questionText.Response2) * 2)
	v.QuestionTextInfoTable[2].Response1Offset = v.GetCurrentSize() - uint32(len(questionText.Response1) * 2) - uint32(len(questionText.Response2) * 2)
	v.QuestionTextInfoTable[2].Response2Offset = v.GetCurrentSize() - uint32(len(questionText.Response2) * 2)
}