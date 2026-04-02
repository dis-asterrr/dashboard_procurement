package handlers

import (
	"net/http"
	"strconv"

	"rygell-dashboard/internal/models"
	"rygell-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles login and current-user requests.
type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
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

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"user": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"username": user.Username,
		},
	})
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
