package postgresql

import (
	"fmt"

	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	DB *gorm.DB
}

func NewAdapter(databaseURL string) (*Adapter, error) {
	DB, err := gorm.Open(postgres.Open(databaseURL))
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}
	_ = DB.AutoMigrate(&User{}, &Token{}, &UserIdentity{})

	return &Adapter{
		DB: DB,
	}, nil
}

func (a *Adapter) WithTx(fn func(repo ports.Repository) error) error {
	return a.DB.Transaction(func(tx *gorm.DB) error {
		txRepo := &Adapter{DB: tx} // same struct, but with tx as the db
		return fn(txRepo)
	})
}
