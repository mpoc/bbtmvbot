package main

import tb "gopkg.in/tucnak/telebot.v2"

func handleCommandDisable(m *tb.Message) {
	query := "UPDATE users SET enabled=0 WHERE id=?"
	_, err := db.Exec(query, m.Sender.ID)
	if err != nil {
		sendTo(m.Sender, errorText)
		return
	}

	ActiveSettingsText, err := getActiveSettingsText(m.Sender)
	if err != nil {
		sendTo(m.Sender, errorText)
		return
	}
	sendTo(m.Sender, "Pranešimai išjungti! Naudokite komandą /enable kad juos įjungti.\n\n"+ActiveSettingsText)
}
