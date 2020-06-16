package tokens

import "github.com/huaweicloud/golangsdk"

func tokenURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
