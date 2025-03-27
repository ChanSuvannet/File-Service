package controller

import "time"

// Item represents the item structure in Go
type Item struct {
	ID           int       `json:"id"`
	Status       Status    `json:"status"`
	Name         string    `json:"name"`
	Abbre        string    `json:"abbre"`
	Icon         string    `json:"icon"`
	Size         string    `json:"size"`
	Files        int       `json:"files"`
	Secret       string    `json:"secret"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	AuthorizedIP *string   `json:"authorized_ip,omitempty"` // Optional field
	ModifiedAt   time.Time `json:"modified_at"`
}

// Status represents the status of the item in Go
type Status struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
