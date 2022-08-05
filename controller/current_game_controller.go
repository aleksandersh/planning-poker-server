package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type currentGamePatchRequest struct {
	Name     *string `json:"name"`
	Complete bool    `json:"complete"`
	Reset    bool    `json:"reset"`
}

func PatchCurrentGame(c *gin.Context) {
	roomID := getRoomID(c)
	sessionID := getSessionID(c)
	request := currentGamePatchRequest{Name: nil, Complete: false, Reset: false}
	c.ShouldBindJSON(&request)

	err := storage.UpdateCurrentGame(roomID, sessionID, request.Name, request.Complete, request.Reset)
	if !handleDataError(c, err) {
		c.AbortWithStatus(http.StatusOK)
	}
}
