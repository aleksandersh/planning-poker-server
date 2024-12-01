package roomsdata

import (
	"math/rand"

	"aleksandersh.github.io/planning-poker-server/internal/rooms"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func generateRoomID() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func isPlayerExists(room rooms.Room, userID string) bool {
	for _, player := range room.Players {
		if player.UserID == userID {
			return true
		}
	}
	return false
}

func isInviteCodeAccepted(room rooms.Room, inviteCode string) bool {
	if room.InviteCodeRequired {
		for _, code := range room.InviteCodes {
			if code.Code == inviteCode {
				return true
			}
		}
	}
	return false
}

func getPlayer(room rooms.Room, userID string) (rooms.Player, error) {
	for _, player := range room.Players {
		if player.UserID == userID {
			return player, nil
		}
	}
	return rooms.Player{}, rooms.ErrForbidden
}
