package main

import (
	"unicode/utf16"
)

// QuestionTextInfo contains the offsets for the
// question, responses and language.
type QuestionTextInfo struct {
	LanguageCode    uint8
	QuestionOffset  uint32
	Response1Offset uint32
	Response2Offset uint32
}

type QuestionText struct {
	Question  []uint16
	Response1 []uint16
	Response2 []uint16
}

// MakeQuestionsTable generates the metadata for questions.
func (v *Votes) MakeQuestionsTable() {
	v.Header.QuestionTextInfoTableOffset = v.GetCurrentSize()

	// Get all the questions for the current country.
	for _, _ = range append(append([]Question{}, nationalQuestions...), worldwideQuestion) {
		for _, language := range GetSupportedLanguages(v.currentCountryCode) {
			v.QuestionTextInfoTable = append(v.QuestionTextInfoTable, QuestionTextInfo{
				LanguageCode:    uint8(language),
				QuestionOffset:  0,
				Response1Offset: 0,
				Response2Offset: 0,
			})
		}
	}

	// Now the text
	index := 0
	for _, question := range append(append([]Question{}, nationalQuestions...), worldwideQuestion) {
		for _, language := range GetSupportedLanguages(v.currentCountryCode) {
			v.QuestionTextInfoTable[index].QuestionOffset = v.GetCurrentSize()

			questionText := QuestionText{
				Question:  utf16.Encode([]rune(v.GetQuestionForLanguage(question, language))),
				Response1: []uint16{},
				Response2: []uint16{},
			}

			// Apply 2 bytes of padding
			questionText.Question = append(questionText.Question, uint16(0))
			v.QuestionText = append(v.QuestionText, questionText)

			// Response 1
			v.QuestionTextInfoTable[index].Response1Offset = v.GetCurrentSize()
			v.QuestionText[index].Response1 = utf16.Encode([]rune(v.GetResponse1ForLanguage(question, language)))
			v.QuestionText[index].Response1 = append(v.QuestionText[index].Response1, uint16(0))

			// Response 2
			v.QuestionTextInfoTable[index].Response2Offset = v.GetCurrentSize()
			v.QuestionText[index].Response2 = utf16.Encode([]rune(v.GetResponse2ForLanguage(question, language)))
			v.QuestionText[index].Response2 = append(v.QuestionText[index].Response2, uint16(0))

			index++
		}
	}

	v.Header.NumberOfQuestions = uint8(len(v.QuestionTextInfoTable))
}
