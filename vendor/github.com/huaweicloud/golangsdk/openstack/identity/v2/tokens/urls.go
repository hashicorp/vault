package tokens

import "github.com/huaweicloud/golangsdk"

// CreateURL generates the URL used to create new Tokens.
func CreateURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("tokens")
}

// GetURL generates the URL used to Validate Tokens.
func GetURL(client *golangsdk.ServiceClient, token string) string {
	return client.ServiceURL("tokens", token)
}
