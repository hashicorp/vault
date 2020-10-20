package tenants

import "github.com/huaweicloud/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("tenants")
}

func getURL(client *golangsdk.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("tenants")
}

func deleteURL(client *golangsdk.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}

func updateURL(client *golangsdk.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}
