package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/secondary/provider/local"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8,max=72"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid request",
		})
		return
	}

	ctx := context.Background()
	_, _, err := h.svc.HandleAuth(ctx, map[string]string{
		"method":   "signup",
		"name":     req.Name,
		"email":    req.Email,
		"password": req.Password,
	}, local.LocalProviderName)
	if err == nil || errors.Is(err, domain.ErrUnverifiedAccount) {
		c.JSON(http.StatusCreated, response.Response[any]{
			Message: err.Error(),
		})
		return
	}
	status := http.StatusInternalServerError
	error := err.Error()
	switch {
	case errors.Is(err, domain.ErrInvalidProviderInputs):
		status = http.StatusBadRequest
		error = "invalid request"
	}
	c.JSON(status, response.Response[any]{
		Error: error,
	})
}

func (h *handler) verify(c *gin.Context) {
	token := c.Query("token")
	identityID, err := strconv.ParseUint(c.Query("identity"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid identity id",
		})
		return
	}
	refreshToken, accessToken, err := h.svc.ActivateAccount(token, uint(identityID))
	if err != nil {
		status := http.StatusInternalServerError
		error := err.Error()
		switch {
		case errors.Is(err, domain.ErrAccountAlreadyActivated):
			status = http.StatusForbidden
		case errors.Is(err, domain.ErrRecordNotFound):
			status = http.StatusNotFound
		}
		c.JSON(status, response.Response[any]{
			Error: error,
		})
		return
	}

	c.SetCookie("refreshToken", refreshToken, int(config.GetRefreshTokenTTL().Seconds()), "/", "", config.RefreshTokenIsSecure(), true)
	c.JSON(http.StatusCreated, response.Response[string]{
		Data: accessToken,
	})
}

func (h *handler) login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8,max=72"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid request",
		})
		return
	}

	ctx := c.Request.Context()
	refreshToken, accessToken, err := h.svc.HandleAuth(ctx, map[string]string{
		"method":   "login",
		"email":    req.Email,
		"password": req.Password,
	}, local.LocalProviderName)
	if err != nil {
		status := http.StatusInternalServerError
		error := err.Error()
		switch {
		case errors.Is(err, domain.ErrUnverifiedAccount):
			status = http.StatusOK
		case errors.Is(err, domain.ErrInvalidProviderInputs):
			status = http.StatusBadRequest
			error = "invalid request"
		}
		c.JSON(status, response.Response[any]{
			Error: error,
		})
		return
	}

	c.SetCookie("refreshToken", refreshToken, int(config.GetRefreshTokenTTL().Seconds()), "/", "", config.RefreshTokenIsSecure(), true)
	c.JSON(http.StatusCreated, response.Response[string]{
		Data: accessToken,
	})
}
