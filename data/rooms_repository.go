package data

import (
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"aleksandersh.github.io/planning-poker-server/domain/entity"

	"github.com/google/uuid"
)

const (
	roomsLimit          = 20
	sessionsLimit       = 60
	roomInactiveTime    = 60 * time.Minute
	sessionInactiveTime = 10 * time.Minute
)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var playerColors = []string{
	"FF8B8B",
	"76FFCE",
	"BB86FF",
	"85E2FF",
	"FF86F3",
	"FFCF86",
	"CAFF86",
}

type RoomsRepository struct {
	mutex sync.Mutex
	rooms []entity.Room
}

func (storage *RoomsRepository) AddRoom() (string, *error) {
	storage.lock()
	defer storage.unlock()
	if len(storage.rooms) >= roomsLimit {
		return "", &ErrLimitExceeded
	}
	now := time.Now()
	gameName := "Game 1"
	game := entity.Game{
		ID:                generateID(),
		Name:              gameName,
		Status:            entity.GameStatusActive,
		SuggestedEstimate: 0,
		AverageEstimate:   0,
		Cards:             []entity.Card{},
	}
	room := entity.Room{
		ID:           createRoomID(),
		Commit:       generateID(),
		LastActivity: now,
		Sessions:     []entity.Session{},
		Games:        []entity.Game{game},
		PlayersCount: 1,
	}
	storage.rooms = append(storage.rooms, room)
	return room.ID, nil
}

// true if room still active
func (storage *RoomsRepository) CheckRoomActivity(roomID string) bool {
	now := time.Now()
	storage.lock()
	defer storage.unlock()
	idx, room := storage.getRoomByID(roomID)
	if room == nil {
		return false
	}

	if len(room.Sessions) > 0 {
		if removeExpiredSessions(room, now, sessionInactiveTime) {
			storage.updateRoomByIndex(*room, idx)
		}
	} else if room.LastActivity.Add(roomInactiveTime).Before(now) {
		storage.deleteRoomByIndex(idx)
		return false
	}
	return true
}

func removeExpiredSessions(room *entity.Room, now time.Time, inactivityTime time.Duration) bool {
	controlTime := now.Add(-inactivityTime)
	var sessions []entity.Session
	isChanged := false
	for _, session := range room.Sessions {
		if session.LastActivity.After(controlTime) {
			sessions = append(sessions, session)
		} else {
			room.RemovePlayerResult(session.Player.ID)
			isChanged = true
		}
	}
	if isChanged {
		room.Sessions = sessions
	}
	return isChanged
}

func (storage *RoomsRepository) updateRoomByIndex(room entity.Room, idx int) {
	room.Commit = generateID()
	storage.rooms[idx] = room
}

func (storage *RoomsRepository) getRoomByID(roomID string) (idx int, room *entity.Room) {
	for idx, room := range storage.rooms {
		if room.ID == roomID {
			return idx, &room
		}
	}
	return entity.UnknownIndex, nil
}

func (storage *RoomsRepository) DeleteRoomByID(roomID string, sessionID string) *error {
	storage.lock()
	defer storage.unlock()
	idx, _, err := storage.requireRoomForOwner(roomID, sessionID)
	if err != nil {
		return err
	}
	rooms := storage.rooms
	storage.rooms = append(rooms[:idx], rooms[idx+1:]...)
	return nil
}

func (storage *RoomsRepository) deleteRoomByIndex(idx int) {
	rooms := storage.rooms
	storage.rooms = append(rooms[:idx], rooms[idx+1:]...)
}

func (storage *RoomsRepository) GetRoomState(roomID string, sessionID string) (*entity.Room, *entity.Player, bool, *error) {
	storage.lock()
	defer storage.unlock()
	roomIdx, room := storage.getRoomByID(roomID)
	if room == nil {
		return nil, nil, false, &ErrMissingResource
	}
	sessionIdx, session := room.GetSessionByID(sessionID)
	if session == nil {
		return nil, nil, false, &ErrUnauthorized
	}

	now := time.Now()
	session.LastActivity = now
	room.Sessions[sessionIdx] = *session
	room.LastActivity = now
	storage.rooms[roomIdx] = *room
	owner := isOwnerSession(room, sessionID)
	return room, &session.Player, owner, nil
}

func (storage *RoomsRepository) AddPlayer(roomID string, name string) (*entity.Session, *error) {
	storage.lock()
	defer storage.unlock()

	idx, room := storage.getRoomByID(roomID)
	if room == nil {
		return nil, &ErrMissingResource
	}
	if len(room.Sessions) >= sessionsLimit {
		return nil, &ErrLimitExceeded
	}

	playersCount := room.PlayersCount
	if name == "" {
		playerNumber := playersCount + 1
		name = "Player " + strconv.Itoa(playerNumber)
	}

	color := "000000"
	if playersCount < len(playerColors) {
		color = playerColors[playersCount]
	}
	player := entity.Player{ID: generateID(), Name: name, Color: color}
	now := time.Now()
	playerSession := entity.Session{ID: generateID(), LastActivity: now, Player: player}
	room.Sessions = append(room.Sessions, playerSession)
	if playersCount == math.MaxInt {
		room.PlayersCount = room.PlayersCount + 1
	} else {
		room.PlayersCount = 0
	}
	storage.updateRoomByIndex(*room, idx)
	return &playerSession, nil
}

func (storage *RoomsRepository) AddGame(roomID string, sessionID string, name *string) *error {
	storage.lock()
	defer storage.unlock()
	idx, room, err := storage.requireRoomForOwner(roomID, sessionID)
	if err != nil {
		return err
	}
	gameName := ""
	if name != nil {
		gameName = *name
	} else {
		gameNumber := strconv.Itoa(len(room.Games) + 1)
		gameName = "Game " + gameNumber
	}
	game := entity.Game{ID: generateID(), Name: gameName, Status: entity.GameStatusActive, SuggestedEstimate: 0, AverageEstimate: 0, Cards: []entity.Card{}}
	room.Games = append(room.Games, game)
	storage.updateRoomByIndex(*room, idx)

	return nil
}

func (storage *RoomsRepository) AddCard(roomID string, sessionID string, score int) *error {
	storage.lock()
	defer storage.unlock()

	roomIdx, room := storage.getRoomByID(roomID)
	if room == nil {
		return &ErrMissingResource
	}
	_, session := room.GetSessionByID(sessionID)
	if session == nil {
		return &ErrUnauthorized
	}

	gameIdx, game := room.GetLastGame()
	if game == nil {
		return &ErrMissingResource
	}

	if game.Status != entity.GameStatusActive {
		return &ErrIllegalState
	}

	game.AddCard(session.Player, score)
	room.Games[gameIdx] = *game
	storage.updateRoomByIndex(*room, roomIdx)
	return nil
}

func (storage *RoomsRepository) UpdateCurrentGame(roomID string, sessionID string, name *string, complete bool, reset bool) *error {
	storage.lock()
	defer storage.unlock()
	idx, room, err := storage.requireRoomForOwner(roomID, sessionID)
	if err != nil {
		return err
	}
	gameIdx, game := room.GetLastGame()
	if game == nil {
		return &ErrMissingResource
	}

	if reset {
		game.Reset()
	}
	if complete {
		game.Complete()
	}
	if name != nil {
		game.Name = *name
	}
	room.Games[gameIdx] = *game
	storage.updateRoomByIndex(*room, idx)
	return nil
}

func isOwnerSession(room *entity.Room, sessionID string) bool {
	if len(room.Sessions) >= 0 {
		return room.Sessions[0].ID == sessionID
	} else {
		return false
	}
}

func (storage *RoomsRepository) requireRoomForOwner(roomID string, sessionID string) (idx int, room *entity.Room, err *error) {
	idx, room = storage.getRoomByID(roomID)
	if room == nil {
		return entity.UnknownIndex, nil, &ErrMissingResource
	}
	if !isOwnerSession(room, sessionID) {
		return entity.UnknownIndex, nil, &ErrUnauthorized
	}
	return idx, room, nil
}

func createRoomID() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func generateID() string {
	return uuid.New().String()
}

func (storage *RoomsRepository) lock() {
	storage.mutex.Lock()
}

func (storage *RoomsRepository) unlock() {
	storage.mutex.Unlock()
}
