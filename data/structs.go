package data

import (
	"time"
)

var PassPhrase string

// User struct maps to the 'users' table
type User struct {
	UserID        int64  `db:"userID"`
	TelegramID    string `db:"telegramID"`
	Username      string `db:"username"`
	Points        int64  `db:"points"`
	Alias         string `db:"alias"`
	PointsMonthly int64
}

// Task struct maps to the 'tasks' table
type Task struct {
	TaskID      int64     `db:"taskID"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Points      int64     `db:"points"`
	Periodicity int64     `db:"periodicity"`
	Status      string    `db:"status"`
	DoneLast    time.Time `db:"doneLast"`
	Category    string    `db:"category"`
}

// Action struct maps to the 'actions' table
type Action struct {
	ActionID int64  `db:"actionID"`
	TaskID   int64  `db:"taskID"`
	TaskName string `db:"taskName"`
	Points   int64  `db:"points"`

	UserID    int64     `db:"userID"`
	Alias     string    `db:"alias"`
	Timestamp time.Time `db:"timestamp"`
}

var UserStates = make(map[int64]string)

var UserTasks = make(map[int64]Task)
