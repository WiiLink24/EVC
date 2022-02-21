package main

import "time"

type Header struct {
	Version								uint32
	Filesize							uint32
	CRC32								uint32
	Timestamp							uint32
	CountryCode							uint8
	PublicityFlag						uint8
	QuestionVersion						uint8
	ResultVersion						uint8
	NumberOfNationalQuestions			uint8
	NationalQuestionTableOffset			uint32
	NumberOfWorldWideQuestions			uint8
	WorldWideQuestionTableOffset		uint32
	NumberOfQuestions					uint8
	QuestionTextInfoTableOffset			uint32
	NumberOfNationalResults				uint8
	NationalResultTableOffset			uint32
	NumberOfDetailedNationalResults		uint16
	DetailedNationalResultTableOffset	uint32
	NumberOfPositionTables				uint16
	PositionTableOffset					uint32
	NumberOfWorldWideResults			uint8
	WorldWideResultsTableOffset			uint32
	NumberOfDetailedWorldWideResults	uint16
	DetailedWorldWideResultTableOffset	uint32
	NumberOfCountries					uint16
	CountryTableOffset					uint32
}

func (v *Votes) MakeHeader() {
	v.Header = Header{
		Version:                      0,
		Filesize:                     0,
		CRC32:                        0,
		Timestamp:                    GenerateCurrentTimestamp(),
		CountryCode: 				  uint8(currentCountryCode),
		PublicityFlag:                0,
		QuestionVersion:              1,
		ResultVersion:                0,
		NumberOfNationalQuestions:    0,
		NationalQuestionTableOffset:  0,
		NumberOfWorldWideQuestions:   0,
		WorldWideQuestionTableOffset: 0,
		NumberOfQuestions:            0,
		QuestionTextInfoTableOffset:  0,
		NumberOfNationalResults:      0,
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
}

func GenerateCurrentTimestamp() uint32 {
	return CreateTimestamp(int(time.Now().Unix()))
}

func CreateTimestamp(time int) uint32 {
	return uint32((time - 946684800) / 60)
}
