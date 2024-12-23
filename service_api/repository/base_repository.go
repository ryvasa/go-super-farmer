package repository

import (
	"context"

	"github.com/ryvasa/go-super-farmer/pkg/database/transaction"
	"gorm.io/gorm"
)

type BaseRepositoryImpl struct {
	db *gorm.DB
}

type BaseRepository interface {
	DB(ctx context.Context) *gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &BaseRepositoryImpl{db: db}
}

func (r *BaseRepositoryImpl) DB(ctx context.Context) *gorm.DB {
	if tx := transaction.GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}
