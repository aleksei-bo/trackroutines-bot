package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aleksei-bo/trackroutines-bot/data"
	"github.com/aleksei-bo/trackroutines-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler for /start command
func Start(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	response := "Welcome! Use /help to see available commands."
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, response))
}

// Handler for /help command
func Help(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	response := "Here are the commands you can use:\n"
	for _, cmd := range data.MainCommands {
		response += cmd.Command + ": " + cmd.Description + "\n"
	}
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, response))
}

// new user can select alias
func SetAlias(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user data.User) {
	userInput := strings.ToLower(msg.Text)
	chatID := msg.Chat.ID

	if userInput == "no" || userInput == "нет" {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Your alias is your telegram username: %v", msg.From.UserName)))
		data.UserStates[chatID] = ""
		return
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Your new alias is: %v", msg.Text)))

		database.SetAlias(user.UserID, msg.Text)
		data.UserStates[chatID] = ""
		return
	}
}

// this is a command which will manage all interactions with existing tasks. currently command is /tasks
func SelectTask(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	tasks, err := database.GetWaitingTasks()
	if err != nil {
		log.Println(err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Error retrieving tasks"))
		return
	}

	if len(tasks) == 0 {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "No open tasks available."))
		return
	}

	// Create inline keyboard buttons for tasks
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, task := range tasks {
		// here each button shows task name and points
		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v: %d p.", task.Name, task.Points), fmt.Sprintf("task #%d", task.TaskID))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add a button to return to the previous command
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel")))

	// Create the reply markup
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	// Send the message with inline buttons
	msgText := "Please select a task to preview actions."
	message := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	message.ReplyMarkup = replyMarkup
	bot.Send(message)
}

func MonthlyActions(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	actions, err := database.GetListOfActions("this month", "DESC")
	if err != nil {
		log.Println(err)
		return
	}
	if len(actions) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "No actions yet in this month")
		bot.Send(msg)
		return
	}
	// Format the actions into a message string
	msgText := fmt.Sprintf("Actions in the current month(%v), newest first:\n\n", time.Now().Month())
	for _, action := range actions {
		msgText += fmt.Sprintf("➤ '%v' by %v on %s | + %v p.\n", action.TaskName, action.Alias, action.Timestamp.Format("02 Jan"), action.Points)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	bot.Send(msg)
}

// shows list of participants and their score to telegram user
func Score(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	users, err := database.GetAllUsers()
	if err != nil {
		log.Println(err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Error retrieving score"))
	}

	thisMonthScore, err := database.GetMonthlyScore("this month", users)
	if err != nil {
		log.Println(err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Error retrieving score"))
		return
	}

	prevMonthScore, err := database.GetMonthlyScore("previous month", users)
	if err != nil {
		log.Println(err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Error retrieving score"))
		return
	}
	totalScore := database.GetTotalScore(users)

	msg := fmt.Sprintf(`Current Score in %v:
		%v
	
	Final Score in Previous Month:
	%v

	Total Points since start:
	%v
	`, time.Now().Month(), thisMonthScore, prevMonthScore, totalScore)

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, msg))

}

// shows list of all tasks to telegram user
func ListOfTasks(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	list, err := database.GetListOfTasks()
	if err != nil {
		log.Println(err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Error retrieving score"))
		return
	}
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, list))
}

func UnknownCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Please use one of the available commands.")
	bot.Send(msg)
}
