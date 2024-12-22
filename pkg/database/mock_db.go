package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryvasa/go-super-farmer/service_api/repository"
	"github.com/ryvasa/go-super-farmer/pkg/database/transaction"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MockDB struct {
    DB          *gorm.DB
    SqlDB       *sql.DB
    Mock        sqlmock.Sqlmock
    TxManager   transaction.TransactionManager
    BaseRepo    repository.BaseRepository
}

func NewMockDB(t *testing.T) *MockDB {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatal(err)
    }

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: sqlDB,
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        t.Fatal(err)
    }

    txManager := transaction.NewTransactionManager(gormDB)
    baseRepo := repository.NewBaseRepository(gormDB)

    return &MockDB{
        DB:        gormDB,
        SqlDB:     sqlDB,
        Mock:      mock,
        TxManager: txManager,
        BaseRepo:  baseRepo,
    }
}
