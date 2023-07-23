package models

type UserSearchResult struct {
	Username   string  `json:"name"`
	ResultType string  `json:"result_type"`
	Score      float64 `json:"score"` // ts_rank returns a float8
	URL        string  `json:"url"`
}

type SearchResult struct {
	Name       string  `json:"name"`
	ResultType string  `json:"result_type"`
	Score      float64 `json:"score"` // ts_rank returns a float8
	URL        string  `json:"url"`
}

func (usr *UserSearchResult) ToSearchResult() *SearchResult {
	output := SearchResult{
		Name:       usr.Username,
		ResultType: usr.ResultType,
		Score:      usr.Score,
		URL:        usr.URL,
	}
	return &output
}
