package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/secondary/provider/google"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"github.com/gin-gonic/gin"
)

type handler struct {
	svc  ports.AuthService
	tsvc ports.TokenService
	usvc ports.UserService
}

func NewHandler(svc ports.AuthService, tsvc ports.TokenService, usvc ports.UserService) *handler {
	return &handler{
		svc:  svc,
		tsvc: tsvc,
		usvc: usvc,
	}
}

func (h *handler) googleBegin(c *gin.Context) {
	// get the authURL, state and nonce cookies
	url, state, nonce, err := h.svc.BeginAuth(google.GoogleProviderName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: "an error occurred, please try again later",
		})
		return
	}

	// set state and nonce cookies to response
	c.SetCookie("state", state, int(5*time.Minute.Seconds()), "/", "", false, true)
	c.SetCookie("nonce", nonce, int(5*time.Minute.Seconds()), "/", "", false, true)

	// redirect user to the authURL
	c.Redirect(http.StatusFound, url)
}

func (h *handler) googleCallback(c *gin.Context) {
	// read cookies, call h.svc.HandleCallback, set your own cookie
	state, err := c.Cookie("state")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if state != c.Query("state") {
		c.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: "state doesnt match",
		},
		)
		return
	}

	nonce, err := c.Cookie("nonce")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	refreshToken, _, err := h.svc.HandleAuth(ctx, map[string]string{
		"code":  c.Query("code"),
		"nonce": nonce,
	}, google.GoogleProviderName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	c.SetCookie("refreshToken", refreshToken, int(config.GetRefreshTokenTTL().Seconds()), "/", "", config.RefreshTokenIsSecure(), true)
	c.SetSameSite(config.GetCookieSameSite())
	c.Redirect(http.StatusFound, config.GetFrontendURL())
}
