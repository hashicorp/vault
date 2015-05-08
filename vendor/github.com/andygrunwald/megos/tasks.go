package megos

import (
	"fmt"
)

// GetTaskByID will return a Task by its unique ID.
//
// The list of tasks are provided by a framework.
func (c *Client) GetTaskByID(tasks []Task, taskID string) (*Task, error) {
	for _, task := range tasks {
		if taskID == task.ID {
			return &task, nil
		}
	}

	return nil, fmt.Errorf("No task with id \"%s\" found", taskID)
}
