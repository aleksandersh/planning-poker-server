package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"aleksandersh.github.io/planning-poker-server/internal/users/usersdomain"
	"aleksandersh.github.io/planning-poker-server/internal/users/usersdomain/usersmodel"
	"github.com/gin-gonic/gin"
)

var (
	ErrMissingAccessToken = errors.New("missing access token")
)

type AuthHelper struct {
	usersService *usersdomain.Service
}

func NewAuthHelper(usersService *usersdomain.Service) *AuthHelper {
	return &AuthHelper{usersService: usersService}
}

func (h *AuthHelper) ResolveUser(c *gin.Context) (usersmodel.User, bool) {
	accessToken, err := h.getAccessToken(c)
	if err != nil {
		log.Println(fmt.Errorf("authorization failed: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return usersmodel.User{}, false
	}

	user, err := h.usersService.ResolveUserByAccessToken(accessToken)
	if err != nil {
		log.Println(fmt.Errorf("authorization failed: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return usersmodel.User{}, false
	}

	return user, true
}

func (h *AuthHelper) ResolveUserID(c *gin.Context) (string, bool) {
	user, ok := h.ResolveUser(c)
	return user.ID, ok
}

func (h *AuthHelper) getAccessToken(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if len(header) == 0 {
		return "", ErrMissingAccessToken
	}
	token, found := strings.CutPrefix(header, "Bearer ")
	if !found {
		return "", ErrMissingAccessToken
	}
	return token, nil
}
