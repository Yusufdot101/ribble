package services

import (
	"regexp"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
)

var EmailRX = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

func isValidEmail(email string) bool {
	return EmailRX.MatchString(email)
}

type UserService struct {
	repo ports.Repository
}

func NewUserService(repo ports.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (usvc *UserService) GetUsersByEmail(email string) ([]*domain.User, error) {
	// if len(email) > 0 && !isValidEmail(email) {
	// 	return nil, domain.ErrInvalidEmail
	// }

	return usvc.repo.FindUsersByEmail(email)
}
