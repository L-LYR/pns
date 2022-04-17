package model

type Message struct {
	Title   string `json:"title" copier:"must,nopanic"`
	Content string `json:"content" copier:"must,nopanic"`
}
