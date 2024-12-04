package dto

type RoleCreateDTO struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}
