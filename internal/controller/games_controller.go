package controller

import (
	"net/http"

	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdomain"
	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdomain/roomsmodel"
	"github.com/gin-gonic/gin"
)

type GamesController struct {
	authHelper   *AuthHelper
	gamesService *roomsdomain.GamesService
}

type gamePostRequest struct {
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
}

type cardPostRequest struct {
	Score int `json:"score"`
}

type gameDto struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func NewGamesController(authHelper *AuthHelper, gamesService *roomsdomain.GamesService) *GamesController {
	return &GamesController{authHelper: authHelper, gamesService: gamesService}
}

func (gc *GamesController) Post(c *gin.Context) {
	userID, ok := gc.authHelper.ResolveUserID(c)
	if !ok {
		return
	}

	request := gamePostRequest{Name: ""}
	c.ShouldBindJSON(&request)

	game, err := gc.gamesService.Create(userID, request.RoomID, request.Name)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createGameDto(game))
}

func (gc *GamesController) Complete(c *gin.Context) {
	userID, ok := gc.authHelper.ResolveUserID(c)
	if !ok {
		return
	}

	gameID := c.Param("game_id")
	game, err := gc.gamesService.Complete(userID, gameID)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	c.JSON(http.StatusOK, createGameDto(game))
}

func (gc *GamesController) Reset(c *gin.Context) {
	userID, ok := gc.authHelper.ResolveUserID(c)
	if !ok {
		return
	}

	gameID := c.Param("game_id")
	game, err := gc.gamesService.Reset(userID, gameID)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	c.JSON(http.StatusOK, createGameDto(game))
}

func (gc *GamesController) SendCard(c *gin.Context) {
	userID, ok := gc.authHelper.ResolveUserID(c)
	if !ok {
		return
	}

	gameID := c.Param("game_id")

	request := cardPostRequest{Score: 0}
	c.ShouldBindJSON(&request)

	game, err := gc.gamesService.SendCard(userID, gameID, request.Score)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	c.JSON(http.StatusOK, createGameDto(game))
}

func (gc *GamesController) DropCard(c *gin.Context) {
	userID, ok := gc.authHelper.ResolveUserID(c)
	if !ok {
		return
	}

	gameID := c.Param("game_id")

	game, err := gc.gamesService.DropCard(userID, gameID)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	c.JSON(http.StatusOK, createGameDto(game))
}

func createGameDto(game roomsmodel.Game) gameDto {
	return gameDto{
		ID:     game.ID,
		RoomID: game.RoomID,
		Name:   game.Name,
		Status: game.Status,
	}
}
