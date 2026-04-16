package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"rygell-dashboard/internal/config"
	"rygell-dashboard/internal/models"
	"rygell-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles login and current-user requests.
type AuthHandler struct {
	userService *services.UserService
	cfg         *config.Config
}

func NewAuthHandler(userService *services.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{userService: userService, cfg: cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.userService.Login(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setAuthCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.clearAuthCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"username": user.Username,
	})
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	usernameAny, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	username, ok := usernameAny.(string)
	if !ok || username != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
		return
	}

	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(input.Name, input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"username": user.Username,
	})
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	usernameAny, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	username, ok := usernameAny.(string)
	if !ok || username != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
		return
	}

	users, err := h.userService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	usernameAny, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	username, ok := usernameAny.(string)
	if !ok || username != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
		return
	}

	idParam := c.Param("id")
	userID64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	targetUserID := uint(userID64)

	currentUserAny, ok := c.Get("user_id")
	if ok {
		if currentUserID, valid := currentUserAny.(uint); valid && currentUserID == targetUserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete currently logged-in admin user"})
			return
		}
	}

	if err := h.userService.DeleteUser(targetUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *AuthHandler) UpdateUserPassword(c *gin.Context) {
	usernameAny, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	username, ok := usernameAny.(string)
	if !ok || username != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
		return
	}

	idParam := c.Param("id")
	userID64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var input models.UpdateUserPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.UpdateUserPassword(uint(userID64), input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
	if h.cfg == nil {
		return
	}
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(strings.TrimSpace(h.cfg.CookieSameSite)) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}
	secure := h.cfg.CookieSecure || sameSite == http.SameSiteNoneMode
	maxAge := 60 * 60 * 24
	if h.cfg.JWTExpiryHours != "" {
		if hours, err := strconv.Atoi(h.cfg.JWTExpiryHours); err == nil && hours > 0 {
			maxAge = hours * 60 * 60
		}
	}
	c.SetSameSite(sameSite)
	c.SetCookie(h.cfg.CookieName, token, maxAge, "/", h.cfg.CookieDomain, secure, true)
}

func (h *AuthHandler) clearAuthCookie(c *gin.Context) {
	if h.cfg == nil {
		return
	}
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(strings.TrimSpace(h.cfg.CookieSameSite)) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}
	secure := h.cfg.CookieSecure || sameSite == http.SameSiteNoneMode
	c.SetSameSite(sameSite)
	c.SetCookie(h.cfg.CookieName, "", -1, "/", h.cfg.CookieDomain, secure, true)
}
