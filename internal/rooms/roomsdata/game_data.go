package roomsdata

import (
	"slices"

	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdomain/roomsmodel"
)

func estimateGame(game roomsmodel.Game) roomsmodel.Game {
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

func putUserCard(cards []roomsmodel.Card, card roomsmodel.Card) []roomsmodel.Card {
	idx := slices.IndexFunc(cards, func(c roomsmodel.Card) bool {
		return c.Player.UserID == card.Player.UserID
	})
	if idx > 0 {
		cards[idx] = card
	} else {
		cards = append(cards, card)
	}
	return cards
}

func dropUserCard(cards []roomsmodel.Card, userID string) []roomsmodel.Card {
	return slices.DeleteFunc(cards, func(card roomsmodel.Card) bool {
		return card.Player.UserID == userID
	})
}
