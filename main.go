package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/wii-tools/lz11"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
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
	currentCountryCode uint8
}

// SQL variables.
var (
	pool     *pgxpool.Pool
	ctx      = context.Background()
	fileType FileType
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Everybody Votes Channel file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func main() {
	fileType = GetFileType(os.Args[1])
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

	// Next prepare the questions for every region.
	PrepareQuestions()

	wg := sync.WaitGroup{}
	wg.Add(len(countryCodes))
	for _, countryCode := range countryCodes {
		go func(countryCode uint8) {
			defer wg.Done()

			votes := Votes{}
			votes.currentCountryCode = countryCode

			// Create the file to write to
			strCountryCode := ZFill(countryCode, 3)
			err = os.Mkdir(fmt.Sprintf("votes/%s", strCountryCode), 0755)
			if !os.IsExist(err) {
				// If the folder exists we can just continue
				checkError(err)
			}

			buffer := bytes.NewBuffer(nil)

			// Header
			votes.MakeHeader()

			if fileType == Normal {
				// Questions
				votes.MakeNationalQuestionsTable()
				votes.MakeWorldWideQuestionsTable()
				votes.MakeQuestionsTable()
			}

			// National Results
			if fileType == Normal || fileType == Results {
				// National Results
				votes.MakeNationalResultsTable()
				votes.MakeDetailedNationalResultsTable()
				votes.MakePositionTable()

				// Worldwide Results
				votes.MakeWorldWideResultsTable()
				votes.MakeDetailedWorldWideResults()
			}

			// TODO: Get national results before this, generate if results are 0
			if 1 == 2 {
				// Country Table + Text
				votes.MakeCountryInfoTable()
				votes.MakeCountryTable()
			}

			// Write to byte buffer, add the file size, calculate crc32 then write file
			votes.WriteAll(buffer)

			binary.BigEndian.PutUint32(buffer.Bytes()[4:], uint32(buffer.Len()))

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

			filename := GetFilename()

			create, err := os.Create(fmt.Sprintf("votes/%s/%s", strCountryCode, filename))
			checkError(err)

			_, err = create.Write(signed)
			checkError(err)
		}(countryCode)
	}

	wg.Wait()
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

func SignFile(contents []byte) []byte {
	buffer := bytes.NewBuffer(nil)

	// Get RSA key and sign
	rsaData, err := ioutil.ReadFile("Private.pem")
	checkError(err)

	rsaBlock, _ := pem.Decode(rsaData)

	parsedKey, err := x509.ParsePKCS1PrivateKey(rsaBlock.Bytes)
	checkError(err)

	// Hash our data then sign
	hash := sha1.New()
	_, err = hash.Write(contents)
	checkError(err)

	contentsHashSum := hash.Sum(nil)

	reader := rand.Reader
	signature, err := rsa.SignPKCS1v15(reader, parsedKey, crypto.SHA1, contentsHashSum)
	checkError(err)

	buffer.Write(make([]byte, 64))
	buffer.Write(signature)
	buffer.Write(contents)

	return buffer.Bytes()
}
