package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aleksei-bo/trackroutines-bot/data"
	"github.com/aleksei-bo/trackroutines-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Create(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	data.UserStates[chatID] = "create task"
	bot.Send(tgbotapi.NewMessage(chatID, data.MessageCreateTask))
}

func CreateNewTask(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userInput := message.Text

	chatID := message.Chat.ID

	// Reset the user's state

	if strings.ToLower(userInput) == "cancel" || strings.ToLower(userInput) == "no" || strings.ToLower(userInput) == "нет" {
		bot.Send(tgbotapi.NewMessage(chatID, "Task creation cancelled."))
		data.UserStates[chatID] = ""
		return
	}

	// setting default values
	var task data.Task

	parts := strings.Split(userInput, ";")
	if len(parts) < 2 {
		bot.Send(tgbotapi.NewMessage(chatID, "Invalid input format. Please use: Task Name; Points \nExample: wash dishes; 1"))
		// sql table fields which may present: Task Name; Points; Periodicity; Status; description
		return
	}
	var err error
	task.Name = strings.TrimSpace(parts[0])
	pointsInt, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	task.Points = int64(pointsInt)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Invalid points value. Please enter a number."))
		return
	}

	// Optional fields (Periodicity, Status, and Description) - filled only if provided
	if len(parts) > 2 {
		periodInt, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		task.Periodicity = int64(periodInt)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Invalid periodicity value. Please enter a number between 1 and 30."))
			return
		}
	}
	if len(parts) > 3 {
		task.Description = strings.TrimSpace(parts[3]) // optional
	}

	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"Task preview:\nName: %s\nPoints: %d\nPeriodicity: %d\nDescription: %s", task.Name, task.Points, task.Periodicity, task.Description)))
	bot.Send(tgbotapi.NewMessage(chatID, "Do you want to save the task? Yes/No"))

	data.UserStates[chatID] = "verify task"
	data.UserTasks[chatID] = task

}

// for creating new tasks
func SubmitTask(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	task := data.UserTasks[message.Chat.ID]
	userInput := message.Text
	chatID := message.Chat.ID
	data.UserStates[chatID] = ""
	data.UserTasks[chatID] = data.Task{}

	if strings.ToLower(userInput) == "yes" || strings.ToLower(userInput) == "y" || strings.ToLower(userInput) == "да" {
		err := database.InsertTask(task.Name, task.Description, task.Points, 7)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Error creating task"))
			log.Println("couldn't insert task", err)
			return
		}
		bot.Send(tgbotapi.NewMessage(chatID, "Task successfully created"))

		return
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, "Task creation cancelled."))
	}

}
