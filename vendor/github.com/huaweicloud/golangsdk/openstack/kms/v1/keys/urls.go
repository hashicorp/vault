package keys

import "github.com/huaweicloud/golangsdk"

const (
	resourcePath = "kms"
)

func getURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "describe-key")
}

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "create-key")
}

func deleteURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "schedule-key-deletion")
}

func updateAliasURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "update-key-alias")
}

func updateDesURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "update-key-description")
}

func dataEncryptURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "create-datakey")
}

func encryptDEKURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "encrypt-datakey")
}

func enableKeyURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "enable-key")
}

func disableKeyURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "disable-key")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "list-keys")
}
