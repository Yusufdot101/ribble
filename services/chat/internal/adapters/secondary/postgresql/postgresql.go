package postgresql

import (
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
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
		&Chat{}, &Message{}, &ChatParticipant{}, &Permission{}, &Role{}, &ChatRolePermission{}, &ChatRole{},
		&ChatBan{},
	)
	if err != nil {
		return nil, err
	}

	if err := seedRBAC(db); err != nil {
		return nil, err
	}

	return &Adapter{
		db: db,
	}, nil
}

func seedRBAC(db *gorm.DB) error {
	roles := []Role{
		{Name: domain.Admin},
		{Name: domain.Member},
		{Name: domain.Creator},
	}
	for _, r := range roles {
		if err := db.Where("name = ?", r.Name).FirstOrCreate(&r).Error; err != nil {
			return err
		}
	}

	perms := []Permission{
		{Name: domain.AddToGroup},
		{Name: domain.RemoveUserFromGroup},
		{Name: domain.BanUsers},
		{Name: domain.DeleteMessages},
		{Name: domain.SendMessage},
	}
	for _, p := range perms {
		if err := db.Where("name = ?", p.Name).FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}
	return nil
}
