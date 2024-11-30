package activitydata

import (
	"sync"
	"time"
)

type Repository struct {
	mutex   sync.Mutex
	rooms   map[string]time.Time
	users   map[string]time.Time
	players map[playerKey]time.Time
}

type playerKey struct {
	RoomID string
	UserID string
}

func NewRepository() *Repository {
	return &Repository{
		rooms:   make(map[string]time.Time),
		users:   make(map[string]time.Time),
		players: make(map[playerKey]time.Time),
	}
}

func (r *Repository) AddUserActivity(userID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	now := time.Now()
	r.users[userID] = now
}

func (r *Repository) AddPlayerActivity(roomID string, userID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	r.rooms[roomID] = now
	r.users[userID] = now
	pk := playerKey{RoomID: roomID, UserID: userID}
	r.players[pk] = now
}

func (r *Repository) DeleteActivity(userID string, roomID string) {
}
