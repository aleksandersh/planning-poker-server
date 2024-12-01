package roomsdata

import (
	"slices"

	"aleksandersh.github.io/planning-poker-server/internal/rooms"
)

func estimateGame(game rooms.Game) rooms.Game {
	suggested := 0
	sum := 0
	count := 0
	for _, card := range game.Cards {
		score := card.Score
		if score >= 0 {
			sum = sum + score
			count = count + 1
			if suggested >= 0 && suggested < score {
				suggested = score
			}
		} else {
			if suggested > score {
				suggested = score
			}
		}
	}
	game.MaxScore = suggested
	if count > 0 {
		average := sum / count
		if sum%count > 0 {
			average = average + 1
		}
		game.AverageScore = average
	} else {
		game.AverageScore = 0
	}
	return game
}

func putUserCard(cards []rooms.Card, card rooms.Card) []rooms.Card {
	idx := slices.IndexFunc(cards, func(c rooms.Card) bool {
		return c.Player.UserID == card.Player.UserID
	})
	if idx > 0 {
		cards[idx] = card
	} else {
		cards = append(cards, card)
	}
	return cards
}

func dropUserCard(cards []rooms.Card, userID string) []rooms.Card {
	return slices.DeleteFunc(cards, func(card rooms.Card) bool {
		return card.Player.UserID == userID
	})
}
