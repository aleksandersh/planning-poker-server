package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	roomActivityCheckInterval = 10 * time.Second
)

type roomPostResponse struct {
	RoomID string `json:"room_id"`
}

func PostRoom(c *gin.Context) {
	roomID, err := storage.AddRoom()
	if handleDataError(c, err) {
		return
	}
	response := roomPostResponse{RoomID: roomID}
	c.JSON(http.StatusCreated, response)
	go startRoomActivityWatcher(roomID)
}

func startRoomActivityWatcher(roomID string) {
	ticker := time.NewTicker(roomActivityCheckInterval)
	defer ticker.Stop()
	for {
		<-ticker.C
		if !storage.CheckRoomActivity(roomID) {
			break
		}
	}
}

func DeleteRoom(c *gin.Context) {
	roomID := getRoomID(c)
	sessionID := getSessionID(c)
	err := storage.DeleteRoomByID(roomID, sessionID)
	if !handleDataError(c, err) {
		c.AbortWithStatus(http.StatusOK)
	}
}
