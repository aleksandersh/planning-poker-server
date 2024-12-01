package roomsdata

import (
	"log"
	"strconv"
	"sync"

	"aleksandersh.github.io/planning-poker-server/internal/rooms"
	"aleksandersh.github.io/planning-poker-server/internal/users"
	"aleksandersh.github.io/planning-poker-server/internal/utils/idutils"
)

const (
	roomsLimit    = 300
	playersLimit  = 300
	sessionsLimit = 60
	gamesLimit    = 200
)

var playerColors = []string{
	"FF8B8B",
	"76FFCE",
	"BB86FF",
	"85E2FF",
	"FF86F3",
	"FFCF86",
	"CAFF86",
}

type Repository struct {
	mutex sync.RWMutex
	rooms map[string]rooms.Room
	games map[string]rooms.Game
}

func NewRepo() *Repository {
	return &Repository{
		rooms: make(map[string]rooms.Room),
		games: make(map[string]rooms.Game),
	}
}

func (r *Repository) Create(user users.User, name string, inviteCodeRequired bool) (rooms.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.rooms) >= roomsLimit {
		return rooms.Room{}, rooms.ErrLimitExceeded
	}

	id := r.createRoomID()
	room := rooms.Room{
		ID:                 id,
		Commit:             idutils.GenerateID(),
		Name:               name,
		InviteCodeRequired: inviteCodeRequired,
		Owner:              user.ID,
		Players:            []rooms.Player{newPlayer(user, 0)},
		Games:              []string{},
		VisitorsCount:      1,
	}
	r.rooms[id] = room
	return room, nil
}

func (r *Repository) Get(userID string, roomID string) (rooms.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	room, contains := r.rooms[roomID]
	if !contains {
		return rooms.Room{}, rooms.ErrRoomNotFound
	}
	return room, nil
}

func (r *Repository) Delete(userID string, roomID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, err := r.getOwnedRoom(userID, roomID)
	if err != nil {
		return err
	}

	delete(r.rooms, roomID)
	return nil
}

func (r *Repository) Join(user users.User, roomID string, inviteCode string) (rooms.Room, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, contains := r.rooms[roomID]
	if !contains {
		return rooms.Room{}, rooms.ErrRoomNotFound
	}
	if isPlayerExists(room, user.ID) || !isInviteCodeAccepted(room, inviteCode) {
		return rooms.Room{}, rooms.ErrForbidden
	}

	room.Players = append(room.Players, newPlayer(user, room.VisitorsCount))
	room.VisitorsCount = room.VisitorsCount + 1
	r.saveRoom(room)

	return room, nil
}

func (r *Repository) AddGame(userID string, roomID string, game rooms.Game) (rooms.Game, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, err := r.getOwnedRoom(userID, roomID)
	if err != nil {
		return rooms.Game{}, err
	}

	if len(room.Games) >= gamesLimit {
		return rooms.Game{}, rooms.ErrLimitExceeded
	}

	game.ID = r.createGameID()
	game.RoomID = roomID
	if len(game.Name) == 0 {
		game.Name = "Game " + strconv.Itoa(len(room.Games)+1)
	}

	r.games[game.ID] = game

	room.Games = append(room.Games, game.ID)
	r.saveRoom(room)

	return game, nil
}

func (r *Repository) CompleteGame(userID string, gameID string) (rooms.Game, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, game, err := r.getRoomAndGame(userID, gameID)
	if err != nil {
		return game, err
	}

	if game.Status == rooms.GameStatusCompleted {
		return game, nil
	}

	game = estimateGame(game)
	game.Status = rooms.GameStatusCompleted
	r.games[game.ID] = game

	r.saveRoom(room)

	return game, nil
}

func (r *Repository) ResetGame(userID string, gameID string) (rooms.Game, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, game, err := r.getRoomAndGame(userID, gameID)
	if err != nil {
		return game, err
	}

	game = estimateGame(game)
	game.Status = rooms.GameStatusActive
	game.MaxScore = 0
	game.AverageScore = 0
	game.Cards = []rooms.Card{}
	r.games[game.ID] = game

	r.saveRoom(room)

	return game, nil
}

func (r *Repository) SendCard(userID string, gameID string, score int) (rooms.Game, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, game, err := r.getRoomAndGame(userID, gameID)
	if err != nil {
		return game, err
	}

	if game.Status == rooms.GameStatusCompleted {
		return game, rooms.ErrIllegalGameStatus
	}

	player, err := getPlayer(room, userID)
	if err != nil {
		return game, err
	}

	game.Cards = putUserCard(game.Cards, rooms.Card{Player: player, Score: score})
	r.games[game.ID] = game

	r.saveRoom(room)

	return game, nil
}

func (r *Repository) DropCard(userID string, gameID string) (rooms.Game, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, game, err := r.getRoomAndGame(userID, gameID)
	if err != nil {
		return game, err
	}

	if game.Status == rooms.GameStatusCompleted {
		return game, rooms.ErrIllegalGameStatus
	}

	game.Cards = dropUserCard(game.Cards, userID)
	r.games[game.ID] = game

	r.saveRoom(room)

	return game, nil
}

func (r *Repository) GetRoomState(userID string, roomID string) (rooms.RoomState, error) {
	room, contains := r.rooms[roomID]
	if !contains {
		return rooms.RoomState{}, rooms.ErrRoomNotFound
	}
	if !isPlayerExists(room, userID) {
		return rooms.RoomState{}, rooms.ErrForbidden
	}

	games := make([]rooms.Game, 0, len(room.Games))
	for _, gameID := range room.Games {
		game, contains := r.games[gameID]
		if contains {
			games = append(games, game)
		}
	}

	return rooms.RoomState{Room: room, Games: games}, nil
}

func (r *Repository) createRoomID() string {
	counter := 0
	id := generateRoomID()
	for r.isRoomExists(id) {
		counter++
		if counter == 10_000 {
			log.Panicf("too many attempts to generate next room ID")
		}
		id = generateRoomID()
	}
	return id
}

func (r *Repository) isRoomExists(id string) bool {
	_, contains := r.rooms[id]
	return contains
}

func (r *Repository) createGameID() string {
	counter := 0
	id := idutils.GenerateID()
	for r.isGameExists(id) {
		counter++
		if counter == 10_000 {
			log.Panicf("too many attempts to generate next game ID")
		}
		id = idutils.GenerateID()
	}
	return id
}

func (r *Repository) isGameExists(id string) bool {
	_, contains := r.games[id]
	return contains
}

func (r *Repository) getOwnedRoom(userID string, roomID string) (rooms.Room, error) {
	room, contains := r.rooms[roomID]
	if !contains {
		return rooms.Room{}, rooms.ErrRoomNotFound
	}
	if room.Owner != userID {
		return rooms.Room{}, rooms.ErrForbidden
	}
	return room, nil
}

func newPlayer(user users.User, index int) rooms.Player {
	name := user.Name
	if len(name) == 0 {
		name = "Player " + strconv.Itoa(index+1)
	}

	color := playerColors[index%len(playerColors)]
	return rooms.Player{UserID: user.ID, Name: name, Color: color}
}

func (r *Repository) getRoomAndGame(userID string, gameID string) (rooms.Room, rooms.Game, error) {
	game, contains := r.games[gameID]
	if !contains {
		return rooms.Room{}, rooms.Game{}, rooms.ErrGameNotFound
	}
	room, contains := r.rooms[game.RoomID]
	if !contains {
		return rooms.Room{}, rooms.Game{}, rooms.ErrGameNotFound
	}
	if room.Owner != userID {
		return rooms.Room{}, rooms.Game{}, rooms.ErrForbidden
	}
	return room, game, nil
}

func (r *Repository) saveRoom(room rooms.Room) {
	room.Commit = idutils.GenerateID()
	r.rooms[room.ID] = room
}
