package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func handleCommandStats(m *tb.Message) {
	s, err := getStats()
	if err != nil {
		sendTo(m.Sender, errorText)
		return
	}

	msg := fmt.Sprintf(statsTemplate, s.usersCount,
		s.enabledUsersCount, s.postsCount, s.averagePriceFrom,
		s.averagePriceTo, s.averageRoomsFrom, s.averageRoomsTo)

	sendTo(m.Sender, msg)
}
