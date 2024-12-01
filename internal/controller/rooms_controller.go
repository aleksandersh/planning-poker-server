package controller

import (
	"net/http"
	"strings"

	"aleksandersh.github.io/planning-poker-server/internal/rooms"
	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdomain"
	"github.com/gin-gonic/gin"
)

type RoomsController struct {
	authHelper   *AuthHelper
	roomsService *roomsdomain.RoomsService
}

type roomsPostRequest struct {
	Name               string `json:"name"`
	InviteCodeRequired bool   `json:"invite_code_required"`
}

type roomDto struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type roomStateDto struct {
	RoomID      string          `json:"room_id"`
	Name        string          `json:"name"`
	Owner       string          `json:"owner"`
	Commit      string          `json:"commit"`
	Players     []playerDto     `json:"players"`
	CurrentGame *currentGameDto `json:"current_game"`
	GameResults []gameResultDto `json:"game_results"`
}

type playerDto struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type currentGameDto struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	MaxScore        int       `json:"max_score"`
	AverageScore    int       `json:"average_score"`
	IsCardsRevealed bool      `json:"is_card_revealed"`
	Cards           []cardDto `json:"cards"`
}

type cardDto struct {
	Score  int       `json:"score"`
	Player playerDto `json:"player"`
}

type gameResultDto struct {
	GameID       string `json:"game_id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	MaxScore     int    `json:"max_score"`
	AverageScore int    `json:"average_score"`
}

func NewRoomsController(authHelper *AuthHelper, roomsService *roomsdomain.RoomsService) *RoomsController {
	return &RoomsController{authHelper: authHelper, roomsService: roomsService}
}

func (rc *RoomsController) Post(c *gin.Context) {
	user, ok := rc.authHelper.ResolveUser(c)
	if !ok {
		return
	}

	request := roomsPostRequest{Name: "", InviteCodeRequired: false}
	c.ShouldBindJSON(&request)

	room, err := rc.roomsService.Create(user, request.Name, request.InviteCodeRequired)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	response := roomDto{ID: room.ID, Name: room.Name, Owner: room.Owner}
	c.JSON(http.StatusCreated, response)
}

func (rc *RoomsController) Get(c *gin.Context) {
	userID, ok := rc.authHelper.ResolveUserID(c)
	if !ok {
		return
	}

	roomID, ok := requireRoomIDParam(c)
	if !ok {
		return
	}

	room, err := rc.roomsService.Get(userID, roomID)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	response := roomDto{ID: room.ID, Name: room.Name, Owner: room.Owner}
	c.JSON(http.StatusOK, response)
}

func (rc *RoomsController) Delete(c *gin.Context) {
	userID, ok := rc.authHelper.ResolveUserID(c)
	if !ok {
		return
	}

	roomID, ok := requireRoomIDParam(c)
	if !ok {
		return
	}

	if err := rc.roomsService.Delete(userID, roomID); err != nil {
		handleRoomsError(c, err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (rc *RoomsController) Join(c *gin.Context) {
	roomID, ok := requireRoomIDParam(c)
	if !ok {
		return
	}
	user, ok := rc.authHelper.ResolveUser(c)
	if !ok {
		return
	}
	inviteCode := c.Query("invite-code")

	room, err := rc.roomsService.Join(user, roomID, inviteCode)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	response := roomDto{ID: room.ID, Name: room.Name, Owner: room.Owner}
	c.JSON(http.StatusOK, response)
}

func (rc *RoomsController) GetState(c *gin.Context) {
	userID, ok := rc.authHelper.ResolveUserID(c)
	if !ok {
		return
	}

	roomID, ok := requireRoomIDParam(c)
	if !ok {
		return
	}

	commit := c.Query("commit")

	roomState, err := rc.roomsService.GetState(userID, roomID)
	if err != nil {
		handleRoomsError(c, err)
		return
	}

	if roomState.Room.Commit == commit {
		c.AbortWithStatus(http.StatusNotModified)
	}

	players := make([]playerDto, 0, len(roomState.Room.Players))
	for _, player := range roomState.Room.Players {
		players = append(players, mapPlayerToDto(player))
	}
	var currentGame *currentGameDto = nil
	if len(roomState.Games) > 0 {
		game := roomState.Games[len(roomState.Games)-1]
		maxScore := 0
		averageScore := 0
		isCardsRevealed := false
		cards := []cardDto{}
		if game.Status == rooms.GameStatusCompleted {
			maxScore = game.MaxScore
			averageScore = game.AverageScore
			isCardsRevealed = true
			cards = make([]cardDto, 0, len(game.Cards))
			for _, card := range game.Cards {
				cards = append(cards, cardDto{Score: card.Score, Player: mapPlayerToDto(card.Player)})
			}
		}
		currentGame = &currentGameDto{
			ID:              game.ID,
			Name:            game.Name,
			Status:          game.Status,
			MaxScore:        maxScore,
			AverageScore:    averageScore,
			IsCardsRevealed: isCardsRevealed,
			Cards:           cards,
		}
	}
	var results []gameResultDto
	if len(roomState.Games) > 1 {
		results = make([]gameResultDto, 0, len(roomState.Games)-1)
		for _, game := range roomState.Games {
			results = append(results, gameResultDto{
				GameID:       game.ID,
				Name:         game.Name,
				Status:       game.Status,
				MaxScore:     game.MaxScore,
				AverageScore: game.AverageScore,
			})
		}
	} else {
		results = []gameResultDto{}
	}
	response := roomStateDto{
		RoomID:      roomState.Room.ID,
		Name:        roomState.Room.Name,
		Owner:       roomState.Room.Owner,
		Commit:      roomState.Room.Commit,
		Players:     players,
		CurrentGame: currentGame,
		GameResults: results,
	}
	c.JSON(http.StatusOK, response)
}

func requireRoomIDParam(c *gin.Context) (roomID string, ok bool) {
	roomID = c.Param("room_id")
	ok = true
	if len(strings.TrimSpace(roomID)) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		ok = false
	}
	return
}

func mapPlayerToDto(player rooms.Player) playerDto {
	return playerDto{
		ID:    player.UserID,
		Name:  player.Name,
		Color: player.Color,
	}
}
