package tmdb

const (
	vikingsID = 44217
	flashID   = 60735
)

func (suite *TMBDTestSuite) TestGetTVDetails() {
	vikings, err := suite.GetTVDetails(vikingsID, nil)
	suite.Nil(err)
	suite.Equal("Vikings", vikings.Name)
}

func (suite *TMBDTestSuite) TestGetTVDetailsFail() {
	_, err := suite.GetTVDetails(0, nil)
	suite.Equal("The resource you requested could not be found.", err.Error())
}

func (suite *TMBDTestSuite) TestGetTVDetailsWithOptions() {
	options := make(map[string]string)
	options["language"] = "pt-BR"
	vikings, err := suite.GetTVDetails(vikingsID, options)
	suite.Nil(err)
	suite.Equal("Vikings", vikings.Name)
}

func (suite *TMBDTestSuite) TestGetTVAccountStatesFail() {
	_, err := suite.GetTVAccountStates(0, nil)
	suite.Equal("Authentication failed: You do not have permissions to access the service.", err.Error())
}

func (suite *TMBDTestSuite) TestGetTVAccountStatesWithOptions() {
	suite.Client.APIKey = ""
	options := make(map[string]string)
	options["session_id"] = "koQubnkaZ"
	_, err := suite.GetTVAccountStates(vikingsID, options)
	suite.Equal("Invalid API key: You must be granted a valid key.", err.Error())
}

func (suite *TMBDTestSuite) TestGetTVAlternativeTitles() {
	flash, err := suite.GetTVAlternativeTitles(flashID, nil)
	suite.Nil(err)
	for _, v := range flash.Results {
		if v.Iso3166_1 == "GR" {
			suite.Equal("Ο Φλας", v.Title)
		}
	}
}

func (suite *TMBDTestSuite) TestGetTVAlternativeTitlesFail() {
	_, err := suite.GetTVAlternativeTitles(0, nil)
	suite.Equal("The resource you requested could not be found.", err.Error())
}

func (suite *TMBDTestSuite) TestGetTVAlternativeTitlesWithOptions() {
	options := make(map[string]string)
	options["language"] = "pt-BR"
	flash, err := suite.GetTVAlternativeTitles(flashID, options)
	suite.Nil(err)
	suite.Equal("GR", flash.Results[0].Iso3166_1)
}

func (suite *TMBDTestSuite) TestGetTVChanges() {
	flash, err := suite.GetTVChanges(flashID, nil)
	suite.Nil(err)
	for _, v := range flash.Changes {
		for _, v := range v.Items {
			suite.NotNil(v.ID)
		}
	}
}

// The API isn't handling zero values for this end-point.
// TODO: Fix this test later.
func (suite *TMBDTestSuite) TestGetTVChangesFail() {
	_, err := suite.GetTVChanges(0, nil)
	suite.Nil(err)
}

func (suite *TMBDTestSuite) TestGetTVChangesWithOptions() {
	options := make(map[string]string)
	options["start_date"] = "2019-01-01"
	options["end_date"] = "2019-01-12"
	options["page"] = "1"
	flash, err := suite.GetTVChanges(flashID, options)
	suite.Nil(err)
	for _, v := range flash.Changes {
		for _, v := range v.Items {
			suite.NotNil(v.ID)
		}
	}
}

func (suite *TMBDTestSuite) TestGetTVContentRatings() {
	vikings, err := suite.GetTVContentRatings(vikingsID, nil)
	suite.Nil(err)
	suite.NotNil(vikings.Results[0].Iso3166_1)
}

func (suite *TMBDTestSuite) TestGetTVContentRatingsFail() {
	_, err := suite.GetTVContentRatings(0, nil)
	suite.Equal("The resource you requested could not be found.", err.Error())
}

func (suite *TMBDTestSuite) TestGetTVContentRatingsWithOptions() {
	options := make(map[string]string)
	options["language"] = "en-US"
	vikings, err := suite.GetTVContentRatings(vikingsID, options)
	suite.Nil(err)
	suite.NotNil(vikings.Results[0].Iso3166_1)
}
