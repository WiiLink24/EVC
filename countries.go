package main

import (
	"unicode/utf16"
)

type CountryInfoTable struct {
	LanguageCode	uint8
	TextOffset		uint32
}

// languages are all the languages the Everybody Votes Channel supports.
var languages = []uint8{0, 1, 2, 3, 4, 5, 6}
var countries = make([][]string, 34)

func (v *Votes) MakeCountryInfoTable()  {
	CreateCountriesMap()
	v.Header.CountryTableOffset = v.GetCurrentSize()
	for _, _ = range countries {
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
			v.CountryTable = append(v.CountryTable, uint16(0))
			i += 1
		}
	}
}

func CreateCountriesMap() {
	countries[0] = []string{"日本", "Japan", "Japan", "Japon", "Japón", "Giappone", "Japan"}
	countries[1] = []string{"アルゼンチン", "Argentina", "Argentinien", "Argentine", "Argentina", "Argentina", "Argentinië"}
	countries[2] = []string{"ブラジル", "Brazil", "Brasilien", "Brésil", "Brasil", "Brasile", "Brazilië"}
	countries[3] = []string{"カナダ", "Canada", "Kanada", "Canada", "Canadá", "Canada", "Canada"}
	countries[4] = []string{"チリ", "Chile", "Chile", "Chili", "Chile", "Cile", "Chili"}
	countries[5] = []string{"コロンビア", "Colombia", "Kolumbien", "Colombie", "Colombia", "Colombia", "Colombia"}
	countries[6] = []string{"コスタリカ", "Costa Rica", "Costa Rica", "Costa Rica", "Costa Rica", "Costa Rica", "Costa Rica"}
	countries[7] = []string{"エクアドル", "Ecuador", "Ecuador", "Equateur", "Ecuador", "Ecuador", "Ecuador"}
	countries[8] = []string{"グアテマラ", "Guatemala", "Guatemala", "Guatemala", "Guatemala", "Guatemala", "Guatemala"}
	countries[9] = []string{"メキシコ", "Mexico", "Mexiko", "Mexique", "México", "Messico", "Mexico"}
	countries[10] = []string{"パナマ", "Panama", "Panama", "Panama", "Panamá", "Panamá", "Panama"}
	countries[11] = []string{"ペルー", "Peru", "Peru", "Pérou", "Perú", "Perù", "Peru"}
	countries[12] = []string{"アメリカ", "United States", "Vereinigte Staaten", "Etats-Unis d’Amérique", "Estados Unidos de América", "Stati Uniti d'America", "Verenigde Staten"}
	countries[13] = []string{"ベネズエラ", "Venezuela", "Venezuela", "Venezuela", "Venezuela", "Venezuela", "Venezuela"}
	countries[14] = []string{"オーストラリア", "Australia", "Australien", "Australie", "Australia", "Australia", "Australië"}
	countries[15] = []string{"オーストリア", "Austria", "Österreich", "Autriche", "Austria", "Austria", "Oostenrijk"}
	countries[16] = []string{"ベルギー", "Belgium", "Belgien", "Belgique", "Bélgica", "Belgio", "België"}
	countries[17] = []string{"デンマーク", "Denmark", "Dänemark", "Danemark", "Dinamarca", "Danimarca", "Denemarken"}
	countries[18] = []string{"フィンランド", "Finland", "Finnland", "Finlande", "Finlandia", "Finlandia", "Finland"}
	countries[19] = []string{"フランス", "France", "Frankreich", "France", "Francia", "Francia", "Frankrijk"}
	countries[20] = []string{"ドイツ", "Germany", "Deutschland", "Allemagne", "Alemania", "Germania", "Duitsland"}
	countries[21] = []string{"ギリシャ", "Greece", "Griechenland", "Grèce", "Grecia", "Grecia", "Griekenland"}
	countries[22] = []string{"アイルランド", "Ireland", "Irland", "Irlande", "Irlanda", "Irlanda", "Ierland"}
	countries[23] = []string{"イタリア", "Italy", "Italien", "Italie", "Italia", "Italia", "Italië"}
	countries[24] = []string{"ルクセンブルク", "Luxembourg", "Luxemburg", "Luxembourg", "Luxemburgo", "Lussemburgo", "Luxemburg"}
	countries[25] = []string{"オランダ", "Netherlands", "Niederlande", "Pays-Bas", "Países Bajos", "Paesi Bassi", "Nederland"}
	countries[26] = []string{"ニュージーランド", "New Zealand", "Neuseeland", "Nouvelle-Zélande", "Nueva Zelanda", "Nuova Zelanda", "Nieuw-Zeeland"}
	countries[27] = []string{"ノルウェー", "Norway", "Norwegen", "Norvège", "Noruega", "Norvegia", "Noorwegen"}
	countries[28] = []string{"ポーランド", "Poland", "Polen", "Pologne", "Polonia", "Polonia", "Polen"}
	countries[29] = []string{"ポルトガル", "Portugal", "Portugal", "Portugal", "Portugal", "Portogallo", "Portugal"}
	countries[30] = []string{"スペイン", "Spain", "Spanien", "Espagne", "España", "Spagna", "Spanje"}
	countries[31] = []string{"スウェーデン", "Sweden", "Schweden", "Suède", "Suecia", "Svezia", "Zweden"}
	countries[32] = []string{"スイス", "Switzerland", "Schweiz", "Suisse", "Suiza", "Svizzera", "Zwitserland"}
	countries[33] = []string{"イギリス", "United Kingdom", "Großbritannien", "Royaume-Uni", "Reino Unido", "Regno Unito", "Verenigd Koninkrijk"}
}