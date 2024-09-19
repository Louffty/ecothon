package entities

import "time"

type Admin struct {
	UUID      string    `json:"uuid" gorm:"primaryKey,unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name  string `json:"name" gorm:"unique"`
	Value bool   `json:"value" gorm:"text"`
}
