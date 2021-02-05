package linodego

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/linode/linodego/pkg/errors"
	"github.com/linode/linodego/pkg/retry"
)

const (
	// APIHost Linode API hostname
	APIHost = "api.linode.com"
	// APIHostVar environment var to check for alternate API URL
	APIHostVar = "LINODE_URL"
	// APIHostCert environment var containing path to CA cert to validate against
	APIHostCert = "LINODE_CA"
	// APIVersion Linode API version
	APIVersion = "v4"
	// APIVersionVar environment var to check for alternate API Version
	APIVersionVar = "LINODE_API_VERSION"
	// APIProto connect to API with http(s)
	APIProto = "https"
	// Version of linodego
	Version = "0.12.0"
	// APIEnvVar environment var to check for API token
	APIEnvVar = "LINODE_TOKEN"
	// APISecondsPerPoll how frequently to poll for new Events or Status in WaitFor functions
	APISecondsPerPoll = 3
	// Maximum wait time for retries
	APIRetryMaxWaitTime = time.Duration(30) * time.Second
	// DefaultUserAgent is the default User-Agent sent in HTTP request headers
	DefaultUserAgent = "linodego " + Version + " https://github.com/linode/linodego"
)

var (
	envDebug = false

	defaultRetryConditionals = []retry.ConditionFunc{
		retry.LinodeBusyRetryCondition,
		retry.ServiceUnavailableRetryCondition,
		retry.TooManyRequestsRetryCondition,
	}
)

// Client is a wrapper around the Resty client
type Client struct {
	resty             *resty.Client
	userAgent         string
	resources         map[string]*Resource
	debug             bool
	retryConditionals []retry.ConditionFunc

	millisecondsPerPoll time.Duration

	Account                *Resource
	AccountSettings        *Resource
	DomainRecords          *Resource
	Domains                *Resource
	Events                 *Resource
	Firewalls              *Resource
	FirewallDevices        *Resource
	FirewallRules          *Resource
	IPAddresses            *Resource
	IPv6Pools              *Resource
	IPv6Ranges             *Resource
	Images                 *Resource
	InstanceConfigs        *Resource
	InstanceDisks          *Resource
	InstanceIPs            *Resource
	InstanceSnapshots      *Resource
	InstanceStats          *Resource
	InstanceVolumes        *Resource
	Instances              *Resource
	InvoiceItems           *Resource
	Invoices               *Resource
	Kernels                *Resource
	LKEClusters            *Resource
	LKEClusterAPIEndpoints *Resource
	LKEClusterPools        *Resource
	LKEVersions            *Resource
	Longview               *Resource
	LongviewClients        *Resource
	LongviewSubscriptions  *Resource
	Managed                *Resource
	NodeBalancerConfigs    *Resource
	NodeBalancerNodes      *Resource
	NodeBalancerStats      *Resource
	NodeBalancers          *Resource
	Notifications          *Resource
	OAuthClients           *Resource
	ObjectStorageBuckets   *Resource
	ObjectStorageClusters  *Resource
	ObjectStorageKeys      *Resource
	Payments               *Resource
	Profile                *Resource
	Regions                *Resource
	SSHKeys                *Resource
	StackScripts           *Resource
	Tags                   *Resource
	Tickets                *Resource
	Token                  *Resource
	Tokens                 *Resource
	Types                  *Resource
	Users                  *Resource
	Volumes                *Resource
}

func init() {
	// Wether or not we will enable Resty debugging output
	if apiDebug, ok := os.LookupEnv("LINODE_DEBUG"); ok {
		if parsed, err := strconv.ParseBool(apiDebug); err == nil {
			envDebug = parsed
			log.Println("[INFO] LINODE_DEBUG being set to", envDebug)
		} else {
			log.Println("[WARN] LINODE_DEBUG should be an integer, 0 or 1")
		}
	}
}

// SetUserAgent sets a custom user-agent for HTTP requests
func (c *Client) SetUserAgent(ua string) *Client {
	c.userAgent = ua
	c.resty.SetHeader("User-Agent", c.userAgent)

	return c
}

// R wraps resty's R method
func (c *Client) R(ctx context.Context) *resty.Request {
	return c.resty.R().
		ExpectContentType("application/json").
		SetHeader("Content-Type", "application/json").
		SetContext(ctx).
		SetError(errors.APIError{})
}

// SetDebug sets the debug on resty's client
func (c *Client) SetDebug(debug bool) *Client {
	c.debug = debug
	c.resty.SetDebug(debug)

	return c
}

// SetBaseURL sets the base URL of the Linode v4 API (https://api.linode.com/v4)
func (c *Client) SetBaseURL(url string) *Client {
	c.resty.SetHostURL(url)
	return c
}

// SetAPIVersion sets the version of the API to interface with
func (c *Client) SetAPIVersion(apiVersion string) *Client {
	c.SetBaseURL(fmt.Sprintf("%s://%s/%s", APIProto, APIHost, apiVersion))
	return c
}

// SetRootCertificate adds a root certificate to the underlying TLS client config
func (c *Client) SetRootCertificate(path string) *Client {
	c.resty.SetRootCertificate(path)
	return c
}

// SetToken sets the API token for all requests from this client
// Only necessary if you haven't already provided an http client to NewClient() configured with the token.
func (c *Client) SetToken(token string) *Client {
	c.resty.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	return c
}

// SetRetryConditionals sets the conditions that will be evaluated on a request error to
// determine whether the request should be retried.
func (c *Client) SetRetryConditionals(conditionals []retry.ConditionFunc) *Client {
	c.retryConditionals = conditionals
	return c
}

// AddRetryConditional adds a request retry conditional to the underlying client.
func (c *Client) AddRetryConditional(retryConditional retry.ConditionFunc) *Client {
	c.retryConditionals = append(c.retryConditionals, retryConditional)
	return c
}

func (c *Client) SetRetryMaxWaitTime(max time.Duration) *Client {
	c.resty.SetRetryMaxWaitTime(max)
	return c
}

// SetPollDelay sets the number of milliseconds to wait between events or status polls.
// Affects all WaitFor* functions and retries.
func (c *Client) SetPollDelay(delay time.Duration) *Client {
	c.millisecondsPerPoll = delay
	c.resty.SetRetryWaitTime(delay * time.Millisecond)
	return c
}

// Resource looks up a resource by name
func (c Client) Resource(resourceName string) *Resource {
	selectedResource, ok := c.resources[resourceName]
	if !ok {
		log.Fatalf("Could not find resource named '%s', exiting.", resourceName)
	}

	return selectedResource
}

// NewClient factory to create new Client struct
func NewClient(hc *http.Client) (client Client) {
	if hc != nil {
		client.resty = resty.NewWithClient(hc)
	} else {
		client.resty = resty.New()
	}

	client.SetUserAgent(DefaultUserAgent)

	baseURL, baseURLExists := os.LookupEnv(APIHostVar)

	if baseURLExists {
		client.SetBaseURL(baseURL)
	} else {
		apiVersion, apiVersionExists := os.LookupEnv(APIVersionVar)
		if apiVersionExists {
			client.SetAPIVersion(apiVersion)
		} else {
			client.SetAPIVersion(APIVersion)
		}
	}

	certPath, certPathExists := os.LookupEnv(APIHostCert)

	if certPathExists {
		cert, err := ioutil.ReadFile(certPath)
		if err != nil {
			log.Fatalf("[ERROR] Error when reading cert at %s: %s\n", certPath, err.Error())
		}

		client.SetRootCertificate(certPath)

		if envDebug {
			log.Printf("[DEBUG] Set API root certificate to %s with contents %s\n", certPath, cert)
		}
	}

	client.
		SetPollDelay(1000 * APISecondsPerPoll).
		SetDebug(envDebug).
		configureRetries()

	addResources(&client)

	return
}

func (c *Client) configureRetries() {
	c.SetRetryConditionals(defaultRetryConditionals)
	c.resty.SetRetryCount(1000).
		AddRetryCondition(resty.RetryConditionFunc(c.retryConditional())).
		SetRetryAfter(retry.RespectRetryAfter)
}

// retryConditional composes a retry.ConditionFunc, which evaluates all the underlying
// conditions in c.retryConditionals.
func (c Client) retryConditional() retry.ConditionFunc {
	return func(r *resty.Response, err error) bool {
		for _, retryConditional := range c.retryConditionals {
			retry := retryConditional(r, err)
			if retry {
				log.Printf("[INFO] Received error %s - Retrying", r.Error())
				return true
			}
		}
		return false
	}
}

// nolint
func addResources(client *Client) {
	resources := map[string]*Resource{
		accountName:                NewResource(client, accountName, accountEndpoint, false, Account{}, nil),                         // really?
		accountSettingsName:        NewResource(client, accountSettingsName, accountSettingsEndpoint, false, AccountSettings{}, nil), // really?
		domainRecordsName:          NewResource(client, domainRecordsName, domainRecordsEndpoint, true, DomainRecord{}, DomainRecordsPagedResponse{}),
		domainsName:                NewResource(client, domainsName, domainsEndpoint, false, Domain{}, DomainsPagedResponse{}),
		eventsName:                 NewResource(client, eventsName, eventsEndpoint, false, Event{}, EventsPagedResponse{}),
		firewallsName:              NewResource(client, firewallsName, firewallsEndpoint, false, Firewall{}, FirewallsPagedResponse{}),
		firewallDevicesName:        NewResource(client, firewallDevicesName, firewallDevicesEndpoint, true, FirewallDevice{}, FirewallDevicesPagedResponse{}),
		firewallRulesName:          NewResource(client, firewallRulesName, firewallRulesEndpoint, true, FirewallRule{}, nil),
		imagesName:                 NewResource(client, imagesName, imagesEndpoint, false, Image{}, ImagesPagedResponse{}),
		instanceConfigsName:        NewResource(client, instanceConfigsName, instanceConfigsEndpoint, true, InstanceConfig{}, InstanceConfigsPagedResponse{}),
		instanceDisksName:          NewResource(client, instanceDisksName, instanceDisksEndpoint, true, InstanceDisk{}, InstanceDisksPagedResponse{}),
		instanceIPsName:            NewResource(client, instanceIPsName, instanceIPsEndpoint, true, InstanceIP{}, nil), // really?
		instanceSnapshotsName:      NewResource(client, instanceSnapshotsName, instanceSnapshotsEndpoint, true, InstanceSnapshot{}, nil),
		instanceStatsName:          NewResource(client, instanceStatsName, instanceStatsEndpoint, true, InstanceStats{}, nil),
		instanceVolumesName:        NewResource(client, instanceVolumesName, instanceVolumesEndpoint, true, nil, InstanceVolumesPagedResponse{}), // really?
		instancesName:              NewResource(client, instancesName, instancesEndpoint, false, Instance{}, InstancesPagedResponse{}),
		invoiceItemsName:           NewResource(client, invoiceItemsName, invoiceItemsEndpoint, true, InvoiceItem{}, InvoiceItemsPagedResponse{}),
		invoicesName:               NewResource(client, invoicesName, invoicesEndpoint, false, Invoice{}, InvoicesPagedResponse{}),
		ipaddressesName:            NewResource(client, ipaddressesName, ipaddressesEndpoint, false, nil, IPAddressesPagedResponse{}), // really?
		ipv6poolsName:              NewResource(client, ipv6poolsName, ipv6poolsEndpoint, false, nil, IPv6PoolsPagedResponse{}),       // really?
		ipv6rangesName:             NewResource(client, ipv6rangesName, ipv6rangesEndpoint, false, IPv6Range{}, IPv6RangesPagedResponse{}),
		kernelsName:                NewResource(client, kernelsName, kernelsEndpoint, false, LinodeKernel{}, LinodeKernelsPagedResponse{}),
		lkeClusterAPIEndpointsName: NewResource(client, lkeClusterAPIEndpointsName, lkeClusterAPIEndpointsEndpoint, true, LKEClusterAPIEndpoint{}, LKEClusterAPIEndpointsPagedResponse{}),
		lkeClustersName:            NewResource(client, lkeClustersName, lkeClustersEndpoint, false, LKECluster{}, LKEClustersPagedResponse{}),
		lkeClusterPoolsName:        NewResource(client, lkeClusterPoolsName, lkeClusterPoolsEndpoint, true, LKEClusterPool{}, LKEClusterPoolsPagedResponse{}),
		lkeVersionsName:            NewResource(client, lkeVersionsName, lkeVersionsEndpoint, false, LKEVersion{}, LKEVersionsPagedResponse{}),
		longviewName:               NewResource(client, longviewName, longviewEndpoint, false, nil, nil), // really?
		longviewclientsName:        NewResource(client, longviewclientsName, longviewclientsEndpoint, false, LongviewClient{}, LongviewClientsPagedResponse{}),
		longviewsubscriptionsName:  NewResource(client, longviewsubscriptionsName, longviewsubscriptionsEndpoint, false, LongviewSubscription{}, LongviewSubscriptionsPagedResponse{}),
		managedName:                NewResource(client, managedName, managedEndpoint, false, nil, nil), // really?
		nodebalancerconfigsName:    NewResource(client, nodebalancerconfigsName, nodebalancerconfigsEndpoint, true, NodeBalancerConfig{}, NodeBalancerConfigsPagedResponse{}),
		nodebalancernodesName:      NewResource(client, nodebalancernodesName, nodebalancernodesEndpoint, true, NodeBalancerNode{}, NodeBalancerNodesPagedResponse{}),
		nodebalancerStatsName:      NewResource(client, nodebalancerStatsName, nodebalancerStatsEndpoint, true, NodeBalancerStats{}, nil),
		nodebalancersName:          NewResource(client, nodebalancersName, nodebalancersEndpoint, false, NodeBalancer{}, NodeBalancerConfigsPagedResponse{}),
		notificationsName:          NewResource(client, notificationsName, notificationsEndpoint, false, Notification{}, NotificationsPagedResponse{}),
		oauthClientsName:           NewResource(client, oauthClientsName, oauthClientsEndpoint, false, OAuthClient{}, OAuthClientsPagedResponse{}),
		objectStorageBucketsName:   NewResource(client, objectStorageBucketsName, objectStorageBucketsEndpoint, false, ObjectStorageBucket{}, ObjectStorageBucketsPagedResponse{}),
		objectStorageClustersName:  NewResource(client, objectStorageClustersName, objectStorageClustersEndpoint, false, ObjectStorageCluster{}, ObjectStorageClustersPagedResponse{}),
		objectStorageKeysName:      NewResource(client, objectStorageKeysName, objectStorageKeysEndpoint, false, ObjectStorageKey{}, ObjectStorageKeysPagedResponse{}),
		paymentsName:               NewResource(client, paymentsName, paymentsEndpoint, false, Payment{}, PaymentsPagedResponse{}),
		profileName:                NewResource(client, profileName, profileEndpoint, false, nil, nil), // really?
		regionsName:                NewResource(client, regionsName, regionsEndpoint, false, Region{}, RegionsPagedResponse{}),
		sshkeysName:                NewResource(client, sshkeysName, sshkeysEndpoint, false, SSHKey{}, SSHKeysPagedResponse{}),
		stackscriptsName:           NewResource(client, stackscriptsName, stackscriptsEndpoint, false, Stackscript{}, StackscriptsPagedResponse{}),
		tagsName:                   NewResource(client, tagsName, tagsEndpoint, false, Tag{}, TagsPagedResponse{}),
		ticketsName:                NewResource(client, ticketsName, ticketsEndpoint, false, Ticket{}, TicketsPagedResponse{}),
		tokensName:                 NewResource(client, tokensName, tokensEndpoint, false, Token{}, TokensPagedResponse{}),
		typesName:                  NewResource(client, typesName, typesEndpoint, false, LinodeType{}, LinodeTypesPagedResponse{}),
		usersName:                  NewResource(client, usersName, usersEndpoint, false, User{}, UsersPagedResponse{}),
		volumesName:                NewResource(client, volumesName, volumesEndpoint, false, Volume{}, VolumesPagedResponse{}),
	}

	client.resources = resources

	client.Account = resources[accountName]
	client.DomainRecords = resources[domainRecordsName]
	client.Domains = resources[domainsName]
	client.Events = resources[eventsName]
	client.Firewalls = resources[firewallsName]
	client.FirewallDevices = resources[firewallDevicesName]
	client.FirewallRules = resources[firewallRulesName]
	client.IPAddresses = resources[ipaddressesName]
	client.IPv6Pools = resources[ipv6poolsName]
	client.IPv6Ranges = resources[ipv6rangesName]
	client.Images = resources[imagesName]
	client.InstanceConfigs = resources[instanceConfigsName]
	client.InstanceDisks = resources[instanceDisksName]
	client.InstanceIPs = resources[instanceIPsName]
	client.InstanceSnapshots = resources[instanceSnapshotsName]
	client.InstanceStats = resources[instanceStatsName]
	client.InstanceVolumes = resources[instanceVolumesName]
	client.Instances = resources[instancesName]
	client.Invoices = resources[invoicesName]
	client.Kernels = resources[kernelsName]
	client.LKEClusterAPIEndpoints = resources[lkeClusterAPIEndpointsName]
	client.LKEClusters = resources[lkeClustersName]
	client.LKEClusterPools = resources[lkeClusterPoolsName]
	client.LKEVersions = resources[lkeVersionsName]
	client.Longview = resources[longviewName]
	client.LongviewSubscriptions = resources[longviewsubscriptionsName]
	client.Managed = resources[managedName]
	client.NodeBalancerConfigs = resources[nodebalancerconfigsName]
	client.NodeBalancerNodes = resources[nodebalancernodesName]
	client.NodeBalancerStats = resources[nodebalancerStatsName]
	client.NodeBalancers = resources[nodebalancersName]
	client.Notifications = resources[notificationsName]
	client.OAuthClients = resources[oauthClientsName]
	client.ObjectStorageBuckets = resources[objectStorageBucketsName]
	client.ObjectStorageClusters = resources[objectStorageClustersName]
	client.ObjectStorageKeys = resources[objectStorageKeysName]
	client.Payments = resources[paymentsName]
	client.Profile = resources[profileName]
	client.Regions = resources[regionsName]
	client.SSHKeys = resources[sshkeysName]
	client.StackScripts = resources[stackscriptsName]
	client.Tags = resources[tagsName]
	client.Tickets = resources[ticketsName]
	client.Tokens = resources[tokensName]
	client.Types = resources[typesName]
	client.Users = resources[usersName]
	client.Volumes = resources[volumesName]
}

func copyBool(bPtr *bool) *bool {
	if bPtr == nil {
		return nil
	}

	var t = *bPtr

	return &t
}

func copyInt(iPtr *int) *int {
	if iPtr == nil {
		return nil
	}

	var t = *iPtr

	return &t
}

func copyString(sPtr *string) *string {
	if sPtr == nil {
		return nil
	}

	var t = *sPtr

	return &t
}

func copyTime(tPtr *time.Time) *time.Time {
	if tPtr == nil {
		return nil
	}

	var t = *tPtr

	return &t
}
