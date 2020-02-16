package main

import tb "gopkg.in/tucnak/telebot.v2"

// phase -1 = not configured (user cannot enable notifications)
// phase 0 = all set and user can enable notifications
// phase 1 = choose type
// phase 2 = choose city
// phase 3 = choose price range
// phase 4 = choose built years range
// phase 5 = choose rooms count range
// phase 6 = choose broker fee (if type is rent)

var typeKeys = [][]tb.ReplyButton{
	// Nuomotis = 0
	// Pirkti = 1
	[]tb.ReplyButton{tb.ReplyButton{Text: "Nuomotis butą"}, tb.ReplyButton{Text: "Pirkti butą"}},
}

var cityKeys = [][]tb.ReplyButton{
	// Vilniuje = 0
	// Kaune = 1
	// Klaipėdoje = 2
	// Šiauliuose = 3
	// Panevėžyje = 4
	// Alytuje = 5
	[]tb.ReplyButton{tb.ReplyButton{Text: "Vilniuje"}, tb.ReplyButton{Text: "Kaune"}},
	[]tb.ReplyButton{tb.ReplyButton{Text: "Klaipėdoje"}, tb.ReplyButton{Text: "Šiauliuose"}},
	[]tb.ReplyButton{tb.ReplyButton{Text: "Panevėžyje"}, tb.ReplyButton{Text: "Alytuje"}},
}

var feeKeys = [][]tb.ReplyButton{
	[]tb.ReplyButton{tb.ReplyButton{Text: "Tik su tarpininkavimo mokesčiu"}},
	[]tb.ReplyButton{tb.ReplyButton{Text: "Tik be tarpininkavimo mokesčio"}},
	[]tb.ReplyButton{tb.ReplyButton{Text: "Tarpininkavimo mokestis nėra svarbus"}},
}

func handleCommandConfig(m *tb.Message) {
	userID := m.Sender.ID
	setConfPhase(userID, 1)

	bot.Send(m.Sender, "Ką norite daryti?", &tb.ReplyMarkup{
		ReplyKeyboard: typeKeys,
	})
}

// Bot will listen to every message text from user, and ignore if conf_phase is -1
func handleCommandText(m *tb.Message) {
	userID := m.Sender.ID
	confPhase := getConfPhase(userID)

	switch confPhase {
	case 0:
		return
	case 1:
		switch m.Text {
		case "Nuomotis butą":
			setType(userID, 0)
			bot.Send(m.Sender, "Kuriame mieste norite nuomotis butą?", &tb.ReplyMarkup{
				ReplyKeyboard: cityKeys,
			})
		case "Pirkti butą":
			setType(userID, 1)
			bot.Send(m.Sender, "Kuriame mieste norite pirkti butą?", &tb.ReplyMarkup{
				ReplyKeyboard: cityKeys,
			})
		default:
			sendTo(&tb.User{ID: userID}, "Pasirinkite, ką norite daryti, naudodamiesi pateiktais mygtukais")
			return
		}

		setConfPhase(userID, 2)
	case 2:
		selectedType := getType(userID)
		var firstSentence string

		switch m.Text {
		case "Vilniuje":
			setCity(userID, 0)
			firstSentence = "Pasirinkote Vilnių."
		case "Kaune":
			setCity(userID, 1)
			firstSentence = "Pasirinkote Kauną."
		case "Klaipėdoje":
			setCity(userID, 2)
			firstSentence = "Pasirinkote Klaipėdą."
		case "Šiauliuose":
			setCity(userID, 3)
			firstSentence = "Pasirinkote Šiaulius."
		case "Panevėžyje":
			setCity(userID, 4)
			firstSentence = "Pasirinkote Panevėžį."
		case "Alytuje":
			setCity(userID, 5)
			firstSentence = "Pasirinkote Alytų."
		default:
			if selectedType == 0 {
				sendTo(&tb.User{ID: userID}, "Pasirinkite, kuriame mieste norite nuomotis butą, naudodamiesi pateiktais mygtukais")
			} else /*if selectedType == 1*/ {
				sendTo(&tb.User{ID: userID}, "Pasirinkite, kuriame mieste norite pirkti butą, naudodamiesi pateiktais mygtukais")
			}
			return
		}

		if selectedType == 0 {
			sendTo(&tb.User{ID: userID}, firstSentence+" Įveskite norimo nuomotis buto kainos intervalą eurais (pvz. \"200 - 500\" arba \"- 500\" jei norite matyti visus butus iki 500 eur). Taip pat galite rašyti \"nesvarbu\" jei nenorite filtruoti skelbimų pagal kainą.", &tb.ReplyMarkup{
				ReplyKeyboardRemove: true,
			})
		} else /*if selectedType == 1*/ {
			sendTo(&tb.User{ID: userID}, firstSentence+" Įveskite norimo pirkti buto kainos intervalą eurais (pvz. \"50000 - 150000\" arba \"- 150000\" jei norite matyti visus butus iki 150 000 eur). Taip pat galite rašyti \"nesvarbu\" jei nenorite filtruoti skelbimų pagal kainą.", &tb.ReplyMarkup{
				ReplyKeyboardRemove: true,
			})
		}

		setConfPhase(userID, 3)
	case 3:
		// TODO
	}
}

func getConfPhase(userID int) int {
	query := "SELECT conf_phase FROM users WHERE id=? LIMIT 1"
	var confPhase int
	err := db.QueryRow(query, userID).Scan(&confPhase)
	if err != nil {
		panic(err)
	}
	return confPhase
}

func setConfPhase(userID, confPhase int) {
	query := "UPDATE users SET conf_phase=? WHERE id=?"
	_, err := db.Exec(query, confPhase, userID)
	if err != nil {
		panic(err)
	}
}

func getType(userID int) int {
	query := "SELECT type FROM users WHERE id=? LIMIT 1"
	var selectedType int
	err := db.QueryRow(query, userID).Scan(&selectedType)
	if err != nil {
		panic(err)
	}
	return selectedType
}

func setType(userID, selectedType int) {
	query := "UPDATE users SET type=? WHERE id=?"
	_, err := db.Exec(query, selectedType, userID)
	if err != nil {
		panic(err)
	}
}

func setCity(userID, city int) {
	query := "UPDATE users SET city=? WHERE id=?"
	_, err := db.Exec(query, city, userID)
	if err != nil {
		panic(err)
	}
}
