// Package tencentcloud provides node discovery for TencentCloud.
package tencentcloud

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `TencentCloud:

	provider:          "tencentcloud"
	region:            The TencentCloud region
	tag_key:           The tag key to filter on
	tag_value:         The tag value to filter on
	address_type:      "private_v4", "public_v4". (default: "private_v4")
	access_key_id:     The secret id of TencentCloud
	access_key_secret: The secret key of TencentCloud

	This required permission to 'cvm:DescribeInstances'.
	It is recommended you make a dedicated key used only for auto-joining.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "tencentcloud" {
		return nil, fmt.Errorf("discover-tencentcloud: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]
	addressType := args["address_type"]
	accessKeyID := args["access_key_id"]
	accessKeySecret := args["access_key_secret"]

	l.Printf("[DEBUG] discover-tencentcloud: Using region=%s, tag_key=%s, tag_value=%s", region, tagKey, tagValue)
	if accessKeyID == "" {
		l.Printf("[DEBUG] discover-tencentcloud: No static credentials provided")
	} else {
		l.Printf("[DEBUG] discover-tencentcloud: Static credentials provided")
	}

	if region == "" {
		l.Printf("[DEBUG] discover-tencentcloud: Region not provided")
		return nil, fmt.Errorf("discover-tencentcloud: region missing")
	}
	l.Printf("[DEBUG] discover-tencentcloud: region is %s", region)

	if addressType == "" {
		addressType = "private_v4"
	}

	if addressType != "private_v4" && addressType != "public_v4" {
		l.Printf("[DEBUG] discover-tencentcloud: Address type %s invalid", addressType)
		return nil, fmt.Errorf("discover-tencentcloud: invalid address_type " + addressType)
	}
	l.Printf("[DEBUG] discover-tencentcloud: address type is %s", addressType)

	credential := common.NewCredential(
		accessKeyID,
		accessKeySecret,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"
	cvmClient, _ := cvm.NewClient(credential, region, cpf)

	l.Printf("[DEBUG] discover-tencentcloud: Filter instances with %s=%s", tagKey, tagValue)
	request := cvm.NewDescribeInstancesRequest()
	request.Filters = []*cvm.Filter{
		{
			Name:   stringToPointer("instance-state"),
			Values: []*string{stringToPointer("RUNNING")},
		},
		{
			Name:   stringToPointer("tag:" + tagKey),
			Values: []*string{stringToPointer(tagValue)},
		},
	}

	response, err := cvmClient.DescribeInstances(request)
	if err != nil {
		l.Printf("[DEBUG] discover-tencentcloud: DescribeInstances failed, %s", err)
		return nil, fmt.Errorf("discover-tencentcloud: DescribeInstances failed, %s", err)
	}
	l.Printf("[DEBUG] discover-tencentcloud: Found %d instances", len(response.Response.InstanceSet))

	var addrs []string
	for _, v := range response.Response.InstanceSet {
		switch addressType {
		case "public_v4":
			if len(v.PublicIpAddresses) == 0 {
				l.Printf("[DEBUG] discover-tencentcloud: Instance %s has no public_v4", *v.InstanceId)
				continue
			}
			l.Printf("[DEBUG] discover-tencentcloud: Instance %s has public_v4 %v", *v.InstanceId, *v.PublicIpAddresses[0])
			addrs = append(addrs, *v.PublicIpAddresses[0])
		case "private_v4":
			if len(v.PrivateIpAddresses) == 0 {
				l.Printf("[DEBUG] discover-tencentcloud: Instance %s has no private_v4", *v.InstanceId)
				continue
			}
			l.Printf("[DEBUG] discover-tencentcloud: Instance %s has private_v4 %v", *v.InstanceId, *v.PrivateIpAddresses[0])
			addrs = append(addrs, *v.PrivateIpAddresses[0])
		}
	}

	l.Printf("[DEBUG] discover-tencentcloud: Found address: %v", addrs)
	return addrs, nil
}

func stringToPointer(s string) *string {
	return &s
}
