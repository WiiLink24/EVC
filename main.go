package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/wii-tools/lz11"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

// Votes contains all the children structs needed to
// make a voting.bin file.
type Votes struct {
	Header                   Header
	NationalQuestionTable    []QuestionInfo
	WorldWideQuestionTable   []QuestionInfo
	QuestionTextInfoTable    []QuestionTextInfo
	QuestionText             []QuestionText
	NationalResults          []NationalResult
	DetailedNationalResults  []DetailedNationalResult
	PositionEntryTable       []byte
	WorldwideResults         []WorldWideResult
	WorldwideResultsDetailed []DetailedWorldwideResult
	CountryInfoTable         []CountryInfoTable
	CountryTable             []uint16

	// Static values
	currentCountryCode  uint8
	tempDetailedResults [][]DetailedNationalResult
}

// SQL variables.
var (
	pool        *pgxpool.Pool
	ctx         = context.Background()
	fileType    FileType
	locality    Locality
	currentTime time.Time
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Everybody Votes Channel file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func main() {
	currentTime = time.Now()
	err := os.WriteFile("votes/first_data.bin", MakeFirstData(), 0666)
	checkError(err)

	fileType = GetFileType(os.Args[1])
	if len(os.Args) >= 3 {
		locality = GetLocality(os.Args[2])
	} else {
		locality = All
	}

	// Get config
	config := GetConfig()

	// Start SQL
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Username, config.Password, config.DatabaseAddress, config.DatabaseName)
	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)
	pool, err = pgxpool.ConnectConfig(ctx, dbConf)
	checkError(err)

	defer pool.Close()

	// First, we will create a housing directory for all our files.
	err = os.Mkdir("votes", 0755)
	if !os.IsExist(err) {
		checkError(err)
	}

	if fileType == Normal {
		// voting.bin requires all questions and all applicable results.
		PrepareNationalQuestions()
		PrepareWorldWideQuestion()

		PrepareWorldWideResults()
	} else if fileType == Results {
		// National results will generate themselves
		if locality == Worldwide {
			PrepareWorldWideResults()
		}
	} else if fileType == _Question {
		if locality == Worldwide {
			PrepareWorldWideQuestion()
		} else {
			PrepareNationalQuestions()
		}
	}

	for _, countryCode := range countryCodes {
		// NOTE: Usually for bulk files, I want to use sync.WaitGroup.
		// However, it seems that the amount of files we generate for this
		// will not give us faster speeds, in fact the opposite has occurred with deadlocks at unknown positions.
		Generate(countryCode)
	}
}

func Generate(countryCode uint8) {
	votes := Votes{}
	votes.currentCountryCode = countryCode

	// Create the file to write to
	strCountryCode := ZFill(countryCode, 3)
	err := os.Mkdir(fmt.Sprintf("votes/%s", strCountryCode), 0755)
	if !os.IsExist(err) {
		// If the folder exists we can just continue
		checkError(err)
	}

	buffer := bytes.NewBuffer(nil)

	// Header
	votes.MakeHeader()

	if fileType == Normal || fileType == _Question {
		// Questions
		if len(nationalQuestions) != 0 {
			votes.MakeNationalQuestionsTable()
		}

		if worldwideQuestion.ID != 0 {
			votes.MakeWorldWideQuestionsTable()
		}

		if worldwideQuestion.ID != 0 || len(nationalQuestions) != 0 {
			votes.MakeQuestionsTable()
		}
	}

	// National Results
	if fileType == Normal || fileType == Results {
		if locality != Worldwide {
			votes.MakeNationalResultsTable()
			votes.MakeDetailedNationalResultsTable()
			votes.MakePositionTable()
		}

		if locality != National {
			votes.MakeWorldWideResultsTable()
			votes.MakeDetailedWorldWideResults()
		}
	}

	if (fileType == Normal || fileType == Results) && locality != National {
		// Country Table + Text
		votes.MakeCountryInfoTable()
		votes.MakeCountryTable()
	}

	// Write to byte buffer, add the file size, calculate crc32 then write file
	votes.WriteAll(buffer)

	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(buffer.Bytes()[12:], crcTable)
	votes.Header.CRC32 = checksum
	votes.Header.Filesize = uint32(buffer.Len())

	// Reset the temp buffer and compress
	buffer.Reset()
	votes.WriteAll(buffer)

	compressed, err := lz11.Compress(buffer.Bytes())
	checkError(err)

	signed := SignFile(compressed)

	filename := GetFilename(strCountryCode)

	err = os.WriteFile(fmt.Sprintf("votes/%s/%s", strCountryCode, filename), signed, 0666)
	checkError(err)
}

// Write writes the current values in Votes to an io.Writer method.
// This is required as Go cannot write structs with non-fixed slice sizes,
// but can write them individually.
func (v *Votes) Write(writer io.Writer, data interface{}) {
	err := binary.Write(writer, binary.BigEndian, data)
	checkError(err)
}

func (v *Votes) WriteAll(writer io.Writer) {
	v.Write(writer, v.Header)

	// Questions
	v.Write(writer, v.NationalQuestionTable)
	v.Write(writer, v.WorldWideQuestionTable)
	v.Write(writer, v.QuestionTextInfoTable)

	// Go doesn't like nested slices in structs.
	for _, question := range v.QuestionText {
		v.Write(writer, question.Question)
		v.Write(writer, question.Response1)
		v.Write(writer, question.Response2)
	}

	// National Results
	v.Write(writer, v.NationalResults)
	v.Write(writer, v.DetailedNationalResults)
	v.Write(writer, v.PositionEntryTable)

	// Worldwide Results
	v.Write(writer, v.WorldwideResults)
	v.Write(writer, v.WorldwideResultsDetailed)

	v.Write(writer, v.CountryInfoTable)
	v.Write(writer, v.CountryTable)
}

// GetCurrentSize returns the current size of our Votes struct.
// This is useful for calculating the current offset of Votes.
func (v *Votes) GetCurrentSize() uint32 {
	buffer := bytes.NewBuffer([]byte{})
	v.WriteAll(buffer)

	return uint32(buffer.Len())
}
