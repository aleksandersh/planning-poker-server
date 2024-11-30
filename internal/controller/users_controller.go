package controller

import (
	"net/http"
	"strings"

	"aleksandersh.github.io/planning-poker-server/internal/users/usersdomain"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	service *usersdomain.Service
}

type usersRegisterRequest struct {
	Name string `json:"name"`
}

type usersRegisterResponse struct {
	User        userDto `json:"user"`
	AccessToken string  `json:"access_token"`
}

type userDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewUsersController(service *usersdomain.Service) *UsersController {
	return &UsersController{service: service}
}

func (uc *UsersController) Register(c *gin.Context) {
	request := usersRegisterRequest{Name: ""}
	c.ShouldBindJSON(&request)

	if len(strings.TrimSpace(request.Name)) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, accessToken := uc.service.Add(request.Name)

	userDto := userDto{ID: user.ID, Name: user.Name}
	response := usersRegisterResponse{User: userDto, AccessToken: accessToken}
	c.JSON(http.StatusCreated, response)
}
