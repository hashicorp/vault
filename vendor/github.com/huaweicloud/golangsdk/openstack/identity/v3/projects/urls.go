package projects

import "github.com/huaweicloud/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("projects")
}

func getURL(client *golangsdk.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("projects")
}

func deleteURL(client *golangsdk.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func updateURL(client *golangsdk.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}
