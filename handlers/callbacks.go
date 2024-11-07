package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aleksei-bo/trackroutines-bot/data"
	"github.com/aleksei-bo/trackroutines-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ShowTaskCommands(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) (string, error) {

	taskID := strings.TrimPrefix(update.CallbackQuery.Data, "task #")

	var buttons [][]tgbotapi.InlineKeyboardButton

	task, err := database.GetTask(taskID)
	if err != nil {
		log.Println("couldn't get task with id:", taskID, err)
		return "", err
	}

	for _, cmd := range data.SingleTaskCommands {
		button := tgbotapi.NewInlineKeyboardButtonData(cmd.Name, fmt.Sprintf("command:%s", cmd.Command))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(button))
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel")))
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	// Send the message with inline buttons
	msgText := fmt.Sprintf("TASK: %v\nSelect action to proceed.", task.Name)
	message := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
	message.ReplyMarkup = replyMarkup
	bot.Send(message)
	data.UserTasks[update.CallbackQuery.Message.Chat.ID] = task

	return fmt.Sprintf("actions shown for %v", task.Name), nil
}

func ProcessTaskCommand(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) (string, error) {
	command := strings.TrimPrefix(update.CallbackQuery.Data, "command:")
	task := data.UserTasks[update.CallbackQuery.Message.Chat.ID]
	chatID := update.CallbackQuery.Message.Chat.ID

	switch command {
	case "info":
		msg := tgbotapi.NewMessage(chatID, TaskDescriptionLong(task))
		msg.ParseMode = "markdown"
		bot.Send(msg)

	case "last":
		msgText, err := database.GetSingleTaskActionsAsMessage(task)
		if err != nil {
			log.Println(err)
			bot.Send(tgbotapi.NewMessage(chatID, "Error retrieving actions"))
		} else {
			msg := tgbotapi.NewMessage(chatID, msgText)
			msg.ParseMode = "markdown"
			bot.Send(msg)
		}

		// do nothing
	case "today":
		if msg, err := database.MarkTaskAsDone(ctx, "", task, "today"); err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, msg))
			log.Println(err)
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, msg))
		}
	case "yesterday":
		if msg, err := database.MarkTaskAsDone(ctx, "", task, "yesterday"); err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, msg))
			log.Println(err)
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, msg))
		}
	case "modify":
		msg := tgbotapi.NewMessage(chatID, TaskDescriptionShort(task))
		msg.ParseMode = "markdown"
		bot.Send(tgbotapi.NewMessage(chatID, "copy message below and modify it:"))
		bot.Send(msg)
		data.UserStates[chatID] = "modify task"
		data.UserTasks[chatID] = task

		// do nothing
	case "delete":
		// do nothing
		bot.Send(tgbotapi.NewMessage(chatID, "You don't have access to this function"))
	}
	return "command processed", nil
}
