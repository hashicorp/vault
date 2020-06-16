package services

import "github.com/huaweicloud/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("services")
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("services")
}

func serviceURL(client *golangsdk.ServiceClient, serviceID string) string {
	return client.ServiceURL("services", serviceID)
}

func updateURL(client *golangsdk.ServiceClient, serviceID string) string {
	return client.ServiceURL("services", serviceID)
}
