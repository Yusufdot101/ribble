package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	TokenString string
	UserID      uint
	Expires     time.Time
	Use         domain.TokenUse
	TokenType   domain.TokenType
}

func (a *Adapter) InsertToken(token *domain.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokenModel := &Token{
		TokenString: token.TokenString,
		UserID:      token.UserID,
		Use:         token.Use,
		Expires:     token.Expires,
		TokenType:   token.TokenType,
	}
	res := a.DB.WithContext(ctx).Create(tokenModel)
	if res.Error == nil {
		token.ID = tokenModel.ID
	}

	return res.Error
}

func (a *Adapter) GetTokenByStringAndUse(tokenString string, tokenUse domain.TokenUse) (
	*domain.Token, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokenModel := &Token{}
	res := a.DB.WithContext(ctx).First(tokenModel, "token_string = ? AND use = ? AND expires > ?", tokenString, tokenUse, time.Now())
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, res.Error
	}
	token := domain.NewToken(tokenModel.Use, tokenModel.TokenType, tokenModel.UserID, tokenModel.TokenString, tokenModel.Expires)
	return token, nil
}

func (a *Adapter) DeleteTokenByStringAndUse(tokenString string, tokenUse domain.TokenUse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := a.DB.WithContext(ctx).Where("token_string = ? AND use = ? AND expires > NOW()", tokenString, tokenUse).Delete(&Token{})

	if res.RowsAffected == 0 {
		return domain.ErrRecordNotFound
	}
	return res.Error
}
