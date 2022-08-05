package controller

import (
	"net/http"

	"aleksandersh.github.io/planning-poker-server/domain/entity"
	"github.com/gin-gonic/gin"
)

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

func GetRoom(c *gin.Context) {
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
