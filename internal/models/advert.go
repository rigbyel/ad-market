package models

import (
	"time"
)

type Advert struct {
	Id          int64
	Header      string
	Body        string
	ImageURL    string
	Price       int
	Date        time.Time
	AuthorLogin string
}
