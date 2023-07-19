package models

type SearchResult struct {
	Username   string  `json:"name"`
	ResultType string  `json:"result_type"`
	Score      float64 `json:"score"` // ts_rank returns a float8
	URL        string  `json:"url"`
}
