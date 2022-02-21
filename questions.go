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

// MakeQuestionTextInfoTable generates the metadata for questions.
func (v *Votes) MakeQuestionTextInfoTable() {
	v.Header.QuestionTextInfoTableOffset = v.GetCurrentSize()

	// We iterate over the amount of questions, then the supported languages map to create our info.
	for i := 0; i < len(questions[countriesSupportedLanguages[currentCountryCode][0]]); i++ {
		for _, code := range countriesSupportedLanguages[currentCountryCode] {
			questionTextInfoTable := QuestionTextInfo{
				LanguageCode: uint8(code),
				QuestionOffset:  0,
				Response1Offset: 0,
				Response2Offset: 0,
			}

			v.QuestionTextInfoTable = append(v.QuestionTextInfoTable, questionTextInfoTable)
		}
	}

	v.Header.NumberOfQuestions = uint8(len(v.QuestionTextInfoTable))
}

// MakeQuestionTextTable writes the text of the questions and choices.
func (v *Votes) MakeQuestionTextTable() {
	index := 0

	// We iterate over the amount of questions, then the supported languages map to write our text.
	for i := 0; i < len(questions[countriesSupportedLanguages[currentCountryCode][0]]); i++ {
		for _, code := range countriesSupportedLanguages[currentCountryCode] {
			questionText := QuestionText{
				Question:  utf16.Encode([]rune(questions[code][i])),
				Response1: utf16.Encode([]rune(response1s[code][i])),
				Response2: utf16.Encode([]rune(response2s[code][i])),
			}

			// Apply 2 bytes of padding
			questionText.Question = append(questionText.Question, uint16(0))
			questionText.Response1 = append(questionText.Response1, uint16(0))
			questionText.Response2 = append(questionText.Response2, uint16(0))

			v.QuestionText = append(v.QuestionText, questionText)

			v.QuestionTextInfoTable[index].QuestionOffset = v.GetCurrentSize() - uint32(len(questionText.Question) * 2) - uint32(len(questionText.Response1) * 2) - uint32(len(questionText.Response2) * 2)
			v.QuestionTextInfoTable[index].Response1Offset = v.GetCurrentSize() - uint32(len(questionText.Response1) * 2) - uint32(len(questionText.Response2) * 2)
			v.QuestionTextInfoTable[index].Response2Offset = v.GetCurrentSize() - uint32(len(questionText.Response2) * 2)

			index += 1
		}
	}
}