package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/services"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := ports.NewMockRepository(t)
	userID := 1
	repo.EXPECT().GetTokenByStringAndUse("refreshToken", domain.REFRESH).Return(&domain.Token{
		UserID: uint(userID),
	}, nil)

	tsvc := services.NewTokenService(repo)
	asvc := ports.NewMockAuthService(t)
	h := NewHandler(asvc, tsvc, nil)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/refreshtoken", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refreshToken",
		Value: "refreshToken",
	})
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	h.RefreshToken(c)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "accessToken")
}
