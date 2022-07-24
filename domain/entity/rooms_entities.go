package entity

import (
	"time"
)

const (
	GameStatusActive    = "active"
	GameStatusCompleted = "completed"
	UnknownIndex        = -1
)

type Room struct {
	ID           string
	Commit       string
	LastActivity time.Time
	Sessions     []Session
	Games        []Game
	PlayersCount int
}

type Game struct {
	ID                string
	Name              string
	Status            string
	SuggestedEstimate int
	AverageEstimate   int
	Cards             []Card
}

type Card struct {
	Player Player
	Score  int
}

type Player struct {
	ID    string
	Name  string
	Color string
}

type Session struct {
	ID           string
	LastActivity time.Time
	Player       Player
}

func (room *Room) GetSessionByID(sessionID string) (idx int, session *Session) {
	for sessionIdx, session := range room.Sessions {
		if session.ID == sessionID {
			return sessionIdx, &session
		}
	}
	return UnknownIndex, nil
}

func (room *Room) GetLastGame() (idx int, game *Game) {
	lastGameIndex := len(room.Games) - 1
	if lastGameIndex >= 0 {
		game := room.Games[lastGameIndex]
		return lastGameIndex, &game
	}
	return UnknownIndex, nil
}

func (room *Room) RemovePlayerResult(playerID string) bool {
	lastGameIndex := len(room.Games) - 1
	if lastGameIndex >= 0 {
		game := room.Games[lastGameIndex]
		if game.Status == GameStatusCompleted {
			return false
		}
		cards := game.Cards
		for cardIdx, card := range cards {
			if card.Player.ID == playerID {
				game.Cards = append(cards[:cardIdx], cards[cardIdx+1:]...)
				game.recalculateEstimate()
				room.Games[lastGameIndex] = game
				return true
			}
		}
	}
	return false
}

func (game *Game) AddCard(player Player, score int) {
	targetCard := Card{Player: player, Score: score}
	for idx, card := range game.Cards {
		if card.Player.ID == player.ID {
			game.Cards[idx] = targetCard
			return
		}
	}
	game.Cards = append(game.Cards, targetCard)
	game.recalculateEstimate()
}

func (game *Game) recalculateEstimate() {
	suggested := 0
	sum := 0
	count := 0
	for _, card := range game.Cards {
		score := card.Score
		if score >= 0 {
			sum = sum + score
			count = count + 1
			if suggested >= 0 && suggested < score {
				suggested = score
			}
		} else {
			if suggested > score {
				suggested = score
			}
		}
	}
	game.SuggestedEstimate = suggested
	if count > 0 {
		average := sum / count
		if sum%count > 0 {
			average = average + 1
		}
		game.AverageEstimate = average
	} else {
		game.AverageEstimate = 0
	}
}

func (game *Game) Reset() {
	game.Cards = []Card{}
	game.AverageEstimate = 0
	game.SuggestedEstimate = 0
	game.Status = GameStatusActive
}

func (game *Game) Complete() {
	game.Status = GameStatusCompleted
}
