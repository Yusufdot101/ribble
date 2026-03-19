package api

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yusufdot101/ribble/services/user/internal/adapters/secondary/postgresql"
	"github.com/Yusufdot101/ribble/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ribble/services/user/internal/application/core/services"
	"github.com/Yusufdot101/ribble/services/user/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestGoogleCallbackHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// create postgres test container
	ctx := context.Background()
	ctr, err := postgres.Run(
		ctx,
		"postgres:18.3-alpine",
		postgres.WithDatabase("ripple_user_service_test"),
		postgres.WithUsername("user_service"),
		postgres.WithPassword("verystrongpassword"),
		postgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, ctr)
	if err != nil {
		log.Fatalf("failed to start postgresql: %v", err)
	}

	databaseURL, err := ctr.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to obtain connection string: %v", err)
	}

	repo, err := postgresql.NewAdapter(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	tsvc := services.NewTokenService(repo)

	// mock the google oidc
	provider := ports.NewMockOAuthProvider(t)
	provider.EXPECT().GetUserInfo(mock.Anything, "fake-code", "fake-nonce").Return(&domain.User{
		Name:     "yusuf",
		Email:    "example@gmail.com",
		Provider: "google",
		Sub:      "1",
	}, nil)
	svc := services.NewAuthService(repo, provider, tsvc)
	h := NewHandler(svc)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=fake-code&state=fake-state", nil)
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	req.AddCookie(&http.Cookie{
		Name:  "state",
		Value: "fake-state",
	})
	req.AddCookie(&http.Cookie{
		Name:  "nonce",
		Value: "fake-nonce",
	})

	h.googleCallback(c)

	cookies := resp.Result().Cookies()
	log.Println("cookies: ", cookies)
	for _, c := range cookies {
		if c.Name == "refreshToken" {
			return
		}
	}
	t.Error("expected refreshToken cookie to be set")
}
