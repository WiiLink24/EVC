package main

import (
	"bytes"
	"encoding/binary"
	"github.com/wii-tools/lz11"
	"hash/crc32"
	"io"
	"unicode/utf16"
)

type FirstData struct {
	Version            uint32
	Filesize           uint32
	CRC32              uint32
	NumberOfCountries  uint8
	CountryTableOffset uint32
	NumberOfLanguages  uint8
	LanguageTable      []uint32
	CountryTable       []CountryTable
	LanguageText       []uint16
}

type CountryTable struct {
	CountryCode                uint8
	NumberOfSupportedLanguages uint8
	SupportedLanguages         [4]LanguageCode
}

func MakeFirstData() []byte {
	buffer := new(bytes.Buffer)

	data := FirstData{
		Version:            1,
		Filesize:           0,
		CRC32:              0,
		NumberOfCountries:  uint8(len(countryCodes)),
		CountryTableOffset: uint32(18 + len(supportedLanguages)*4),
		NumberOfLanguages:  uint8(len(supportedLanguages)),
		LanguageTable:      make([]uint32, len(supportedLanguages)),
	}

	for _, code := range countryCodes {
		var languageCodes [4]LanguageCode
		copy(languageCodes[:], countriesSupportedLanguages[code])

		data.CountryTable = append(data.CountryTable, CountryTable{
			CountryCode:                code,
			NumberOfSupportedLanguages: uint8(len(countriesSupportedLanguages[code])),
			SupportedLanguages:         languageCodes,
		})
	}

	for i, language := range supportedLanguages {
		data.LanguageTable[i] = data.GetCurrentSize()
		data.LanguageText = append(data.LanguageText, utf16.Encode([]rune(language))...)
		data.LanguageText = append(data.LanguageText, 0)
	}

	data.WriteAll(buffer)

	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(buffer.Bytes()[12:], crcTable)
	data.CRC32 = checksum
	data.Filesize = uint32(buffer.Len())

	// Reset the temp buffer and compress
	buffer.Reset()
	data.WriteAll(buffer)

	compressed, err := lz11.Compress(buffer.Bytes())
	checkError(err)

	signed := SignFile(compressed)

	return signed
}

// Write writes the current values in Votes to an io.Writer method.
// This is required as Go cannot write structs with non-fixed slice sizes,
// but can write them individually.
func (f *FirstData) Write(writer io.Writer, data interface{}) {
	err := binary.Write(writer, binary.BigEndian, data)
	checkError(err)
}

func (f *FirstData) WriteAll(writer io.Writer) {
	f.Write(writer, f.Version)
	f.Write(writer, f.Filesize)
	f.Write(writer, f.CRC32)
	f.Write(writer, f.NumberOfCountries)
	f.Write(writer, f.CountryTableOffset)
	f.Write(writer, f.NumberOfLanguages)
	f.Write(writer, f.LanguageTable)
	f.Write(writer, f.CountryTable)
	f.Write(writer, f.LanguageText)
}

func (f *FirstData) GetCurrentSize() uint32 {
	buffer := bytes.NewBuffer([]byte{})
	f.WriteAll(buffer)

	return uint32(buffer.Len())
}

var supportedLanguages = []string{
	"日本語",
	"English",
	"Deutsch",
	"Français",
	"Español",
	"Italiano",
	"Nederlands",
	"Português",
	"Français",
}
