package handler

import (
	"errors"
	"net/http"

	"alioth-hrc/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register 用户注册（无需 Bearer）。
//
//	@Summary		用户注册
//	@Description	创建账号并自动创建默认账本；成功后返回用户基本信息。
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.RegisterInput	true	"注册信息"
//	@Success		201		{object}	handler.RegisterResponse
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		409		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var in service.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.svc.Register(c.Request.Context(), in)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
	})
}

// Login 登录并获取访问令牌与刷新令牌（无需 Bearer）。
//
//	@Summary		登录
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.LoginInput	true	"登录信息"
//	@Success		200		{object}	service.TokenPair
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var in service.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := h.svc.Login(c.Request.Context(), in)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

// Refresh 使用 refresh_token 换取新的令牌对（无需 Bearer）。
//
//	@Summary		刷新令牌
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.RefreshInput	true	"刷新令牌"
//	@Success		200		{object}	service.TokenPair
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		401		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var in service.RefreshInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := h.svc.Refresh(c.Request.Context(), in)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefresh) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh failed"})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

// Logout 吊销 refresh_token（无需 Bearer）。
//
//	@Summary		登出
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.LogoutInput	true	"刷新令牌"
//	@Success		200		{object}	handler.OKStatus
//	@Failure		400		{object}	handler.ErrorJSON
//	@Failure		500		{object}	handler.ErrorJSON
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var in service.LogoutInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Logout(c.Request.Context(), in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
