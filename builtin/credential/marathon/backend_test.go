package marathon

import (
	"errors"
	"os"
	"testing"
	"time"

	marathon "github.com/gambol99/go-marathon"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

const AppID = "test-app"

func buildBackend(t *testing.T) logical.Backend {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 30
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	return b
}

func TestBackend_login(t *testing.T) {
	if os.Getenv("MARATHON_URL") == "" {
		t.Skip("MARATHON_URL not set")
	}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  buildBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccLogin(t, AppID),
		},
		Teardown: func() error {
			return tearDown(AppID)
		},
	})
}

func TestBackend_invalid(t *testing.T) {
	if os.Getenv("MARATHON_URL") == "" {
		t.Skip("MARATHON_URL not set")
	}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  buildBackend(t),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccLoginInvalid(t, AppID),
		},
		Teardown: func() error {
			return tearDown(AppID)
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MARATHON_URL"); v == "" {
		t.Fatal("MARATHON_URL must be set for acceptance tests")
	}
	if v := os.Getenv("MESOS_URL"); v == "" {
		t.Fatal("MESOS_URL must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	marathonURL := os.Getenv("MARATHON_URL")
	mesosURL := os.Getenv("MESOS_URL")

	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"marathon_url": marathonURL,
			"mesos_url":    mesosURL,
		},
	}
}

func testAccLogin(t *testing.T, appID string) logicaltest.TestStep {
	marathonURL := os.Getenv("MARATHON_URL")

	c, err := marathonClient(marathonURL)
	if err != nil {
		t.Fatal(err)
	}

	task, err := startTestTask(c, appID)

	if err != nil {
		t.Fatal(err)
	}

	appVersion := task.Version
	taskID := task.ID

	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"marathon_app_id":      appID,
			"marathon_app_version": appVersion,
			"mesos_task_id":        taskID,
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth([]string{"default", appID}),
	}
}

func testAccLoginInvalid(t *testing.T, appID string) logicaltest.TestStep {
	marathonURL := os.Getenv("MARATHON_URL")

	c, err := marathonClient(marathonURL)
	if err != nil {
		t.Fatal(err)
	}

	task, err := startTestTask(c, appID)

	if err != nil {
		t.Fatal(err)
	}

	appVersion := task.Version
	taskID := task.ID

	time.Sleep(time.Second * 5)

	stopTestTask(c, appID)

	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"marathon_app_id":      appID,
			"marathon_app_version": appVersion,
			"mesos_task_id":        taskID,
		},
		ErrorOk:         true,
		Unauthenticated: true,

		Check: logicaltest.TestCheckError(),
	}
}

func tearDown(appID string) error {
	marathonURL := os.Getenv("MARATHON_URL")

	c, err := marathonClient(marathonURL)
	if err != nil {
		return err
	}

	return stopTestTask(c, appID)
}

func marathonClient(marathonURL string) (marathon.Marathon, error) {
	config := marathon.NewDefaultConfig()
	config.URL = marathonURL
	config.LogOutput = os.Stdout
	c, err := marathon.NewClient(config)

	return c, err
}

func startTestTask(c marathon.Marathon, appID string) (*marathon.Task, error) {
	stopTestTask(c, appID)

	time.Sleep(time.Second * 2)

	application := marathon.NewDockerApplication()
	application.Name(appID)
	application.AddArgs("sleep").AddArgs("10000")
	application.CPU(0.1).Memory(256).Count(1)
	application.Container.Docker.Container("alpine")

	app, err := c.CreateApplication(application)
	if err != nil {
		return nil, err
	}

	err = c.WaitOnApplication(app.ID, time.Second*30)
	if err != nil {
		return nil, err
	}

	appRead, err := c.Application(app.ID)
	if err != nil {
		return nil, err
	}

	if len(appRead.Tasks) == 0 {
		return nil, errors.New("Test App failed to start")
	}

	return appRead.Tasks[0], nil
}

func stopTestTask(c marathon.Marathon, appID string) error {

	deploymentID, err := c.DeleteApplication(appID, true)
	if err != nil {
		// app is already deleted
		return nil
	}

	err = c.WaitOnDeployment(deploymentID.DeploymentID, time.Second*120)

	if err != nil {
		return err
	}

	return nil
}
