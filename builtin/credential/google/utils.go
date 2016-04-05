package google

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"strings"
)

func applicationOauth2Config(applicationId string, applicationSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     applicationId,
		ClientSecret: applicationSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{ "email" },
	}
}

//copied from vault/util... make public?
func strListContains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func sliceToMap(slice []string) map[string]bool{
	m := map[string]bool{}
	for _, element := range slice {
		m[element] = true
	}
	return m
}

func environmentVariable(name string) string {

	return os.Getenv(name)
}

func localPartFromEmail(email string) string {

	var name string

	if index := strings.Index(email, "@") ; index > -1 {
		name = email[:index]
	} else {
		name = email
	}

	return name
}