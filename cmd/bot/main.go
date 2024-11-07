package main

import (
	"io"
	"log"
	"os"
	"github.com/aleksei-bo/trackroutines-bot/data"
	"github.com/aleksei-bo/trackroutines-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is not set")
	}

	data.PassPhrase = os.Getenv("BOT_AUTH")
	if data.PassPhrase == "" {
		data.PassPhrase = "this bot is cool"
		log.Println("BOT_AUTH environment variable is not set, using default value")
	}

	var err error

	err = database.InitDB()
	if err != nil {
		log.Fatal("can't initialize database:", err)
	}
	defer database.DB.Close()

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("can't initialize bot:", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// commands are list of actions shown on the bot page
	_, err = bot.Request(tgbotapi.NewSetMyCommands(data.MainCommands...))
	if err != nil {
		log.Panic("can't set commands:", err)
	}

	// Open or create a log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	defer file.Close()

	// Set log output to the file
	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)

	// Configure update channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	handleUpdates(bot, updates)

}

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64, username string) {
	msg := tgbotapi.NewMessage(chatID, "Welcome, "+username+"! Do you want to set a custom username to use in the bot? Type your preferred username or 'no' to skip.")
	bot.Send(msg)
	data.UserStates[chatID] = "set alias"

}

func sendErrorMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "You are not authorized to use this bot.")
	bot.Send(msg)
}
