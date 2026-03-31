package services

import (
	"errors"
	"testing"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		isValid bool
	}{
		{
			name:    "vaild email",
			email:   "example@gmail.com",
			isValid: true,
		},
		{
			name:  "invaild email",
			email: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := isValidEmail(test.email)
			if isValid != test.isValid {
				t.Errorf("expected isValidEmail=%v, got isValidEmail=%v", test.isValid, isValid)
			}
		})
	}
}

func TestGetUsersByEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		setupRepo func(r *ports.MockRepository)
		wantErr   bool
		err       error
	}{
		{
			name:  "valid flow",
			email: "example@gmail.com",
			setupRepo: func(r *ports.MockRepository) {
				r.EXPECT().FindUsersByEmail("example@gmail.com").Return([]*domain.User{}, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := ports.NewMockRepository(t)
			test.setupRepo(repo)
			userSvc := NewUserService(repo)
			_, err := userSvc.GetUsersByEmail(test.email)
			if err != nil {
				if !test.wantErr {
					t.Fatalf("unexpected error: %v", err)
				}
				if !errors.Is(err, test.err) {
					t.Fatalf("expected error=%v, got error=%v", test.err, err)
				}
				return
			}
			if test.wantErr {
				t.Fatalf("expected error=%v", test.err)
			}
		})
	}
}
