package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Task represents a  Task
type Task struct {
	// Identifier is a unique identifier for the task
	Identifier string `json:"id,omitempty"`

	// StartDate is the start date of the task
	StartDate string `json:"started_at,omitempty"`

	// TerminationDate is the termination date of the task
	TerminationDate string `json:"terminated_at,omitempty"`

	HrefFrom string `json:"href_from,omitempty"`

	Description string `json:"description,omitempty"`

	Status string `json:"status,omitempty"`

	Progress int `json:"progress,omitempty"`
}

// oneTask represents the response of a GET /tasks/UUID API call
type oneTask struct {
	Task Task `json:"task,omitempty"`
}

// Tasks represents a group of  tasks
type Tasks struct {
	// Tasks holds  tasks of the response
	Tasks []Task `json:"tasks,omitempty"`
}

// GetTasks get the list of tasks from the API
func (s *API) GetTasks() ([]Task, error) {
	query := url.Values{}
	// TODO per_page=20&page=2
	resp, err := s.GetResponsePaginate(s.computeAPI, "tasks", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var tasks Tasks

	if err = json.Unmarshal(body, &tasks); err != nil {
		return nil, err
	}
	return tasks.Tasks, nil
}

// GetTask fetches a specific task
func (s *API) GetTask(id string) (*Task, error) {
	query := url.Values{}
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("tasks/%s", id), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}

	var t oneTask
	if err = json.Unmarshal(body, &t); err != nil {
		return nil, err
	}
	return &t.Task, nil
}
