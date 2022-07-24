package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"aleksandersh.github.io/planning-poker-server/data"
	"aleksandersh.github.io/planning-poker-server/domain/entity"

	"github.com/gin-gonic/gin"
)

const (
	roomActivityCheckInterval = 10 * time.Second
)

var storage data.RoomsRepository

func main() {
	port := os.Getenv("POKER_PORT")

	if port == "" {
		log.Fatal("variable $POKER_PORT must be set")
	}

	router := gin.Default()

	router.POST("/v1/rooms", postRoom)
	router.GET("/v1/rooms/:room", getRoom)
	router.DELETE("/v1/rooms/:room", deleteRoom)
	router.POST("/v1/rooms/:room/games", postGame)
	router.PATCH("/v1/rooms/:room/games", patchGame)
	router.POST("/v1/rooms/:room/players", postPlayer)
	router.PUT("/v1/rooms/:room/cards", putCard)

	router.Run(":" + port)
}

type roomPostResponse struct {
	RoomID string `json:"room_id"`
}

func postRoom(c *gin.Context) {
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

func deleteRoom(c *gin.Context) {
	roomID := getRoomID(c)
	sessionID := getSessionID(c)
	err := storage.DeleteRoomByID(roomID, sessionID)
	if !handleDataError(c, err) {
		c.AbortWithStatus(http.StatusOK)
	}
}

type gamePatchRequest struct {
	Name     *string `json:"name"`
	Complete bool    `json:"complete"`
	Reset    bool    `json:"reset"`
}

func patchGame(c *gin.Context) {
	roomID := getRoomID(c)
	sessionID := getSessionID(c)
	request := gamePatchRequest{Name: nil, Complete: false, Reset: false}
	c.ShouldBindJSON(&request)

	err := storage.UpdateCurrentGame(roomID, sessionID, request.Name, request.Complete, request.Reset)
	if !handleDataError(c, err) {
		c.AbortWithStatus(http.StatusOK)
	}
}

type gamePostRequest struct {
	Name *string `json:"name"`
}

func postGame(c *gin.Context) {
	sessionID := getSessionID(c)
	roomID := getRoomID(c)
	request := gamePostRequest{Name: nil}
	c.ShouldBindJSON(&request)

	err := storage.AddGame(roomID, sessionID, request.Name)
	if !handleDataError(c, err) {
		c.AbortWithStatus(http.StatusCreated)
	}
}

type playerPostRequest struct {
	Name string `json:"name"`
}

type playerPostResponse struct {
	AccessToken string `json:"access_token"`
}

func postPlayer(c *gin.Context) {
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

type roomDTO struct {
	PlayerID    string         `json:"player_id"`
	Owner       bool           `json:"owner"`
	Commit      string         `json:"commit"`
	CurrentGame currentGameDTO `json:"current_game"`
	Players     []playerDTO    `json:"players"`
	Games       []gameDTO      `json:"games"`
}

type currentGameDTO struct {
	Name              *string   `json:"name"`
	Status            string    `json:"status"`
	SuggestedEstimate int       `json:"suggested_estimate"`
	AverageEstimate   int       `json:"average_estimate"`
	IsCardsRevealed   bool      `json:"is_cards_revealed"`
	Cards             []cardDTO `json:"cards"`
}

type playerDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type gameDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type cardDTO struct {
	Player playerDTO `json:"player"`
	Score  int       `json:"score"`
}

func getRoom(c *gin.Context) {
	roomID := getRoomID(c)
	sessionID := getSessionID(c)
	room, player, owner, err := storage.GetRoomState(roomID, sessionID)
	if handleDataError(c, err) {
		return
	}

	commit := c.Query("commit")
	if room.Commit == commit {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	players := convertPlayersToDTO(room)
	games := []gameDTO{}
	lastGameIndex := len(room.Games) - 1
	if lastGameIndex < 0 {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var currentGame currentGameDTO
	for gameIdx, game := range room.Games {
		if gameIdx == lastGameIndex {
			currentGame = convertGameToDto(&game)
		} else {
			gameDTO := gameDTO{
				ID:    game.ID,
				Name:  game.Name,
				Score: game.SuggestedEstimate,
			}
			games = append(games, gameDTO)
		}
	}
	roomDTO := roomDTO{
		CurrentGame: currentGame,
		PlayerID:    player.ID,
		Owner:       owner,
		Commit:      room.Commit,
		Players:     players,
		Games:       games,
	}
	c.JSON(http.StatusOK, roomDTO)
}

func convertPlayersToDTO(room *entity.Room) []playerDTO {
	players := []playerDTO{}
	for _, session := range room.Sessions {
		playerDTO := playerDTO{
			ID:    session.Player.ID,
			Name:  session.Player.Name,
			Color: session.Player.Color,
		}
		players = append(players, playerDTO)
	}
	return players
}

func convertGameToDto(game *entity.Game) currentGameDTO {
	var suggestedEstimate = 0
	var averageEstimate = 0
	isCardsRevealed := game.Status == entity.GameStatusCompleted
	if isCardsRevealed {
		suggestedEstimate = game.SuggestedEstimate
		averageEstimate = game.AverageEstimate
	}
	cards := []cardDTO{}
	for _, card := range game.Cards {
		playerDTO := playerDTO{
			ID:    card.Player.ID,
			Name:  card.Player.Name,
			Color: card.Player.Color,
		}
		cardDTO := cardDTO{
			Player: playerDTO,
			Score:  card.Score,
		}
		cards = append(cards, cardDTO)
	}
	return currentGameDTO{
		Name:              &game.Name,
		Status:            game.Status,
		SuggestedEstimate: suggestedEstimate,
		AverageEstimate:   averageEstimate,
		IsCardsRevealed:   isCardsRevealed,
		Cards:             cards,
	}
}

type postCardRequest struct {
	Score int `json:"score"`
}

func putCard(c *gin.Context) {
	roomID := getRoomID(c)
	sessionID := getSessionID(c)
	request := postCardRequest{Score: 0}
	c.ShouldBindJSON(&request)

	err := storage.AddCard(roomID, sessionID, request.Score)
	if !handleDataError(c, err) {
		c.AbortWithStatus(http.StatusOK)
	}
}

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
