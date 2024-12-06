package dto

type RoleCreateDTO struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type RoleResponseDTO struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
