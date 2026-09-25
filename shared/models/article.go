package models

import "time"

type Article struct {
	ID      int       `json:"id"`
	Link    string    `json:"link"`
	Title   string    `json:"title"`
	Text    string    `json:"text"`
	Authors []string  `json:"authors"`
	Date    time.Time `json:"date"`
	Source  string    `json:"source"`
}

type ArticleWithScore struct {
	Article Article `json:"article"`
	Score   float32 `json:"score"`
}
