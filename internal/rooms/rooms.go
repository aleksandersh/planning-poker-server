package rooms

import (
	"errors"
	"time"
)

const (
	UnknownIndex = -1
)

var (
	ErrForbidden     = errors.New("forbidden")
	ErrRoomNotFound  = errors.New("room not found")
	ErrLimitExceeded = errors.New("resource limit exceeded")
)

type RoomState struct {
	Room  Room
	Games []Game
}

type Room struct {
	ID                 string
	Commit             string
	Name               string
	InviteCodeRequired bool
	Owner              string
	Players            []Player
	InviteCodes        []InviteCode
	Games              []string
	VisitorsCount      int
}

type InviteCode struct {
	Code      string
	CreatedAt time.Time
}
