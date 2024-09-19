package entities

import "time"

type EventType string

const (
	Exhibition  EventType = "Выставки"
	Conferencee EventType = "Конференции"
	RoundTable  EventType = "Круглые столы"
	Forum       EventType = "Форумы"
	Seminar     EventType = "Семинары"
	Other       EventType = "Другое"
)

type Event struct {
	UUID      string    `json:"uuid" gorm:"primaryKey,unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_date"`
	AuthorUUID  string    `json:"author_uuid"`
	Archived    bool      `json:"archived"`
	Longitude   float64   `json:"longitude"`
	Latitude    float64   `json:"latitude"`
	Address     string    `json:"address"`
	Type        EventType `json:"type" gorm:"type:event_type;not null; default:'Другое'"`
}

type EventsUser struct {
	UUID      string    `json:"uuid" gorm:"primaryKey,unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	EventUUID string `json:"event_uuid"`
	UserUUID  string `json:"user_uuid"`
	Liked     bool   `json:"liked" gorm:"default:false"`
}
