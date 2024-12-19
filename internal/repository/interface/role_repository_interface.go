package repository_interface

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type RoleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	FindAll(ctx context.Context) ([]*domain.Role, error)
	FindByID(ctx context.Context, id int64) (*domain.Role, error)
}
