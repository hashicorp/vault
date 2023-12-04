package uicustommessages

import "time"

type Entry struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`

	Authenticated bool   `json:"authenticated"`
	MessageType   string `json:"type"`

	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	active    *bool
}

func (e *Entry) Active() bool {

}
