package roomsdomain

import (
	"aleksandersh.github.io/planning-poker-server/internal/activity/activitydata"
	"aleksandersh.github.io/planning-poker-server/internal/rooms"
	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdata"
)

type GamesService struct {
	roomsRepository    *roomsdata.Repository
	activityRepository *activitydata.Repository
}

func NewGamesService(roomsRepository *roomsdata.Repository, activityRepository *activitydata.Repository) *GamesService {
	return &GamesService{roomsRepository: roomsRepository, activityRepository: activityRepository}
}

func (s *GamesService) Create(userID string, roomID string, name string) (rooms.Game, error) {
	game := rooms.Game{
		RoomID:       roomID,
		Name:         name,
		Status:       rooms.GameStatusActive,
		MaxScore:     0,
		AverageScore: 0,
		Cards:        []rooms.Card{},
	}
	game, err := s.roomsRepository.AddGame(userID, roomID, game)
	if err == nil {
		s.activityRepository.AddPlayerActivity(game.RoomID, userID)
	}
	return game, err
}

func (s *GamesService) Complete(userID string, gameID string) (rooms.Game, error) {
	game, err := s.roomsRepository.CompleteGame(userID, gameID)
	if err == nil {
		s.activityRepository.AddPlayerActivity(game.RoomID, userID)
	}
	return game, err
}

func (s *GamesService) Reset(userID string, gameID string) (rooms.Game, error) {
	game, err := s.roomsRepository.ResetGame(userID, gameID)
	if err == nil {
		s.activityRepository.AddPlayerActivity(game.RoomID, userID)
	}
	return game, err
}

func (s *GamesService) SendCard(userID string, gameID string, score int) (rooms.Game, error) {
	game, err := s.roomsRepository.SendCard(userID, gameID, score)
	if err == nil {
		s.activityRepository.AddPlayerActivity(game.RoomID, userID)
	}
	return game, err
}

func (s *GamesService) DropCard(userID string, gameID string) (rooms.Game, error) {
	game, err := s.roomsRepository.DropCard(userID, gameID)
	if err == nil {
		s.activityRepository.AddPlayerActivity(game.RoomID, userID)
	}
	return game, err
}
