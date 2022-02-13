package main

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"
)

// Votes contains all the children structs needed to
// make a voting.bin file
type Votes struct {
	Header					Header
	NationalQuestionTable	[]QuestionInfo
	WorldWideQuestionTable	[]QuestionInfo
	QuestionTextInfoTable	[]QuestionTextInfo
	// NationalResults
	// NationalResultsDetailed
	// PositionEntryTable
	// WorldwideResults
	// WorldwideResultsDetailed
	QuestionText			[]QuestionText
	CountryInfoTable		[]CountryInfoTable
	CountryTable			[]uint16
}

func main() {
	create, err := os.Create("voting.bin")
	if err != nil {
		panic(err)
	}

	// Create a byte buffer for calculating crc32 later on
	buffer := bytes.NewBuffer([]byte{})

	header := Header{
		Version:                            0,
		Filesize:                           0,
		CRC32:                              0,
		Timestamp:                          GenerateTimestamp(),
		CountryCode:                        18,
		PublicityFlag:                      0,
		QuestionVersion:                    1,
		ResultVersion:                      0,
		NumberOfNationalQuestions:          0,
		NationalQuestionTableOffset:        0,
		NumberOfWorldWideQuestions:         0,
		WorldWideQuestionTableOffset:       0,
		NumberOfQuestions:                  0,
		QuestionTextInfoTableOffset:        0,
		NumberOfNationalResults:            0,
		NationalResultTableOffset:          0,
		NumberOfDetailedNationalResults:    0,
		DetailedNationalResultTableOffset:  0,
		NumberOfPositionTables:             0,
		PositionTableOffset:                0,
		NumberOfWorldWideResults:           0,
		WorldWideResultsTableOffset:        0,
		NumberOfDetailedWorldWideResults:   0,
		DetailedWorldWideResultTableOffset: 0,
		NumberOfCountries:                  0,
		CountryTableOffset:                 0,
	}

	votes := Votes{Header: header}
	votes.MakeNationalQuestionsTable()
	votes.MakeWorldWideQuestionsTable()
	votes.MakeQuestionTextInfoTable()
	votes.MakeQuestionTextTable()
	votes.MakeCountryInfoTable()
	votes.MakeCountryTable()

	// Write to byte buffer, add the file size, calculate crc32 then write file
	votes.WriteAll(buffer)

	binary.BigEndian.PutUint32(buffer.Bytes()[4:8], uint32(buffer.Len()))

	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(buffer.Bytes(), crcTable)
	votes.Header.CRC32 = checksum
	votes.Header.Filesize = uint32(buffer.Len())

	votes.WriteAll(create)
}

// Write writes the current values in Votes to an io.Writer method.
// This is required as Go cannot write structs with non-fixed slice sizes,
// but can write them individually.
func (v *Votes) Write(writer io.Writer, data interface{}) {
	err := binary.Write(writer, binary.BigEndian, data)
	if err != nil {
		panic(err)
	}
}

func (v *Votes) WriteAll(writer io.Writer)  {
	v.Write(writer, v.Header)
	v.Write(writer, v.NationalQuestionTable)
	v.Write(writer, v.WorldWideQuestionTable)
	v.Write(writer, v.QuestionTextInfoTable)

	// Go doesn't like nested slices in structs.
	for _, question := range v.QuestionText {
		v.Write(writer, question.Question)
		v.Write(writer, question.Response1)
		v.Write(writer, question.Response2)
	}

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