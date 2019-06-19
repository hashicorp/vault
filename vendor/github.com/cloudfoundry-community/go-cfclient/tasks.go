package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// TaskListResponse is the JSON response from the API.
type TaskListResponse struct {
	Pagination Pagination `json:"pagination"`
	Tasks      []Task     `json:"resources"`
}

// Task is a description of a task element.
type Task struct {
	GUID       string `json:"guid"`
	SequenceID int    `json:"sequence_id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	State      string `json:"state"`
	MemoryInMb int    `json:"memory_in_mb"`
	DiskInMb   int    `json:"disk_in_mb"`
	Result     struct {
		FailureReason string `json:"failure_reason"`
	} `json:"result"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DropletGUID string    `json:"droplet_guid"`
	Links       struct {
		Self    Link `json:"self"`
		App     Link `json:"app"`
		Droplet Link `json:"droplet"`
	} `json:"links"`
}

// TaskRequest is a v3 JSON object as described in:
// http://v3-apidocs.cloudfoundry.org/version/3.0.0/index.html#create-a-task
type TaskRequest struct {
	Command          string `json:"command"`
	Name             string `json:"name"`
	MemoryInMegabyte int    `json:"memory_in_mb"`
	DiskInMegabyte   int    `json:"disk_in_mb"`
	DropletGUID      string `json:"droplet_guid"`
}

func (c *Client) makeTaskListRequestWithParams(baseUrl string, query url.Values) ([]byte, error) {
	requestUrl := baseUrl + "?" + query.Encode()
	req := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting tasks")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(err, "Error requesting tasks: status code not 200, it was %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func parseTaskListRespones(answer []byte) (TaskListResponse, error) {
	var response TaskListResponse
	err := json.Unmarshal(answer, &response)
	if err != nil {
		return response, errors.Wrap(err, "Error unmarshaling response %v")
	}
	return response, nil
}

func (c *Client) handleTasksApiCall(apiUrl string, query url.Values) ([]Task, error) {
	body, err := c.makeTaskListRequestWithParams(apiUrl, query)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting tasks")
	}
	response, err := parseTaskListRespones(body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading tasks")
	}
	return response.Tasks, nil
}

// ListTasks returns all tasks the user has access to.
// See http://v3-apidocs.cloudfoundry.org/version/3.12.0/index.html#list-tasks
func (c *Client) ListTasks() ([]Task, error) {
	return c.handleTasksApiCall("/v3/tasks", url.Values{})
}

// ListTasksByQuery returns all tasks the user has access to, with query parameters.
// See http://v3-apidocs.cloudfoundry.org/version/3.12.0/index.html#list-tasks
func (c *Client) ListTasksByQuery(query url.Values) ([]Task, error) {
	return c.handleTasksApiCall("/v3/tasks", query)
}

// TasksByApp returns task structures which aligned to an app identified by the given guid.
// See: http://v3-apidocs.cloudfoundry.org/version/3.12.0/index.html#list-tasks-for-an-app
func (c *Client) TasksByApp(guid string) ([]Task, error) {
	return c.TasksByAppByQuery(guid, url.Values{})
}

// TasksByAppByQuery returns task structures which aligned to an app identified by the given guid
// and filtered by the given query parameters.
// See: http://v3-apidocs.cloudfoundry.org/version/3.12.0/index.html#list-tasks-for-an-app
func (c *Client) TasksByAppByQuery(guid string, query url.Values) ([]Task, error) {
	uri := fmt.Sprintf("/v3/apps/%s/tasks", guid)
	return c.handleTasksApiCall(uri, query)
}

func createReader(tr TaskRequest) (io.Reader, error) {
	rmap := make(map[string]string)
	rmap["command"] = tr.Command
	if tr.Name != "" {
		rmap["name"] = tr.Name
	}
	// setting droplet GUID causing issues
	if tr.MemoryInMegabyte != 0 {
		rmap["memory_in_mb"] = fmt.Sprintf("%d", tr.MemoryInMegabyte)
	}
	if tr.DiskInMegabyte != 0 {
		rmap["disk_in_mb"] = fmt.Sprintf("%d", tr.DiskInMegabyte)
	}

	bodyReader := bytes.NewBuffer(nil)
	enc := json.NewEncoder(bodyReader)
	if err := enc.Encode(rmap); err != nil {
		return nil, errors.Wrap(err, "Error during encoding task request")
	}
	return bodyReader, nil
}

// CreateTask creates a new task in CF system and returns its structure.
func (c *Client) CreateTask(tr TaskRequest) (task Task, err error) {
	bodyReader, err := createReader(tr)
	if err != nil {
		return task, err
	}

	request := fmt.Sprintf("/v3/apps/%s/tasks", tr.DropletGUID)
	req := c.NewRequestWithBody("POST", request, bodyReader)

	resp, err := c.DoRequest(req)
	if err != nil {
		return task, errors.Wrap(err, "Error creating task")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return task, errors.Wrap(err, "Error reading task after creation")
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		return task, errors.Wrap(err, "Error unmarshaling task")
	}
	return task, err
}

// GetTaskByGuid returns a task structure by requesting it with the tasks GUID.
func (c *Client) GetTaskByGuid(guid string) (task Task, err error) {
	request := fmt.Sprintf("/v3/tasks/%s", guid)
	req := c.NewRequest("GET", request)

	resp, err := c.DoRequest(req)
	if err != nil {
		return task, errors.Wrap(err, "Error requesting task")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return task, errors.Wrap(err, "Error reading task")
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		return task, errors.Wrap(err, "Error unmarshaling task")
	}
	return task, err
}

func (c *Client) TaskByGuid(guid string) (task Task, err error) {
	return c.GetTaskByGuid(guid)
}

// TerminateTask cancels a task identified by its GUID.
func (c *Client) TerminateTask(guid string) error {
	req := c.NewRequest("PUT", fmt.Sprintf("/v3/tasks/%s/cancel", guid))
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error terminating task")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return errors.Wrapf(err, "Failed terminating task, response status code %d", resp.StatusCode)
	}
	return nil
}
