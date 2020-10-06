// Package os provides node discovery for Openstack.
package os

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Openstack:

    provider:   "os"
    auth_url:   The endpoint of OS identity
    project_id: The id of the project (tenant id)
    tag_key:    The tag key to filter on
    tag_value:  The tag value to filter on
    user_name:  The user used to authenticate
    password:   The password of the provided user
    token:      The token to use
    insecure:   Sets if the api certificate shouldn't be check. Any value means true

    Variables can also be provided by environmental variables.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "os" {
		return nil, fmt.Errorf("discover-os: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	projectID := args["project_id"]
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]
	var err error

	if projectID == "" { // Use the one on the instance if not provided either by parameter or env
		l.Printf("[INFO] discover-os: ProjectID not provided. Looking up in metadata...")
		projectID, err = getProjectID()
		if err != nil {
			return nil, err
		}
		l.Printf("[INFO] discover-os: ProjectID is %s", projectID)
		args["project_id"] = projectID
	}

	log.Printf("[DEBUG] discover-os: Using project_id=%s tag_key=%s tag_value=%s", projectID, tagKey, tagValue)
	client, err := newClient(args, l)
	if err != nil {
		return nil, err
	}

	if p.userAgent != "" {
		client.UserAgent.Prepend(p.userAgent)
	}

	pager := servers.List(client, ListOpts{ListOpts: servers.ListOpts{Status: "ACTIVE"}, ProjectID: projectID})
	if err := pager.Err; err != nil {
		return nil, fmt.Errorf("discover-os: ListServers failed: %s", err)
	}

	var addrs []string
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		srvs, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}
		for _, srv := range srvs {
			for key, value := range srv.Metadata {
				l.Printf("[INFO] discover-os: Filter instances with %s=%s", tagKey, tagValue)
				if key == tagKey && value == tagValue {
					// Loop over the server address and append any fixed one to the list
					for _, v := range srv.Addresses {
						if addrsInfo, ok := v.([]interface{}); ok {
							for _, addrInfo := range addrsInfo {
								if info, ok := addrInfo.(map[string]interface{}); ok {
									if info["OS-EXT-IPS:type"] == "fixed" {
										addrs = append(addrs, info["addr"].(string))
									}
								}
							}
						}
					}
				}
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover-os: ExtractServerInfo failed: %s", err)
	}

	l.Printf("[DEBUG] discover-os: Found ip addresses: %v", addrs)
	return addrs, nil
}

func newClient(args map[string]string, l *log.Logger) (*gophercloud.ServiceClient, error) {
	username := argsOrEnv(args, "user_name", "OS_USERNAME")
	password := argsOrEnv(args, "password", "OS_PASSWORD")
	token := argsOrEnv(args, "token", "OS_AUTH_TOKEN")
	url := argsOrEnv(args, "auth_url", "OS_AUTH_URL")
	region := argsOrEnv(args, "region", "OS_REGION_NAME")
	if region == "" {
		region = "RegionOne"
	}
	projectID := argsOrEnv(args, "project_id", "OS_PROJECT_ID")
	insecure := argsOrEnv(args, "insecure", "OS_INSECURE")
	domain_id := argsOrEnv(args, "domain_id", "OS_DOMAIN_ID")
	domain_name := argsOrEnv(args, "domain_name", "OS_DOMAIN_NAME")

	if url == "" {
		return nil, fmt.Errorf("discover-os: Auth url must be provided")
	}

	ao := gophercloud.AuthOptions{
		DomainID:         domain_id,
		DomainName:       domain_name,
		IdentityEndpoint: url,
		Username:         username,
		Password:         password,
		TokenID:          token,
		TenantID:         projectID,
	}

	client, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return nil, fmt.Errorf("discover-os: Client initialization failed: %s", err)
	}

	config := &tls.Config{InsecureSkipVerify: insecure != ""}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       config,
	}
	transport.TLSClientConfig = config
	client.HTTPClient = *http.DefaultClient
	client.HTTPClient.Transport = transport

	l.Printf("[DEBUG] discover-os: Authenticating...")
	if err = openstack.Authenticate(client, ao); err != nil {
		return nil, fmt.Errorf("discover-os: Authentication failed: %s", err)
	}

	l.Printf("[DEBUG] discover-os: Creating client...")
	computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{Region: region})
	if err != nil {
		return nil, fmt.Errorf("discover-os: ComputeClient initialization failed: %s", err)
	}
	return computeClient, nil
}

func argsOrEnv(args map[string]string, key, env string) string {
	if value := args[key]; value != "" {
		return value
	}
	return os.Getenv(env)
}

// ListOpts add the project to the parameters of servers.ListOpts
type ListOpts struct {
	servers.ListOpts
	ProjectID string `q:"project_id"`
}

// ToServerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToServerListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func getProjectID() (string, error) {
	resp, err := http.Get("http://169.254.169.254/openstack/latest/meta_data.json")
	if err != nil {
		return "", fmt.Errorf("discover-os: Error asking metadata for project_id: %s", err)
	}
	data := struct {
		ProjectID string `json:"project_id"`
	}{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("discover-os: Can't read response body: %s", err)
	}
	resp.Body.Close()
	if err = json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("discover-os: Can't convert project_id: %s", err)
	}
	if data.ProjectID == "" {
		return "", fmt.Errorf("discover-os: Couln't find project_id on metadata")
	}
	return data.ProjectID, nil
}
