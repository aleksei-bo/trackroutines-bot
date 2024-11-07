// here are some unused functions which may be returned later

package handlers

import (
	"fmt"
	"log"

	"github.com/aleksei-bo/trackroutines-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Not used in current version of a bot
// Handles "/task_done" command and shows list of tasks as buttons.
func TaskDone(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
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
		button := tgbotapi.NewInlineKeyboardButtonData(task.Name, fmt.Sprintf("task_done_%d", task.TaskID))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Add a button to return to the previous command
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel")))

	// Create the reply markup
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	// Send the message with inline buttons
	msgText := "Please select the task that you have completed:"
	message := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	message.ReplyMarkup = replyMarkup
	bot.Send(message)
}
