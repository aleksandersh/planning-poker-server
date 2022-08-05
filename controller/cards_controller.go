package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type postCardRequest struct {
	Score int `json:"score"`
}

func PutCard(c *gin.Context) {
	roomID := getRoomID(c)
	sessionID := getSessionID(c)
	request := postCardRequest{Score: 0}
	c.ShouldBindJSON(&request)

	err := storage.AddCard(roomID, sessionID, request.Score)
	if !handleDataError(c, err) {
		c.AbortWithStatus(http.StatusOK)
	}
}
