package ports

import "github.com/Yusufdot101/ribble/services/chat/internal/application/core/domain"

type Repository interface {
	InsertChat(*domain.Chat) error
}
