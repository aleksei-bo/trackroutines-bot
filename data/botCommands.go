package data

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

// MainCommands are list of actions shown on the bot page
var MainCommands = []tgbotapi.BotCommand{
	{Command: "/help", Description: "Show help message"},
	{Command: "/create", Description: "Create new task"},
	{Command: "/tasks", Description: "Interact with existing tasks"},
	{Command: "/score", Description: "Show score of all participants"},
	{Command: "/monthly", Description: "Show completed tasks in that month"},
	{Command: "/list", Description: "Show all tasks"},
}

type TaskCommand struct {
	Name    string
	Command string
}

// commands appear when selecting a task after "/tasks" command
var SingleTaskCommands = []TaskCommand{
	{"Show task info", "info"},
	{"Show history of actions", "last"},
	{"Done today", "today"},
	{"Done yesterday", "yesterday"},
	{"Modify task", "modify"},
	{"Delete task", "delete"},
}
