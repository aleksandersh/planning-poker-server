package usersdata

import (
	"errors"
	"sync"

	"aleksandersh.github.io/planning-poker-server/internal/users"
	"aleksandersh.github.io/planning-poker-server/internal/utils/idutils"
)

var (
	ErrAccessTokenNotFound = errors.New("access token not found")
)

type Repository struct {
	mutex        sync.RWMutex
	users        map[string]users.User
	accessTokens map[string]string
}

func NewRepo() *Repository {
	return &Repository{
		users:        make(map[string]users.User),
		accessTokens: make(map[string]string),
	}
}

func (r *Repository) CreateUser(user users.User) users.User {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := idutils.GenerateID()
	for r.isUserExists(id) {
		id = idutils.GenerateID()
	}

	user.ID = id
	r.users[user.ID] = user
	return user
}

func (r *Repository) CreateToken(userID string) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	accessToken := idutils.GenerateID()
	for r.isAccessTokenExists(accessToken) {
		accessToken = idutils.GenerateID()
	}

	r.accessTokens[accessToken] = userID
	return accessToken
}

func (r *Repository) ResolveUserByAccessToken(accessToken string) (users.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userID, contains := r.accessTokens[accessToken]
	if !contains {
		return users.User{}, ErrAccessTokenNotFound
	}

	user, contains := r.users[userID]
	if !contains {
		return users.User{}, ErrAccessTokenNotFound
	}

	return user, nil
}

func (r *Repository) isUserExists(id string) bool {
	_, contains := r.users[id]
	return contains
}

func (r *Repository) isAccessTokenExists(accessToken string) bool {
	_, contains := r.accessTokens[accessToken]
	return contains
}
