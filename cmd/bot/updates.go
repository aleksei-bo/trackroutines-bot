package main

import (
	"context"
	"log"
	"strings"
	"github.com/aleksei-bo/trackroutines-bot/data"
	"github.com/aleksei-bo/trackroutines-bot/database"

	"github.com/aleksei-bo/trackroutines-bot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.CallbackQuery == nil && update.Message == nil {
			continue
		}

		if update.CallbackQuery != nil {
			user, err := database.GetUser(update.CallbackQuery.From.UserName)
			if err != nil {
				log.Println(err)
				sendErrorMessage(bot, update.CallbackQuery.Message.Chat.ID)
				continue
			}
			ctx := context.WithValue(context.Background(), "userCtx", user)

			// Now you can call handleCallbackQuery with the context
			handleCallbackQuery(ctx, bot, update)
			continue
		}

		database.StoreMessagesToDB(update)

		user, err := database.GetUser(update.Message.From.UserName)
		if err != nil {
			log.Println(err)
			if !database.CheckPassPhrase(update) {
				sendErrorMessage(bot, update.Message.Chat.ID)
			} else {
				sendWelcomeMessage(bot, update.Message.Chat.ID, update.Message.From.UserName)
			}
			continue
		}
		ctx := context.WithValue(context.Background(), "userCtx", user)
		handleCommands(ctx, bot, update)

	}
}

// Handle callback queries from inline keyboard
func handleCallbackQuery(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	callbackData := update.CallbackQuery.Data
	switch {
	case strings.HasPrefix(callbackData, "task #"):
		handlers.ShowTaskCommands(ctx, bot, update)

	case strings.HasPrefix(callbackData, "command:"):
		handlers.ProcessTaskCommand(ctx, bot, update)

	case callbackData == "cancel":
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Cancelled."))
	}

	delMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	if _, err := bot.Request(delMsg); err != nil {
		log.Println("couldnt delete message", err)
	}

}

func handleCommands(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: Implement command handling
	user := ctx.Value("userCtx").(data.User)

	_ = user

	message := update.Message
	chatID := message.Chat.ID

	// if user state is already affected (in process of something)
	switch data.UserStates[chatID] {
	case "create task":
		handlers.CreateNewTask(bot, message)
		return
	case "verify task":
		handlers.SubmitTask(bot, message)
		return
	case "set alias":
		handlers.SetAlias(bot, message, user)
		return
	case "modify task":
		handlers.ModifyTask(bot, message)
		return
	case "submit updated task":
		handlers.SubmitUpdatedTask(bot, message)
		return
	}

	/* if !checkPassPhrase(update) {
		sendErrorMessage(bot, chatID)
		return
	} */

	command := message.Command()

	switch command {
	case "start":
		handlers.Start(bot, message)
	case "help":
		handlers.Help(bot, message)
	case "score":
		handlers.Score(bot, message)
	case "list":
		handlers.ListOfTasks(bot, message)
	case "create":
		handlers.Create(bot, message)
	case "tasks":
		handlers.SelectTask(bot, message)
	case "monthly":
		handlers.MonthlyActions(bot, message)
	default:
		handlers.UnknownCommand(bot, message)
	}

}
