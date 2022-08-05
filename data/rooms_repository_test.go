package data

import (
	"testing"
	"time"
)

func TestCheckRoomActivity_activeSession(t *testing.T) {
	var repository RoomsRepository
	roomID, err := repository.AddRoom()
	if err != nil {
		t.Error("Failed to create room")
	}
	createdSession, err := repository.AddPlayer(roomID, "")
	if err != nil {
		t.Error("Failed to create player")
	}
	roomState, _, _, err := repository.GetRoomState(roomID, createdSession.ID)
	if err != nil || roomState == nil {
		t.Error("Failed to retrieve room")
	}

	isRoomActive := repository.CheckRoomActivity(roomID)

	if !isRoomActive {
		t.Error("Inactive room")
	}
	if len(repository.rooms) != 1 {
		t.Error("Should exist only one room")
	}
	room := repository.rooms[0]
	if room.Commit != roomState.Commit {
		t.Error("Illegal room commit")
	}
	if len(room.Sessions) != 1 {
		t.Error("Should exist only one session")
	}
	session := &room.Sessions[0]
	if session.ID != createdSession.ID {
		t.Error("Unexpected session ID")
	}
}

func TestCheckRoomActivity_removeInactiveSession(t *testing.T) {
	var repository RoomsRepository
	roomID, err := repository.AddRoom()
	if err != nil {
		t.Error("Failed to create room")
	}
	createdSession, err := repository.AddPlayer(roomID, "")
	if err != nil || createdSession == nil {
		t.Error("Failed to create player")
	}
	roomState, _, _, err := repository.GetRoomState(roomID, createdSession.ID)
	if err != nil || roomState == nil {
		t.Error("Failed to retrieve room")
	}
	now := time.Now()
	repository.rooms[0].Sessions[0].LastActivity = now.Add(-sessionInactiveTime)

	isRoomActive := repository.CheckRoomActivity(roomID)

	if !isRoomActive {
		t.Error("Inactive room")
	}
	if len(repository.rooms) != 1 {
		t.Error("Should exist only one room")
	}
	room := repository.rooms[0]
	if room.Commit == roomState.Commit {
		t.Error("Illegal room commit")
	}
	if len(room.Sessions) > 0 {
		t.Error("Inactive session survived")
	}
}

func TestCheckRoomActivity_removeInactiveRoom(t *testing.T) {
	var repository RoomsRepository
	roomID, err := repository.AddRoom()
	if err != nil {
		t.Error("Failed to create room")
	}
	now := time.Now()
	repository.rooms[0].LastActivity = now.Add(-roomInactiveTime)

	isRoomActive := repository.CheckRoomActivity(roomID)

	if isRoomActive {
		t.Error("Room is active: missing method result")
	}
	if len(repository.rooms) != 0 {
		t.Error("Room is active: illegal room count")
	}
}
