package main

// WorldWideResult contains the overall results for a worldwide question.
type WorldWideResult struct {
	PollID                          uint32
	MaleVotersResponse1             uint32
	MaleVotersResponse2             uint32
	FemaleVotersResponse1           uint32
	FemaleVotersResponse2           uint32
	PredictorsResponse1             uint32
	PredictorsResponse2             uint32
	NumberOfWorldWideDetailedTables uint8
	WorldWideDetailedTableNumber    uint32
}

// DetailedWorldwideResult contains the results for a specific country.
type DetailedWorldwideResult struct {
	Unknown               uint32
	MaleVotersResponse1   uint32
	MaleVotersResponse2   uint32
	FemaleVotersResponse1 uint32
	FemaleVotersResponse2 uint32
	CountryTableCount     uint16
	CountryTableNumber    uint32
}

// MakeWorldWideQuestionsTable gets the available questions from
// the database and forms the metadata.
func (v *Votes) MakeWorldWideQuestionsTable() {
	v.Header.WorldWideQuestionTableOffset = v.GetCurrentSize()

	entryNum := len(v.NationalQuestionTable) * len(countriesSupportedLanguages[v.currentCountryCode])

	for _, question := range worldwideQuestions {
		v.WorldWideQuestionTable = append(v.WorldWideQuestionTable, QuestionInfo{
			PollID: uint32(question.ID),
			// TODO: Implement categories within db
			PollCategory1:              0,
			PollCategory2:              0,
			StartingTimestamp:          CreateTimestamp(question.StartTime),
			EndingTimestamp:            CreateTimestamp(question.EndTime),
			NumberOfSupportedLanguages: uint8(len(countriesSupportedLanguages[v.currentCountryCode])),
			QuestionTableEntryNumber:   uint32(entryNum),
		})

		entryNum += len(countriesSupportedLanguages[v.currentCountryCode])
	}

	v.Header.NumberOfWorldWideQuestions = uint8(len(v.WorldWideQuestionTable))
}

// MakeWorldWideResultsTable creates the results for the current national question.
func (v *Votes) MakeWorldWideResultsTable() {
	result := PrepareWorldWideResults()

	if result != nil {
		v.Header.WorldWideResultsTableOffset = v.GetCurrentSize()
		v.WorldwideResults = append(v.WorldwideResults, *result)
	}

	v.Header.NumberOfWorldWideResults = uint8(len(v.WorldwideResults))
}

// MakeDetailedWorldWideResults creates the detailed results for the current national question.
func (v *Votes) MakeDetailedWorldWideResults() {
	v.Header.DetailedWorldWideResultTableOffset = v.GetCurrentSize()

	v.WorldwideResultsDetailed = worldWideDetailedResults
	v.Header.NumberOfDetailedWorldWideResults = uint16(len(v.WorldwideResultsDetailed))
}
