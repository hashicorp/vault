package marathon

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	var data struct {
		MarathonAppId string `mapstructure:"marathon_app_id"`
		MarathonAppVersion string `mapstructure:"marathon_app_version"`
		MesosTaskId string `mapstructure:"mesos_task_id"`
		Mount    string `mapstructure:"mount"`
	}
	if err := mapstructure.WeakDecode(m, &data); err != nil {
		return "", err
	}

	if data.MarathonAppId == "" || data.MarathonAppVersion == "" || data.MesosTaskId == "" {
		return "", fmt.Errorf("All of 'marathon_app_id', 'marathon_app_version' and 'mesos_task_id' must be specified")
	}
	if data.Mount == "" {
		data.Mount = "marathon"
	}

	path := fmt.Sprintf("auth/%s/login", data.Mount)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"marathon_app_id":      data.MarathonAppId,
		"marathon_app_version": data.MarathonAppVersion,
		"mesos_task_id":        data.MesosTaskId,
	})
	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("empty response from credential provider")
	}
	
	return secret.Auth.ClientToken, nil
}

func (h *CLIHandler) Help() string {
	help := `
The Marathon credential provider allows you to authenticate via Marathon tasks.
To use it, specify "marathon_app_id", "marathon_app_version" and "mesos_task_id.

    Example: vault auth -method=marathon marathon_app_id=<app_id> marathon_app_version=<app_version> mesos_task_id=<task_id>

	`

	return strings.TrimSpace(help)
}
