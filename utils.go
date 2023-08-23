package main

import (
	"strconv"
	"strings"
	"time"
)

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

func ZFill(value uint8, size int) string {
	str := strconv.FormatInt(int64(value), 10)
	temp := ""

	for i := 0; i < size-len(str); i++ {
		temp += "0"
	}

	return temp + str
}

func GetSupportedLanguages(countryCode uint8) []LanguageCode {
	return countriesSupportedLanguages[countryCode]
}

func (v *Votes) GetQuestionForLanguage(question Question, language LanguageCode) string {
	switch language {
	case Japanese:
		return question.QuestionText.Japanese
	case English:
		return question.QuestionText.English
	case German:
		return question.QuestionText.German
	case French:
		return question.QuestionText.French
	case Spanish:
		return question.QuestionText.Spanish
	case Italian:
		return question.QuestionText.Italian
	case Dutch:
		return question.QuestionText.Dutch
	case Portuguese:
		return question.QuestionText.Portuguese
	case FrenchCanadian:
		return question.QuestionText.FrenchCanadian
	case Catalan:
		return question.QuestionText.Catalan
	case Russian:
		return question.QuestionText.Russian
	}

	return question.QuestionText.English
}

func (v *Votes) GetResponse1ForLanguage(question Question, language LanguageCode) string {
	switch language {
	case Japanese:
		return question.Response1.Japanese
	case English:
		return question.Response1.English
	case German:
		return question.Response1.German
	case French:
		return question.Response1.French
	case Spanish:
		return question.Response1.Spanish
	case Italian:
		return question.Response1.Italian
	case Dutch:
		return question.Response1.Dutch
	case Portuguese:
		return question.Response1.Portuguese
	case FrenchCanadian:
		return question.Response1.FrenchCanadian
	case Catalan:
		return question.Response1.Catalan
	case Russian:
		return question.Response1.Russian
	}

	return question.QuestionText.English
}

func (v *Votes) GetResponse2ForLanguage(question Question, language LanguageCode) string {
	switch language {
	case Japanese:
		return question.Response2.Japanese
	case English:
		return question.Response2.English
	case German:
		return question.Response2.German
	case French:
		return question.Response2.French
	case Spanish:
		return question.Response2.Spanish
	case Italian:
		return question.Response2.Italian
	case Dutch:
		return question.Response2.Dutch
	case Portuguese:
		return question.Response2.Portuguese
	case FrenchCanadian:
		return question.Response2.FrenchCanadian
	case Catalan:
		return question.Response2.Catalan
	case Russian:
		return question.Response2.Russian
	}

	return question.QuestionText.English
}

func GetFileType(str string) FileType {
	switch str {
	case "v":
		return Normal
	case "r":
		return Results
	default:
		return Normal
	}
}

func GetFilename() string {
	if fileType == Normal {
		return "voting.bin"
	} else {
		year := strconv.Itoa(time.Now().Year())
		month := ZFill(uint8(time.Now().Month()), 2)
		day := ZFill(uint8(time.Now().Day()), 2)
		return year + "/" + month + day + "_r.bin"
	}
}
