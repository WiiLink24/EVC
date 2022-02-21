package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"hash/crc32"
	"io"
	"log"
	"os"
	"strconv"
)

// Votes contains all the children structs needed to
// make a voting.bin file.
type Votes struct {
	Header						Header
	NationalQuestionTable		[]QuestionInfo
	WorldWideQuestionTable		[]QuestionInfo
	QuestionTextInfoTable		[]QuestionTextInfo
	NationalResults				[]NationalResult
	DetailedNationalResults		[]DetailedNationalResult
	PositionEntryTable			[]byte
	WorldwideResults			[]WorldWideResult
	WorldwideResultsDetailed	[]DetailedWorldwideResult
	QuestionText				[]QuestionText
	CountryInfoTable			[]CountryInfoTable
	CountryTable				[]uint16
}

// SQL variables.
var (
	pool *pgxpool.Pool
	ctx = context.Background()
)

// currentCountryCode stores the country code we are currently generating for.
var currentCountryCode int

func checkError(err error) {
	if err != nil {
		log.Fatalf("Everybody Votes Channel file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func main() {
	// First, we will create a housing directory for all our files.
	err := os.Mkdir("votes", 0755)
	if !os.IsExist(err) {
		// If the folder exists we can just continue
		checkError(err)
	}

	for _, regionCode := range countryCodes {
		currentCountryCode = regionCode

		// Create the file to write to
		strCountryCode := strconv.FormatInt(int64(currentCountryCode), 10)
		temp := ""

		for i := 0; i < 3-len(strCountryCode); i++ {
			temp += "0"
		}

		err = os.Mkdir(fmt.Sprintf("votes/%s", temp + strCountryCode), 0755)
		if !os.IsExist(err) {
			// If the folder exists we can just continue
			checkError(err)
		}

		create, err := os.Create(fmt.Sprintf("votes/%s/VOTING.BIN", temp + strCountryCode))
		checkError(err)

		// Create a byte buffer for calculating crc32 later on
		buffer := bytes.NewBuffer([]byte{})

		// Get config
		config := GetConfig()

		// Start SQL
		dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
		dbConf, err := pgxpool.ParseConfig(dbString)
		pool, err = pgxpool.ConnectConfig(ctx, dbConf)

		var votes Votes

		// Header
		PrepareQuestions()
		votes.MakeHeader()

		// Questions
		votes.MakeNationalQuestionsTable()
		votes.MakeWorldWideQuestionsTable()
		votes.MakeQuestionTextInfoTable()

		// National Results
		votes.MakeNationalResultsTable()
		votes.MakeDetailedNationalResultsTable()
		votes.MakePositionTable()

		// Worldwide Results
		votes.MakeWorldWideResultsTable()
		votes.MakeDetailedWorldWideResults()

		// Country Table + Text
		votes.MakeCountryInfoTable()
		votes.MakeQuestionTextTable()
		votes.MakeCountryTable()

		// Write to byte buffer, add the file size, calculate crc32 then write file
		votes.WriteAll(buffer)

		binary.BigEndian.PutUint32(buffer.Bytes()[4:8], uint32(buffer.Len()))

		crcTable := crc32.MakeTable(crc32.IEEE)
		checksum := crc32.Checksum(buffer.Bytes()[12:], crcTable)
		votes.Header.CRC32 = checksum
		votes.Header.Filesize = uint32(buffer.Len())

		votes.WriteAll(create)

		CleanVariables()
	}
}

// Write writes the current values in Votes to an io.Writer method.
// This is required as Go cannot write structs with non-fixed slice sizes,
// but can write them individually.
func (v *Votes) Write(writer io.Writer, data interface{}) {
	err := binary.Write(writer, binary.BigEndian, data)
	checkError(err)
}

func (v *Votes) WriteAll(writer io.Writer)  {
	v.Write(writer, v.Header)
	
	// Questions
	v.Write(writer, v.NationalQuestionTable)
	v.Write(writer, v.WorldWideQuestionTable)
	v.Write(writer, v.QuestionTextInfoTable)

	// National Results
	v.Write(writer, v.NationalResults)
	v.Write(writer, v.DetailedNationalResults)
	v.Write(writer, v.PositionEntryTable)

	// Worldwide Results
	v.Write(writer, v.WorldwideResults)
	v.Write(writer, v.WorldwideResultsDetailed)

	v.Write(writer, v.CountryInfoTable)

	// Go doesn't like nested slices in structs.
	for _, question := range v.QuestionText {
		v.Write(writer, question.Question)
		v.Write(writer, question.Response1)
		v.Write(writer, question.Response2)
	}

	v.Write(writer, v.CountryTable)
}

// GetCurrentSize returns the current size of our Votes struct.
// This is useful for calculating the current offset of Votes.
func (v *Votes) GetCurrentSize() uint32 {
	buffer := bytes.NewBuffer([]byte{})
	v.WriteAll(buffer)

	return uint32(buffer.Len())
}