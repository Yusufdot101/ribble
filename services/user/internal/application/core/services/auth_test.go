package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Yusufdot101/ribble/services/user/internal/application/core/domain"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (r *mockRepo) InsertUser(user *domain.User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *mockRepo) FindUserByProviderAndSub(provider, sub string) (*domain.User, error) {
	args := r.Called(provider, sub)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (r *mockRepo) InsertToken(token *domain.Token) error {
	args := r.Called(token)
	return args.Error(0)
}

type mockProvider struct {
	mock.Mock
}

func (p *mockProvider) GetAuthURL(state, nonce string) string {
	args := p.Called(state, nonce)
	return args.Get(0).(string)
}

func (p *mockProvider) GetUserInfo(ctx context.Context, code, nonce string) (*domain.User, error) {
	args := p.Called(ctx, code, nonce)
	return args.Get(0).(*domain.User), args.Error(1)
}

type mockTokenService struct {
	mock.Mock
}

func (tsv *mockTokenService) New(tokenType domain.TokenType, tokenUse domain.TokenUse, userID uint) (*domain.Token, error) {
	args := tsv.Called(tokenType, tokenUse, userID)
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (tsv *mockTokenService) Save(token *domain.Token) error {
	args := tsv.Called(token)
	return args.Error(0)
}

func TestHandleCallback(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		err        error
		setupMocks func(r *mockRepo, tsvc *mockTokenService, pr *mockProvider)
	}{
		{
			name: "valid flow",
			setupMocks: func(r *mockRepo, tsvc *mockTokenService, pr *mockProvider) {
				pr.On("GetUserInfo", context.Background(), "code", "nonce").Return(&domain.User{
					Sub:      "1",
					Provider: "google",
					Name:     "yusuf",
					Email:    "ym@gmail.com",
				}, nil)

				r.On("FindUserByProviderAndSub", "google", "1").Return(&domain.User{
					Sub:      "1",
					ID:       1,
					Provider: "google",
					Name:     "yusuf",
					Email:    "ym@gmail.com",
				}, nil)

				tsvc.On("New", domain.RANDOMSTRING, domain.REFRESH, uint(1)).Return(&domain.Token{
					TokenString: "tokenString",
					UserID:      1,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}, nil)

				tsvc.On("New", domain.JWT, domain.ACCESS, uint(1)).Return(&domain.Token{
					TokenString: "tokenString",
					UserID:      1,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}, nil)

				tsvc.On("Save", &domain.Token{
					TokenString: "tokenString",
					UserID:      1,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}).Return(nil)
			},
		},
		{
			name: "user not found",
			setupMocks: func(r *mockRepo, tsvc *mockTokenService, pr *mockProvider) {
				pr.On("GetUserInfo", context.Background(), "code", "nonce").Return(&domain.User{
					Sub:      "1",
					Provider: "google",
					Name:     "yusuf",
					Email:    "ym@gmail.com",
				}, nil)

				r.On("FindUserByProviderAndSub", "google", "1").Return(&domain.User{}, domain.ErrUserNotFound)

				r.On("InsertUser", &domain.User{
					Sub:      "1",
					Provider: "google",
					Name:     "yusuf",
					Email:    "ym@gmail.com",
				}).Return(nil)

				tsvc.On("New", domain.RANDOMSTRING, domain.REFRESH, uint(0)).Return(&domain.Token{
					TokenString: "tokenString",
					UserID:      0,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}, nil)

				tsvc.On("New", domain.JWT, domain.ACCESS, uint(0)).Return(&domain.Token{
					TokenString: "tokenString",
					UserID:      0,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}, nil)

				tsvc.On("Save", &domain.Token{
					TokenString: "tokenString",
					UserID:      0,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		repo := &mockRepo{}
		provider := &mockProvider{}
		tsvc := &mockTokenService{}

		tt.setupMocks(repo, tsvc, provider)

		asvc := NewAuthService(repo, provider, tsvc)
		ctx := context.Background()
		_, _, err := asvc.HandleCallback(ctx, "code", "nonce")
		if err != nil {
			if !tt.wantErr {
				t.Fatalf("unexpected error %v\n", err)
			}
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error:%v got error: %v\n", tt.err, err)
			}
			return
		}
	}
}
