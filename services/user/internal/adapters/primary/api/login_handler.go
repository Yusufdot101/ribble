package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/secondary/provider/local"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
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
	if err != nil {
		status := http.StatusInternalServerError
		error := err.Error()
		switch {
		case errors.Is(err, domain.ErrUnverifiedAccount):
			status = http.StatusOK
			c.JSON(status, response.Response[any]{
				Message: err.Error(),
			})
			return
		case errors.Is(err, domain.ErrInvalidProviderInputs):
			status = http.StatusBadRequest
			error = "invalid request"
		}
		c.JSON(status, response.Response[any]{
			Error: error,
		})
		return
	}
	panic("shouldnt reach this place")
}

func (h *handler) verify(c *gin.Context) {
	token := c.Query("token")
	refreshToken, accessToken, err := h.svc.ActivateAccount(token)
	if err != nil {
		status := http.StatusInternalServerError
		error := err.Error()
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
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid request",
		})
		return
	}

	ctx := context.Background()
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
