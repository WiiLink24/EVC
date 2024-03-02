package main

// countryCodes is a list of supported countries.
var countryCodes = []uint8{
	1,
	10,
	16,
	18,
	20,
	21,
	22,
	25,
	30,
	36,
	40,
	42,
	49,
	52,
	65,
	66,
	67,
	74,
	76,
	77,
	78,
	79,
	82,
	83,
	88,
	94,
	95,
	96,
	97,
	98,
	105,
	107,
	108,
	110,
}

// numberOfRegions is the amount of provinces/states/prefectures each country has
var numberOfRegions = map[uint8]uint8{
	1:   47,
	10:  24,
	16:  27,
	18:  13,
	20:  13,
	21:  33,
	22:  7,
	25:  22,
	30:  22,
	36:  32,
	40:  10,
	42:  25,
	49:  52,
	52:  25,
	65:  14,
	66:  9,
	67:  3,
	74:  17,
	76:  6,
	77:  26,
	78:  16,
	79:  13,
	82:  8,
	83:  20,
	88:  3,
	94:  12,
	95:  17,
	96:  11,
	97:  16,
	98:  7,
	105: 19,
	107: 21,
	108: 23,
	110: 5,
}

// languages are all the languages the Everybody Votes Channel supports.
var languages = []LanguageCode{Japanese, English, German, French, Spanish, Italian, Dutch}

// countries are all the countries EVC supports in all languages.
var countries = map[int][]string{
	1:   {"日本", "Japan", "Japan", "Japon", "Japón", "Giappone", "Japan"},
	10:  {"アルゼンチン", "Argentina", "Argentinien", "Argentine", "Argentina", "Argentina", "Argentinië"},
	16:  {"ブラジル", "Brazil", "Brasilien", "Brésil", "Brasil", "Brasile", "Brazilië"},
	18:  {"カナダ", "Canada", "Kanada", "Canada", "Canadá", "Canada", "Canada"},
	20:  {"チリ", "Chile", "Chile", "Chili", "Chile", "Cile", "Chili"},
	21:  {"コロンビア", "Colombia", "Kolumbien", "Colombie", "Colombia", "Colombia", "Colombia"},
	22:  {"コスタリカ", "Costa Rica", "Costa Rica", "Costa Rica", "Costa Rica", "Costa Rica", "Costa Rica"},
	25:  {"エクアドル", "Ecuador", "Ecuador", "Equateur", "Ecuador", "Ecuador", "Ecuador"},
	30:  {"グアテマラ", "Guatemala", "Guatemala", "Guatemala", "Guatemala", "Guatemala", "Guatemala"},
	36:  {"メキシコ", "Mexico", "Mexiko", "Mexique", "México", "Messico", "Mexico"},
	40:  {"パナマ", "Panama", "Panama", "Panama", "Panamá", "Panamá", "Panama"},
	42:  {"ペルー", "Peru", "Peru", "Pérou", "Perú", "Perù", "Peru"},
	49:  {"アメリカ", "United States", "Vereinigte Staaten", "Etats-Unis d’Amérique", "Estados Unidos de América", "Stati Uniti d'America", "Verenigde Staten"},
	52:  {"ベネズエラ", "Venezuela", "Venezuela", "Venezuela", "Venezuela", "Venezuela", "Venezuela"},
	65:  {"オーストラリア", "Australia", "Australien", "Australie", "Australia", "Australia", "Australië"},
	66:  {"オーストリア", "Austria", "Österreich", "Autriche", "Austria", "Austria", "Oostenrijk"},
	67:  {"ベルギー", "Belgium", "Belgien", "Belgique", "Bélgica", "Belgio", "België"},
	74:  {"デンマーク", "Denmark", "Dänemark", "Danemark", "Dinamarca", "Danimarca", "Denemarken"},
	76:  {"フィンランド", "Finland", "Finnland", "Finlande", "Finlandia", "Finlandia", "Finland"},
	77:  {"フランス", "France", "Frankreich", "France", "Francia", "Francia", "Frankrijk"},
	78:  {"ドイツ", "Germany", "Deutschland", "Allemagne", "Alemania", "Germania", "Duitsland"},
	79:  {"ギリシャ", "Greece", "Griechenland", "Grèce", "Grecia", "Grecia", "Griekenland"},
	82:  {"アイルランド", "Ireland", "Irland", "Irlande", "Irlanda", "Irlanda", "Ierland"},
	83:  {"イタリア", "Italy", "Italien", "Italie", "Italia", "Italia", "Italië"},
	88:  {"ルクセンブルク", "Luxembourg", "Luxemburg", "Luxembourg", "Luxemburgo", "Lussemburgo", "Luxemburg"},
	94:  {"オランダ", "Netherlands", "Niederlande", "Pays-Bas", "Países Bajos", "Paesi Bassi", "Nederland"},
	95:  {"ニュージーランド", "New Zealand", "Neuseeland", "Nouvelle-Zélande", "Nueva Zelanda", "Nuova Zelanda", "Nieuw-Zeeland"},
	96:  {"ノルウェー", "Norway", "Norwegen", "Norvège", "Noruega", "Norvegia", "Noorwegen"},
	97:  {"ポーランド", "Poland", "Polen", "Pologne", "Polonia", "Polonia", "Polen"},
	98:  {"ポルトガル", "Portugal", "Portugal", "Portugal", "Portugal", "Portogallo", "Portugal"},
	105: {"スペイン", "Spain", "Spanien", "Espagne", "España", "Spagna", "Spanje"},
	107: {"スウェーデン", "Sweden", "Schweden", "Suède", "Suecia", "Svezia", "Zweden"},
	108: {"スイス", "Switzerland", "Schweiz", "Suisse", "Suiza", "Svizzera", "Zwitserland"},
	110: {"イギリス", "United Kingdom", "Großbritannien", "Royaume-Uni", "Reino Unido", "Regno Unito", "Verenigd Koninkrijk"},
}

// countriesSupportedLanguages is a list of languages each country supports.
var countriesSupportedLanguages = map[uint8][]LanguageCode{
	1:   {Japanese, English},
	10:  {English, Spanish, FrenchCanadian},
	16:  {English, Spanish, Portuguese, FrenchCanadian},
	18:  {English, Spanish, FrenchCanadian},
	20:  {English, Spanish, FrenchCanadian},
	21:  {English, Spanish, FrenchCanadian},
	22:  {English, Spanish, FrenchCanadian},
	25:  {English, Spanish, FrenchCanadian},
	30:  {English, Spanish, FrenchCanadian},
	36:  {English, Spanish, FrenchCanadian},
	40:  {English, Spanish, FrenchCanadian},
	42:  {English, Spanish, FrenchCanadian},
	49:  {English, Spanish, FrenchCanadian},
	52:  {English, Spanish, FrenchCanadian},
	65:  {English},
	66:  {German, French, English, Dutch},
	67:  {German, French, English, Dutch},
	74:  {English, German},
	76:  {English},
	77:  {French, German, English},
	78:  {German, English},
	79:  {English, Spanish, Portuguese, German},
	82:  {English},
	83:  {English, Italian},
	88:  {English, German, French, Portuguese},
	94:  {English, Dutch},
	95:  {English},
	96:  {English},
	97:  {English, German},
	98:  {English, Spanish, Portuguese},
	105: {English, Spanish, Portuguese},
	107: {English},
	108: {English, German, French, Italian},
	110: {English},
}

var positionTable = map[uint8][]uint8{
	1:   {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2},
	16:  {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1},
	18:  {1, 1, 2, 1, 1, 3, 1, 1, 1, 1, 1, 4, 3},
	21:  {1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
	36:  {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	40:  {2, 0, 1, 1, 1, 0, 0, 1, 1, 2},
	49:  {1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	77:  {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
	78:  {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	83:  {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	94:  {1, 1, 1, 3, 1, 1, 1, 1, 1, 2, 1, 1},
	105: {1, 1, 1, 1, 3, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	110: {1, 2, 2, 1, 1},
}

// positionData is a list of data for the position table.
// It has something to do with the position of votes on the downloaded map.
var positionData = map[int]string{
	1:   "A2A4C828AF52B964B478AA64AA73AA87AD9BA5969B96A09EADA5A2A987947F8E78A096A5919B9B8782A591AF82AF7AB978AA6EAA6DB364AF73B96BC05AA546AA55AF4BB437B95FC358BA46C350C82DBE26C623CD2DD237C837D728E14849395A",
	16:  "A4862664E8648E1E4141C873D746CD9E7DA0B4467878B99B8746E35385BEC855C2AEE94D82DC4B6996C8A5AAE3699687E15AA064",
	18:  "87BE3CA009981EA064AAC8C3F0A8E1AAC89BD7C3D4BDAAAA50AF1E695C405649505A3C787841647D8E89",
	21:  "7C7D78739BC8695AAA5A71247D468D6B6E6E579887326946969BC896649B9119782D8C8C4BA58D4864B2677B647328194E19875A733E6E825A87",
	36:  "37508FB0786914465A5A69A54B7D98B69B9E8AAF9687E6A07DAF82918C787DA2649B91B476988BA1EBAA5F7D8CBE91A52B6F67B2A5C8C8C899AE738CC8B9D7B4",
	40:  "A05DAF7B1E7373737D5A739BAA5250823AA0",
	49:  "D25E78D252E748E1AA87917D3C7819645A64E04EDC5FC8A0BE872EE628DF18D98C5A3C46A064AA5F7869B46C9191E249DC64EB37A53FAF5087419169A08C5037D2737337735AE440DC55557D2D5AD746E254B95D7D7D2341CD55E84CC87D714BAA7878914164CD69DC3F272F9B46C3645550F0BE",
	77:  "8246DC465AB49196463CA06E28467864AA46E6E6C86E6E3296C87896C84678C88C14505A8C2D508CC8C8BE96",
	78:  "B95A64966EDC9BC8C86E5F417837AF2D7350467841AA3CBEBE919664781E8C8C",
	83:  "7D822328283C324B463264196432821E64466464786E82649682A08CA0A0BE96B9AABEBE96E63CB4",
	94:  "645AC8418C6496288214B40AAA82D223BE08A0C882B4B46E32C8788232C8",
	105: "6E5F64E6A03C3C1EF852E65FCA739AD9A7E6B4E1C8E6EBE1641E7878503CC832AA73468C1E32A0968C28781E783278327832",
	110: "B4B4738732E67846D71E82B4507D",
}

// VoteType is the type of vote sent.
// This can either be an actual vote or prediction.
type VoteType int

const (
	Vote VoteType = iota
	Prediction
)

// LanguageCode is a numerical value that represents
// a supported language in EVC.
type LanguageCode uint8

const (
	Japanese LanguageCode = iota
	English
	German
	French
	Spanish
	Italian
	Dutch
	Portuguese
	FrenchCanadian
)

// FileType is the current type of file we are generating
type FileType uint8

const (
	Normal FileType = iota
	Results
	_Question
)

// Locality is whether it is national or worldwide
type Locality uint8

const (
	Worldwide Locality = iota
	National
	All
)
