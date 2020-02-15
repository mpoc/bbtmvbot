package main

import tb "gopkg.in/tucnak/telebot.v2"

func handleCommandEnable(m *tb.Message) {
	query := "UPDATE users SET enabled=1 WHERE id=?"
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
	sendTo(m.Sender, "Pranešimai įjungti! Naudokite komandą /disable kad juos išjungti.\n\n"+ActiveSettingsText)
}
