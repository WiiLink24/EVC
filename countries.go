package main

import (
	"unicode/utf16"
)

type CountryInfoTable struct {
	LanguageCode LanguageCode
	TextOffset   uint32
}

func (v *Votes) MakeCountryInfoTable() {
	v.Header.CountryTableOffset = v.GetCurrentSize()
	for range countries {
		for _, code := range languages {
			country := CountryInfoTable{
				LanguageCode: code,
				TextOffset:   0,
			}

			v.CountryInfoTable = append(v.CountryInfoTable, country)
		}
	}

	v.Header.NumberOfCountries = uint16(len(v.CountryInfoTable))
}

func (v *Votes) MakeCountryTable() {
	i := 0

	for _, strings := range countries {
		for _, country := range strings {
			v.CountryInfoTable[i].TextOffset = v.GetCurrentSize()
			v.CountryTable = append(v.CountryTable, utf16.Encode([]rune(country))...)

			// Apply 2 bytes padding
			v.CountryTable = append(v.CountryTable, uint16(0))
			i += 1
		}
	}
}
