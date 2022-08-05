package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type gamePostRequest struct {
	Name *string `json:"name"`
}

func PostGame(c *gin.Context) {
	sessionID := getSessionID(c)
	roomID := getRoomID(c)
	request := gamePostRequest{Name: nil}
	c.ShouldBindJSON(&request)

	err := storage.AddGame(roomID, sessionID, request.Name)
	if !handleDataError(c, err) {
		c.AbortWithStatus(http.StatusCreated)
	}
}
