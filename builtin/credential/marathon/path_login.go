package marathon

import (
	"errors"
	"fmt"
	"strings"
	"time"

	marathon "github.com/gambol99/go-marathon"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	StartupThresholdSeconds = time.Second * 300
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"marathon_app_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "MARATHON_APP_ID env var from a Marathon task",
			},
			"marathon_app_version": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "MARATHON_APP_VERSION env var from a Marathon task",
			},
			"mesos_task_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "MESOS_TASK env var from a Marathon task",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	loginTime := time.Now()

	client, err := getMarathonClientFromConfig(b, req)

	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	appId := data.Get("marathon_app_id").(string)
	appVersion := data.Get("marathon_app_version").(string)
	taskId := data.Get("mesos_task_id").(string)

	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	appTask, err := getAppTaskFromValues(client, appId, appVersion)

	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	_, err = appTaskStartedWithinThreshold(appTask, loginTime)

	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	appName := strings.TrimPrefix(appId, "/")

	mesosUrl, err := getMesosUrl(b, req)
	_, err = SlaveTaskIDIsValid(mesosUrl, taskId)

	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: []string{appName},
			Metadata: map[string]string{
				"marathon_app_id":      appName,
				"marathon_app_version": appVersion,
				"mesos_task_id":        taskId,
			},
			DisplayName: appName,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       time.Minute * 5,
			},
		},
	}, nil
}

func validateConfigIsConfigured(b *backend, req *logical.Request) (*config, error) {
	// Get all our stored state
	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	if config.MarathonUrl == "" {
		return nil, errors.New("configure the marathon credential backend first: missing marathon_url")
	}

	if config.MesosUrl == "" {
		return nil, errors.New("configure the marathon credential backend first: missing mesos_url")
	}

	return config, nil
}

func getMarathonClientFromConfig(b *backend, req *logical.Request) (marathon.Marathon, error) {
	config, err := validateConfigIsConfigured(b, req)

	if err != nil {
		return nil, err
	}

	return b.Client(config.MarathonUrl)
}

func getMesosUrl(b *backend, req *logical.Request) (string, error) {
	config, err := validateConfigIsConfigured(b, req)

	if err != nil {
		return "", err
	}

	return config.MesosUrl, nil
}

func getAppTaskFromValues(client marathon.Marathon, appId string, appVersion string) (*marathon.Task, error) {
	// Get marathon task data
	app, err := client.Application(appId)
	if err != nil {
		return nil, err
	}

	for _, task := range app.Tasks {
		if task.Version == appVersion {
			return task, nil
		}
	}

	return nil, errors.New("App version not found")
}

func appTaskStartedWithinThreshold(appTask *marathon.Task, loginTime time.Time) (bool, error) {
	startedAt, e := time.Parse(
		time.RFC3339,
		appTask.StartedAt)

	if e != nil {
		return false, errors.New(fmt.Sprintf("Failed to validate app startup time: %s", e.Error()))
	}

	delta := loginTime.Sub(startedAt)
	if delta > StartupThresholdSeconds {
		return false, errors.New("App did not startup within threshold")
	}
	return true, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	appId := req.Auth.Metadata["marathon_app_id"]
	appVersion := req.Auth.Metadata["marathon_app_version"]

	client, err := getMarathonClientFromConfig(b, req)

	if err != nil {
		return nil, err
	}

	appTask, err := getAppTaskFromValues(client, appId, appVersion)

	if err != nil {
		return nil, err
	}

	if appTask == nil {
		// not sure if this is necessary, but if appTask is nil,
		// do not renew
		return nil, nil
	}

	mesosUrl, err := getMesosUrl(b, req)
	_, err = SlaveTaskIDIsValid(mesosUrl, appTask.ID)

	if err != nil {
		return nil, err
	}

	return framework.LeaseExtend(5*time.Minute, 0, b.System())(req, d)
}
