package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"github.com/aleksei-bo/trackroutines-bot/data"
	"time"
	"unicode/utf8"

	_ "github.com/mattn/go-sqlite3"
)

// AddUserToDB adds a new user to the database
func AddUserToDB(telegramID int64, username string) error {

	_, err := DB.Exec("INSERT INTO users (telegramID, username, points, alias) VALUES (?, ?, ?, ?)", telegramID, username, 0, username)
	log.Println("user added")
	return err
}

// GetUser retrieves a user from the database by their Telegram username
func GetUser(username string) (data.User, error) {
	var u data.User
	if DB == nil {
		log.Fatal(fmt.Errorf("database not initialized"))
	}

	// Use sqlx Get method to map result to struct fields
	query := "SELECT * FROM users WHERE username = ?"
	err := DB.Get(&u, query, username)

	if err == sql.ErrNoRows {
		return u, fmt.Errorf("user not found: %w", err)
	}
	if err != nil {
		return u, fmt.Errorf("error querying user: %w", err)
	}
	u.PointsMonthly, err = GetUserMonthlyScore(u)
	if err != nil {
		return u, fmt.Errorf("error getting monthly score: %w", err)
	}

	return u, nil
}

func GetUserMonthlyScore(user data.User) (int64, error) {
	actions, err := GetListOfActions("this month")
	if err != nil {
		return 0, err
	}
	var points int64
	for _, a := range actions {
		if a.UserID == user.UserID {
			points += a.Points
		}
	}
	return points, nil

}

// Fetches all open tasks from the database
func GetWaitingTasks() ([]data.Task, error) {
	var tasks []data.Task
	err := DB.Select(&tasks, "SELECT * FROM tasks WHERE status = 'waiting'")
	return tasks, err
}

func GetTask(taskID string) (data.Task, error) {
	var task data.Task
	err := DB.Get(&task, "SELECT * FROM tasks WHERE taskID = ?", taskID)

	return task, err
}

func MarkTaskAsDone(ctx context.Context, taskID string, task data.Task, day string) (string, error) {
	// Extract user from context
	user, ok := ctx.Value("userCtx").(data.User)
	msg := "Couldn't mark task as done"
	var date time.Time
	switch day {
	case "today":
		date = time.Now()
	case "yesterday":
		date = time.Now().AddDate(0, 0, -1)
	default:
		date = time.Now()
	}

	var err error
	if !ok {
		return msg, fmt.Errorf("user not found in context")
	}
	if taskID != "" {
		task, err = GetTask(taskID)
		if err != nil {
			return msg, err
		}

	}

	// Create a timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Start a new transaction
	tx, err := DB.BeginTx(ctx, nil) // Use context here to respect the timeout
	if err != nil {
		return msg, fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute the first query
	_, err = tx.ExecContext(ctx, "UPDATE tasks SET doneLast = ? WHERE taskID = ?",
		time.Now(), taskID)
	if err != nil {
		return msg, fmt.Errorf("can't update task status: %w", err)
	}
	// Execute the second query
	_, err = tx.ExecContext(ctx, "INSERT INTO actions (taskID, userID, timestamp) VALUES (?,?,?)", task.TaskID, user.UserID, date)

	/* 	// TESTING
	   	// Use prevMonth instead of time.Now() in your SQL insert
	   	prevMonth := time.Now().AddDate(0, -1, 0)
	   	_, err = tx.ExecContext(ctx, "INSERT INTO actions (taskID, userID, timestamp) VALUES (?,?,?)",
	   		task.TaskID, user.UserID, prevMonth) */

	if err != nil {
		return msg, fmt.Errorf("can't insert new action: %w", err)
	}

	// Execute the third query
	_, err = tx.ExecContext(ctx, "UPDATE users SET points = points + ? WHERE userID = ?",
		task.Points, user.UserID)
	if err != nil {
		return msg, fmt.Errorf("can't update user points: %w", err)
	}

	// Commit the transaction if all queries succeed
	if err = tx.Commit(); err != nil {
		return msg, fmt.Errorf("could not commit transaction: %w", err)
	}
	log.Printf("Task '%v' is done by %v, who has %v points now", task.Name, user.Username, (user.Points + task.Points))
	msg = fmt.Sprintf(
		"Task --%v-- marked as done %v! \nYour monthly score is %v points and total score is %v points.", task.Name, day, user.PointsMonthly+task.Points, (user.Points + task.Points))

	return msg, nil
}

func GetAllUsers() ([]data.User, error) {
	var users []data.User
	err := DB.Select(&users, "SELECT * FROM users ORDER BY points DESC")
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	return users, nil
}

func GetMonthlyScore(period string, users []data.User) (string, error) {

	actions, _ := GetListOfActions(period)

	msg := ""
	for _, u := range users {
		for _, a := range actions {

			if u.UserID == a.UserID {
				u.PointsMonthly += a.Points
			}
		}
		msg += fmt.Sprintf("- %v: \t %v points\n", u.Alias, u.PointsMonthly)
	}

	return msg, nil

}

func GetTotalScore(users []data.User) string {
	msg := ""

	for _, u := range users {
		msg += fmt.Sprintf("- %v: \t %v points\n", u.Alias, u.Points)
	}

	return msg
}

func GetListOfTasks() (string, error) {
	var tasks []data.Task
	err := DB.Select(&tasks, "SELECT * FROM tasks WHERE status = 'waiting'")
	if len(tasks) == 0 {
		return "No tasks added yet", nil
	}
	if err != nil {
		return "", fmt.Errorf("error retrieving tasks: %w", err)
	}

	message := "Tasks:\n"

	for _, task := range tasks {
		// Adjust the name length for non-Latin characters
		nameLen := utf8.RuneCountInString(task.Name)
		if nameLen < 25 {
			task.Name += strings.Repeat(" ", 25-nameLen)
		} else {
			runes := []rune(task.Name)
			task.Name = string(runes[:22]) + "... "
		}

		message += fmt.Sprintf("➤ %s | %4d p. | %-15s\n",
			task.Name, task.Points, task.DoneLast.Format("02 Jan"))
	}
	log.Println(message)
	return message, nil
}

func GetListOfActions(period string, sortBy ...string) ([]data.Action, error) {
	var actions []data.Action

	// Get the current time
	now := time.Now()
	var startOfPeriod time.Time
	var endOfPeriod = time.Now()
	switch period {
	case "this week":

	case "this month":
		startOfPeriod = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "previous month":
		startOfPeriod = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		endOfPeriod = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	default:
		startOfPeriod = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	}

	var query string

	// two queries to avoid sql injection
	if len(sortBy) > 0 && (sortBy[0] == "DESC") {
		query = `
		SELECT a.actionID, a.taskID, a.userID, a.timestamp, t.name AS taskName, t.points AS points, u.alias AS alias
		FROM actions a
		JOIN tasks t ON a.taskID = t.taskID
		JOIN users u ON a.userID = u.userID
		WHERE a.timestamp >= ? AND a.timestamp <= ?
		ORDER BY a.timestamp DESC`
	} else {
		query = `
	SELECT a.actionID, a.taskID, a.userID, a.timestamp, t.name AS taskName, t.points AS points, u.alias AS alias
	FROM actions a
	JOIN tasks t ON a.taskID = t.taskID
	JOIN users u ON a.userID = u.userID
	WHERE a.timestamp >= ? AND a.timestamp <= ?
	ORDER BY a.timestamp ASC`
	}

	// Execute the query to get all actions starting from the first of the current month
	err := DB.Select(&actions, query, startOfPeriod, endOfPeriod) // asc - older first
	if err != nil {
		return []data.Action{}, fmt.Errorf("error fetching actions: %w", err)
	}

	return actions, nil
}

func GetSingleTaskActionsAsMessage(task data.Task) (string, error) {
	taskID := task.TaskID
	var actions []data.Action
	var msgText string

	err := DB.Select(&actions, `SELECT a.actionID, a.taskID, a.userID, a.timestamp, t.name AS taskName, t.points AS points, u.alias AS alias
	FROM actions a
	JOIN tasks t ON a.taskID = t.taskID
	JOIN users u ON a.userID = u.userID
	WHERE a.taskID = ? ORDER BY timestamp DESC`, taskID)

	// err := DB.Select(&actions, "SELECT * FROM actions WHERE taskID = ? ORDER BY timestamp DESC LIMIT 10", taskID)
	if err != nil {
		return msgText, fmt.Errorf("error fetching actions: %w", err)
	}
	if len(actions) == 0 {
		return fmt.Sprintf("No actions for task _%v_ yet", task.Name), nil
	} else {
		msgText = fmt.Sprintf("Last %v actions for task _%v_:\n\n", len(actions), task.Name)
		for _, action := range actions {
			msgText += fmt.Sprintf("➤ Completed by %v on %s | + %v p.\n", action.Alias, action.Timestamp.Format("02 Jan"), action.Points)
		}

		return msgText, nil
	}
}

// UpdateUserPoints updates the points for a given user
func UpdateUserPoints(userID int64, points int) error {
	_, err := DB.Exec("UPDATE users SET points = points + ? WHERE userID = ?", points, userID)
	return err
}

func SetAlias(userID int64, alias string) error {
	_, err := DB.Exec("UPDATE users SET alias = ? WHERE userID = ?", alias, userID)
	return err
}

// InsertTask adds a new task to the database
func InsertTask(name, description string, points, periodicity int64) error {
	_, err := DB.Exec("INSERT INTO tasks (name, description, points, periodicity) VALUES (?, ?, ?, ?)", name, description, points, periodicity)
	log.Println("inserted into tasks: ", name, description, points, periodicity)
	return err
}

// InsertTask adds a new task to the database
func UpdateTask(id int64, name, description string, points, periodicity int64) error {
	_, err := DB.Exec("UPDATE tasks SET name = ?, description = ?, points = ?, periodicity = ? WHERE taskID = ?", name, description, points, periodicity, id)
	return err
}

/* // GetAvailableTasks retrieves all tasks that are not marked as done
func GetAvailableTasks() ([]Task, error) {
	rows, err := DB.Query("SELECT taskID, name, description, points FROM tasks WHERE status = 'waiting'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.TaskID, &t.Name, &t.Description, &t.Points)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
} */
