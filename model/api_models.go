package model

type ApiResponse struct {
	Success      bool
	Data         interface{}
	ErrorMessage string
}

type SentimentRequest struct {
	ArticleTitle   string "json:\"articleTitle\""
	ArticleSummary string "json:\"articleSummary\""
}

type SentimentResponse struct {
	ArticleTitle string  "json:\"articleTitle\""
	Score        float32 "json:\"score\""
}
