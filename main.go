package main

import (
	"log"
	"os"

	"aleksandersh.github.io/planning-poker-server/controller"

	"github.com/gin-gonic/gin"
)

const (
	envPort      = "POKER_PORT"
	envMode      = "POKER_MODE"
	envModeDebug = "debug"
)

func main() {
	port := os.Getenv(envPort)
	if port == "" {
		log.Fatal("variable $" + envPort + " must be set")
	}

	isDebug := os.Getenv(envMode) == envModeDebug
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	log.Printf("Start poker app (port=%s, isDebug=%t)", port, isDebug)

	router := gin.Default()

	router.POST("/v1/rooms", controller.PostRoom)
	router.GET("/v1/rooms/:room", controller.GetRoom)
	router.DELETE("/v1/rooms/:room", controller.DeleteRoom)
	router.POST("/v1/rooms/:room/players", controller.PostPlayer)
	router.POST("/v1/rooms/:room/games", controller.PostGame)
	router.PATCH("/v1/rooms/:room/currentgame", controller.PatchCurrentGame)
	router.PUT("/v1/rooms/:room/currentgame/cards", controller.PutCard)

	router.Run(":" + port)
}
