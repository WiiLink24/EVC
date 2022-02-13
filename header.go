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

func GenerateTimestamp() uint32 {
	return uint32((time.Now().Unix() - 946684800) / 60)
}
