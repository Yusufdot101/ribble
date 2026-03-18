package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Yusufdot101/ribble/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ribble/services/user/internal/ports"
)

func TestHandleCallback(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		err        error
		setupMocks func(r *ports.MockRepository, tsvc *ports.MockTokenService, pr *ports.MockOAuthProvider)
	}{
		{
			name: "valid flow",
			setupMocks: func(r *ports.MockRepository, tsvc *ports.MockTokenService, pr *ports.MockOAuthProvider) {
				pr.EXPECT().GetUserInfo(context.Background(), "code", "nonce").Return(&domain.User{
					Sub:      "1",
					Provider: "google",
					Name:     "yusuf",
					Email:    "ym@gmail.com",
				}, nil)

				r.EXPECT().FindUserByProviderAndSub("google", "1").Return(&domain.User{
					Sub:      "1",
					ID:       1,
					Provider: "google",
					Name:     "yusuf",
					Email:    "ym@gmail.com",
				}, nil)

				tsvc.EXPECT().New(domain.UUID, domain.REFRESH, uint(1)).Return(&domain.Token{
					TokenString: "tokenString",
					UserID:      1,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}, nil)

				tsvc.EXPECT().New(domain.JWT, domain.ACCESS, uint(1)).Return(&domain.Token{
					TokenString: "tokenString",
					UserID:      1,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}, nil)

				tsvc.EXPECT().Save(&domain.Token{
					TokenString: "tokenString",
					UserID:      1,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}).Return(nil)
			},
		},
		{
			name: "user not found",
			setupMocks: func(r *ports.MockRepository, tsvc *ports.MockTokenService, pr *ports.MockOAuthProvider) {
				pr.EXPECT().GetUserInfo(context.Background(), "code", "nonce").Return(&domain.User{
					Sub:      "1",
					Provider: "google",
					Name:     "yusuf",
					Email:    "ym@gmail.com",
				}, nil)

				r.EXPECT().FindUserByProviderAndSub("google", "1").Return(&domain.User{}, domain.ErrUserNotFound)

				r.EXPECT().InsertUser(&domain.User{
					Sub:      "1",
					Provider: "google",
					Name:     "yusuf",
					Email:    "ym@gmail.com",
				}).Return(nil)

				tsvc.EXPECT().New(domain.UUID, domain.REFRESH, uint(0)).Return(&domain.Token{
					TokenString: "tokenString",
					UserID:      0,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}, nil)

				tsvc.EXPECT().New(domain.JWT, domain.ACCESS, uint(0)).Return(&domain.Token{
					TokenString: "tokenString",
					UserID:      0,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}, nil)

				tsvc.EXPECT().Save(&domain.Token{
					TokenString: "tokenString",
					UserID:      0,
					Expires:     time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local),
				}).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := ports.NewMockRepository(t)
			provider := ports.NewMockOAuthProvider(t)
			tsvc := ports.NewMockTokenService(t)

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
		})
	}
}
