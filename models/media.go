package models

import "time"

type ApiMedia struct {
	ID        string
	MediaType string
	DateTime  time.Time
	Duration  string
}
