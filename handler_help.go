package main

import tb "gopkg.in/tucnak/telebot.v2"

func handleCommandHelp(m *tb.Message) {
	ActiveSettingsText, err := getActiveSettingsText(m.Sender)
	if err != nil {
		sendTo(m.Sender, errorText)
		return
	}
	sendTo(m.Sender, helpText+"\n\n"+ActiveSettingsText)
}
