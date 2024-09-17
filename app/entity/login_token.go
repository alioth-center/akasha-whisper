package entity

import "time"

type LoginToken struct {
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}
