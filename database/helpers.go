package database

import (
	"log"
	"github.com/aleksei-bo/trackroutines-bot/data"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// all messages are stored for now
func StoreMessagesToDB(update tgbotapi.Update) {
	// Store message in database
	_, err := DB.Exec(`INSERT INTO allMessages (user_id, username, message_text) VALUES (?, ?, ?)`,
		update.Message.From.ID, update.Message.From.UserName, update.Message.Text)
	if err != nil {
		log.Printf("Error storing message: %v", err)

	}

	log.Printf("Stored message from %s: %s", update.Message.From.UserName, update.Message.Text)

}

func CheckPassPhrase(update tgbotapi.Update) bool {

	if update.Message.Text == data.PassPhrase {
		if err := AddUserToDB(update.Message.From.ID, update.Message.From.UserName); err != nil {
			log.Println(err)
			return false
		}
		return true
	} else {
		return false

	}
}
