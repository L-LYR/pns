package model

type Message struct {
	Topic   string `json:"topic"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
