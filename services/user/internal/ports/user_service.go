package ports

import "github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"

type UserService interface {
	GetUsersByEmail(email string) ([]*domain.User, error)
}
