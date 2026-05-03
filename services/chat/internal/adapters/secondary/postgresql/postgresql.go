package postgresql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(databaseURL string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(databaseURL))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&Chat{}, &Message{}, &ChatParticipant{}, &ChatRolePermission{}, &ChatRole{},
		&ChatBan{},
	)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		db: db,
	}, nil
}
