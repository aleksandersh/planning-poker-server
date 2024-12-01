package usersdomain

import (
	"aleksandersh.github.io/planning-poker-server/internal/activity/activitydata"
	"aleksandersh.github.io/planning-poker-server/internal/users"
	"aleksandersh.github.io/planning-poker-server/internal/users/usersdata"
)

type Service struct {
	usersRepository    *usersdata.Repository
	activityRepository *activitydata.Repository
}

func NewService(usersRepository *usersdata.Repository, activityRepository *activitydata.Repository) *Service {
	return &Service{usersRepository: usersRepository, activityRepository: activityRepository}
}

func (s *Service) Add(name string) (users.User, string) {
	user := users.User{Name: name}
	user = s.usersRepository.CreateUser(user)
	s.activityRepository.AddUserActivity(user.ID)
	// todo: users activity watcher
	accessToken := s.usersRepository.CreateToken(user.ID)
	return user, accessToken
}

func (s *Service) ResolveUserByAccessToken(accessToken string) (users.User, error) {
	return s.usersRepository.ResolveUserByAccessToken(accessToken)
}
