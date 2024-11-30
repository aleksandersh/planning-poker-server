package roomsdomain

import (
	"aleksandersh.github.io/planning-poker-server/internal/activity/activitydata"
	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdata"
	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdomain/roomsmodel"
	"aleksandersh.github.io/planning-poker-server/internal/users/usersdomain/usersmodel"
)

type RoomsService struct {
	roomsRepository    *roomsdata.Repository
	activityRepository *activitydata.Repository
}

func NewRoomsService(roomsRepository *roomsdata.Repository, activityRepository *activitydata.Repository) *RoomsService {
	return &RoomsService{roomsRepository: roomsRepository, activityRepository: activityRepository}
}

func (rs *RoomsService) Create(user usersmodel.User, name string, inviteCodeRequired bool) (roomsmodel.Room, error) {
	room, err := rs.roomsRepository.Create(user, name, inviteCodeRequired)
	if err != nil {
		return room, err
	} else {
		rs.activityRepository.AddPlayerActivity(room.ID, user.ID)
	}

	// todo: start activity watcher
	return room, err
}

func (rs *RoomsService) Get(userID string, roomID string) (roomsmodel.Room, error) {
	room, err := rs.roomsRepository.Get(userID, roomID)
	if err == nil {
		rs.activityRepository.AddPlayerActivity(roomID, userID)
	}
	return room, err
}

func (rs *RoomsService) Delete(userID string, roomID string) error {
	err := rs.roomsRepository.Delete(userID, roomID)
	if err == nil {
		rs.activityRepository.AddUserActivity(userID)
		// todo: stop room activity watcher
	}
	return err
}

func (rs *RoomsService) Join(user usersmodel.User, roomID string, inviteCode string) (roomsmodel.Room, error) {
	room, err := rs.roomsRepository.Join(user, roomID, inviteCode)
	if err == nil {
		rs.activityRepository.AddPlayerActivity(roomID, user.ID)
	}
	return room, err
}

func (rs *RoomsService) GetState(userID string, roomID string) (roomsmodel.RoomState, error) {
	roomState, err := rs.roomsRepository.GetRoomState(userID, roomID)
	if err == nil {
		rs.activityRepository.AddPlayerActivity(roomID, userID)
	}
	return roomState, err
}
