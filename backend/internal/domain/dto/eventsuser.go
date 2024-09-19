package dto

type CreateEventsUser struct {
	EventUUID string `json:"eventUUID" validate:"required"`
}

type ReturnEventsUser struct {
	Event Event     `json:"event"`
	Users []*Author `json:"users"`
}

type ReturnUsersEvents struct {
	Events []*Event `json:"events"`
}
