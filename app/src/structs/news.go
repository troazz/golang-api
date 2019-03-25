package structs

import "time"

// News struct
type News struct {
	ID      int       `json:"id,omitempty"`
	Author  string    `json:"author,omitempty"`
	Body    string    `json:"body,omitempty"`
	Created time.Time `json:"created,omitempty"`
}
