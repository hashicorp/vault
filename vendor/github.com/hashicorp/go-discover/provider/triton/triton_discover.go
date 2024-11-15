// Package aws provides node discovery for Joyent Triton.
package triton

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/compute"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `triton:

    provider:     "triton"
    account: 	  The Triton account name
    key_id:       The Triton KeyID
    url:          The Triton URL
    tag_key:      The tag key to filter on
    tag_value:    The tag value to filter on
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "triton" {
		return nil, fmt.Errorf("discover-triton: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	account := args["account"]
	keyID := args["key_id"]
	url := args["url"]
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]

	l.Printf("[INFO] discover-triton: Account is %q", account)
	l.Printf("[INFO] discover-triton: URL is %q", url)

	input := authentication.SSHAgentSignerInput{
		KeyID:       keyID,
		AccountName: account,
	}
	signer, err := authentication.NewSSHAgentSigner(input)
	if err != nil {
		return nil, fmt.Errorf("error Creating SSH Agent Signer: %v", err)
	}

	config := &triton.ClientConfig{
		TritonURL:   url,
		AccountName: account,
		Signers:     []authentication.Signer{signer},
	}

	c, err := compute.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("error constructing Compute Client: %v", err)
	}

	t := make(map[string]interface{}, 0)
	t[tagKey] = tagValue

	listInput := &compute.ListInstancesInput{
		Tags: t,
	}
	instances, err := c.Instances().List(context.Background(), listInput)
	if err != nil {
		return nil, fmt.Errorf("error getting instance list: %v", err)
	}
	var addrs []string
	for _, instance := range instances {
		l.Printf("[DEBUG] Instance ID: %q", instance.ID)
		l.Printf("[DEBUG] Instance PrimaryIP: %q", instance.PrimaryIP)
		if instance.PrimaryIP == "" {
			l.Printf("[DEBUG] discover-triton: Instance %s has no marked PrimaryIP", instance.ID)
			continue
		} else {
			addrs = append(addrs, instance.PrimaryIP)
		}
	}

	return addrs, nil
}
