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

	_ = db.Migrator().DropTable(
		&Chat{}, &Message{}, &ChatParticipant{}, &ChatRolePermission{}, &ChatRole{},
	)
	err = db.AutoMigrate(
		&Chat{}, &Message{}, &ChatParticipant{}, &Permission{}, &Role{}, &ChatRolePermission{}, &ChatRole{},
	)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		db: db,
	}, nil
}
