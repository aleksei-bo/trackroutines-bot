package handlers

import (
	"fmt"

	"github.com/aleksei-bo/trackroutines-bot/data"
)

// returns full task description
func TaskDescriptionLong(task data.Task) string {
	taskInfo := fmt.Sprintf(
		`_%v_,
		- *Description*: %v 
		- *Points*: %v 
		- *Periodicity*: %v 
		- *Status*: %v 
		- *Done last on*: %v
		`, task.Name, task.Description, task.Points, task.Periodicity, task.Status, task.DoneLast.Format("02 Jan"))

	return taskInfo
}

// returns task info , separated by semicolons
func TaskDescriptionShort(task data.Task) string {
	taskInfo := fmt.Sprintf(
		"%v; %v; %d; %d", task.Name, task.Description, task.Points, task.Periodicity)

	return taskInfo
}
