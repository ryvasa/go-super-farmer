package dto

type AuthDTO struct {
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type AuthResponseDTO struct {
	User  *UserResponseDTO `json:"user"`
	Token string           `json:"token"`
}

type AuthVerifyEmailDTO struct {
	Email string `json:"email" validate:"required,email,min=3,max=255"`
}
