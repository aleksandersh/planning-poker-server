package roomsmodel

import (
	"errors"
)

const (
	GameStatusActive    = "active"
	GameStatusCompleted = "completed"
)

var (
	ErrGameNotFound      = errors.New("game not found")
	ErrIllegalGameStatus = errors.New("illegal game status")
)

type Game struct {
	ID           string
	RoomID       string
	Name         string
	Status       string
	MaxScore     int
	AverageScore int
	Cards        []Card
}

type Card struct {
	Player Player
	Score  int
}
