package controller

import (
	"errors"
	"net/http"

	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdomain/roomsmodel"
	"github.com/gin-gonic/gin"
)

func handleRoomsError(c *gin.Context, err error) {
	if errors.Is(err, roomsmodel.ErrRoomNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
	} else if errors.Is(err, roomsmodel.ErrForbidden) {
		c.AbortWithStatus(http.StatusForbidden)
	} else if errors.Is(err, roomsmodel.ErrLimitExceeded) {
		c.AbortWithStatus(http.StatusTooManyRequests)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
