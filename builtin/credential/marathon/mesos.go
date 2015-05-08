package marathon

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/andygrunwald/megos"
)

// SlaveTaskIDIsValid ensure a valid task is running with taskID
func SlaveTaskIDIsValid(mesosURL string, taskID string) (bool, error) {

	mesosURLs := strings.Split(mesosURL, ",")
	var urls []*url.URL

	for _, mesosURL := range mesosURLs {
		url, err := url.Parse(mesosURL)

		if err != nil {
			return false, errors.New(fmt.Sprintf("Invalid mesos url %s", mesosURL))
		}
		urls = append(urls, url)
	}

	mesos := megos.NewClient(urls, nil)

	state, _ := mesos.GetStateFromCluster()
	framework, _ := mesos.GetFrameworkByPrefix(state.Frameworks, "marathon")
	_, err := mesos.GetTaskByID(framework.Tasks, taskID)

	if err != nil {
		return false, errors.New("Slave Task ID not found!")
	}

	return true, nil
}
