package models

import "time"

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Publisher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Authors       []Author  `json:"authors"`
	Blurb         string    `json:"blurb"`
	Image         string    `json:"image"`
	GoodreadsLink string    `json:"goodreads_link"`
	CustomLink    string    `json:"custom_link"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
}

type ReadingStatus string

const (
	StatusReading  ReadingStatus = "reading"
	StatusFinished ReadingStatus = "finished"
	StatusPaused   ReadingStatus = "paused"
)

type ReadingProgress struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	BookID    string        `json:"book_id"`
	Pages     int           `json:"pages"`
	PagesRead int           `json:"pages_read"`
	Status    ReadingStatus `json:"status"`
	CreateAt  time.Time     `json:"create_at"`
	UpdateAt  time.Time     `json:"update_at"`
}

type BookAuthor struct {
	BookID   string `json:"book_id"`
	AuthorID string `json:"author_id"`
}
