package postgresql

import (
	"fmt"

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
	err = DB.AutoMigrate(&User{}, &Token{}, &UserIdentity{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &Adapter{
		DB: DB,
	}, nil
}
