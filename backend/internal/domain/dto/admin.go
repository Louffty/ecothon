package dto

type CreateAdminValue struct {
	Name  string `json:"name" validate:"required"`
	Value bool   `json:"value"`
}

type UpdateAdminValue struct {
	UUID  string `json:"uuid" validate:"required"`
	Value bool   `json:"value"`
}
