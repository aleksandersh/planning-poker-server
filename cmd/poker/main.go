package main

import (
	"log"
	"os"

	"aleksandersh.github.io/planning-poker-server/internal/server"
	"github.com/gin-gonic/gin"
)

const (
	envAddress   = "POKER_ADDRESS"
	envMode      = "POKER_MODE"
	envModeDebug = "debug"
)

func main() {
	address := os.Getenv(envAddress)
	if address == "" {
		log.Fatal("variable $" + envAddress + " must be set")
	}

	isDebug := os.Getenv(envMode) == envModeDebug
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	log.Printf("Start poker app (address=%s, isDebug=%t)", address, isDebug)

	server.Start(address)
}
