package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type playerPostRequest struct {
	Name string `json:"name"`
}

type playerPostResponse struct {
	AccessToken string `json:"access_token"`
}

func PostPlayer(c *gin.Context) {
	roomID := getRoomID(c)
	request := playerPostRequest{Name: ""}
	c.ShouldBindJSON(&request)
	session, err := storage.AddPlayer(roomID, request.Name)
	if handleDataError(c, err) {
		return
	}
	response := playerPostResponse{AccessToken: session.ID}
	c.JSON(http.StatusCreated, response)
}
