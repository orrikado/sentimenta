package adviceservice

import "time"

type Advice struct {
	UserID int       `json:"user_id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}
