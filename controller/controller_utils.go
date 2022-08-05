package controller

import (
	"errors"
	"net/http"

	"aleksandersh.github.io/planning-poker-server/data"
	"github.com/gin-gonic/gin"
)

var storage data.RoomsRepository

func getRoomID(c *gin.Context) string {
	return c.Param("room")
}

func getSessionID(c *gin.Context) string {
	return c.GetHeader("Authorization")
}

func handleDataError(c *gin.Context, err *error) bool {
	if err == nil {
		return false
	}
	objErr := *err
	if errors.Is(objErr, data.ErrUnauthorized) {
		c.AbortWithError(http.StatusUnauthorized, objErr)
	} else if errors.Is(objErr, data.ErrMissingResource) {
		c.AbortWithError(http.StatusNotFound, objErr)
	} else if errors.Is(objErr, data.ErrLimitExceeded) {
		c.AbortWithError(http.StatusTooManyRequests, objErr)
	} else if errors.Is(objErr, data.ErrIllegalState) {
		c.AbortWithError(http.StatusConflict, objErr)
	} else {
		c.AbortWithError(http.StatusNotImplemented, objErr)
	}
	return true
}
