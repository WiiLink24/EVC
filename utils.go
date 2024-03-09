package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
	"os"
	"strconv"
	"strings"
	"time"
)

func FormatAnsCnt(content string) []uint32 {
	if len(content) > 4 {
		content = content[:4]
	}

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
	// TODO: Larsen never supported Japanese for some reason. Until we are able to translate all 1000+ questions, default to English.
	case Japanese:
		return question.QuestionText.English
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
	}

	return question.QuestionText.English
}

func (v *Votes) GetResponse1ForLanguage(question Question, language LanguageCode) string {
	switch language {
	case Japanese:
		return question.Response1.English
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
	}

	return question.QuestionText.English
}

func (v *Votes) GetResponse2ForLanguage(question Question, language LanguageCode) string {
	switch language {
	case Japanese:
		return question.Response2.English
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
	}

	return question.QuestionText.English
}

func GetFileType(str string) FileType {
	switch str {
	case "v":
		return Normal
	case "r":
		return Results
	case "q":
		return _Question
	default:
		return Normal
	}
}

func GetLocality(str string) Locality {
	switch str {
	case "w":
		return Worldwide
	case "n":
		return National
	default:
		return All
	}
}

func GetTimeDifference() int {
	if locality == National {
		return 7
	} else {
		return 14
	}
}

func GetExtension() string {
	if fileType == Results {
		return "_r.bin"
	} else {
		return "_q.bin"
	}
}

func GetFilename(countryCode string) string {
	if fileType == Normal {
		return "voting.bin"
	} else {
		date := time.Now().AddDate(0, 0, -GetTimeDifference())
		year := strconv.Itoa(date.Year())
		month := ZFill(uint8(date.Month()), 2)
		day := ZFill(uint8(date.Day()), 2)

		// Create underlying directory if needed
		err := os.Mkdir(fmt.Sprintf("votes/%s/%s", countryCode, year), 0755)
		if !os.IsExist(err) {
			checkError(err)
		}

		return year + "/" + month + day + GetExtension()
	}
}

// Get the sum of a slice
func sum(arr []uint8) uint8 {
	var _sum uint8
	for _, valueInt := range arr {
		_sum += valueInt
	}
	return _sum
}

func sanitizeText(text string) string {
	var returnText string
	textList := wordwrap.WrapString(text, 50)
	for i, s := range strings.Split(textList, "\n") {
		if i == 0 {
			returnText = s
		} else {
			returnText += "\n"
			returnText += s
		}
	}

	return returnText
}

// SanitizeText wraps the text into a format suitable for EVC.
// This is a massive function but is necessary.
func (q *Question) SanitizeText() {
	// Question Text
	q.QuestionText.English = sanitizeText(q.QuestionText.English)
	q.QuestionText.German = sanitizeText(q.QuestionText.German)
	q.QuestionText.French = sanitizeText(q.QuestionText.French)
	q.QuestionText.Spanish = sanitizeText(q.QuestionText.Spanish)
	q.QuestionText.Italian = sanitizeText(q.QuestionText.Italian)
	q.QuestionText.Dutch = sanitizeText(q.QuestionText.Dutch)
	q.QuestionText.Portuguese = sanitizeText(q.QuestionText.Portuguese)
	q.QuestionText.FrenchCanadian = sanitizeText(q.QuestionText.FrenchCanadian)

	// Response 1
	q.Response1.English = sanitizeText(q.Response1.English)
	q.Response1.German = sanitizeText(q.Response1.German)
	q.Response1.French = sanitizeText(q.Response1.French)
	q.Response1.Spanish = sanitizeText(q.Response1.Spanish)
	q.Response1.Italian = sanitizeText(q.Response1.Italian)
	q.Response1.Dutch = sanitizeText(q.Response1.Dutch)
	q.Response1.Portuguese = sanitizeText(q.Response1.Portuguese)
	q.Response1.FrenchCanadian = sanitizeText(q.Response1.FrenchCanadian)

	// Response 2
	q.Response2.English = sanitizeText(q.Response2.English)
	q.Response2.German = sanitizeText(q.Response2.German)
	q.Response2.French = sanitizeText(q.Response2.French)
	q.Response2.Spanish = sanitizeText(q.Response2.Spanish)
	q.Response2.Italian = sanitizeText(q.Response2.Italian)
	q.Response2.Dutch = sanitizeText(q.Response2.Dutch)
	q.Response2.Portuguese = sanitizeText(q.Response2.Portuguese)
	q.Response2.FrenchCanadian = sanitizeText(q.Response2.FrenchCanadian)
}

func SignFile(contents []byte) []byte {
	buffer := bytes.NewBuffer(nil)

	// Get RSA key and sign
	rsaData, err := os.ReadFile("Private.pem")
	checkError(err)

	rsaBlock, _ := pem.Decode(rsaData)

	parsedKey, err := x509.ParsePKCS8PrivateKey(rsaBlock.Bytes)
	checkError(err)

	// Hash our data then sign
	hash := sha1.New()
	_, err = hash.Write(contents)
	checkError(err)

	contentsHashSum := hash.Sum(nil)

	reader := rand.Reader
	signature, err := rsa.SignPKCS1v15(reader, parsedKey.(*rsa.PrivateKey), crypto.SHA1, contentsHashSum)
	checkError(err)

	buffer.Write(make([]byte, 64))
	buffer.Write(signature)
	buffer.Write(contents)

	return buffer.Bytes()
}
