package openstack

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/huaweicloud/golangsdk"
	tokens2 "github.com/huaweicloud/golangsdk/openstack/identity/v2/tokens"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/domains"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/endpoints"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/projects"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/services"
	tokens3 "github.com/huaweicloud/golangsdk/openstack/identity/v3/tokens"
	"github.com/huaweicloud/golangsdk/openstack/utils"
	"github.com/huaweicloud/golangsdk/pagination"
)

const (
	// v2 represents Keystone v2.
	// It should never increase beyond 2.0.
	v2 = "v2.0"

	// v3 represents Keystone v3.
	// The version can be anything from v3 to v3.x.
	v3 = "v3"
)

/*
NewClient prepares an unauthenticated ProviderClient instance.
Most users will probably prefer using the AuthenticatedClient function
instead.

This is useful if you wish to explicitly control the version of the identity
service that's used for authentication explicitly, for example.

A basic example of using this would be:

	ao, err := openstack.AuthOptionsFromEnv()
	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	client, err := openstack.NewIdentityV3(provider, golangsdk.EndpointOpts{})
*/
func NewClient(endpoint string) (*golangsdk.ProviderClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	u.RawQuery, u.Fragment = "", ""

	var base string
	versionRe := regexp.MustCompile("v[0-9.]+/?")
	if version := versionRe.FindString(u.Path); version != "" {
		base = strings.Replace(u.String(), version, "", -1)
	} else {
		base = u.String()
	}

	endpoint = golangsdk.NormalizeURL(endpoint)
	base = golangsdk.NormalizeURL(base)

	p := new(golangsdk.ProviderClient)
	p.IdentityBase = base
	p.IdentityEndpoint = endpoint
	p.UseTokenLock()

	return p, nil
}

/*
AuthenticatedClient logs in to an OpenStack cloud found at the identity endpoint
specified by the options, acquires a token, and returns a Provider Client
instance that's ready to operate.

If the full path to a versioned identity endpoint was specified  (example:
http://example.com:5000/v3), that path will be used as the endpoint to query.

If a versionless endpoint was specified (example: http://example.com:5000/),
the endpoint will be queried to determine which versions of the identity service
are available, then chooses the most recent or most supported version.

Example:

	ao, err := openstack.AuthOptionsFromEnv()
	provider, err := openstack.AuthenticatedClient(ao)
	client, err := openstack.NewNetworkV2(client, golangsdk.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
*/
func AuthenticatedClient(options golangsdk.AuthOptions) (*golangsdk.ProviderClient, error) {
	client, err := NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	err = Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authenticate or re-authenticate against the most recent identity service
// supported at the provided endpoint.
func Authenticate(client *golangsdk.ProviderClient, options golangsdk.AuthOptionsProvider) error {
	versions := []*utils.Version{
		{ID: v2, Priority: 20, Suffix: "/v2.0/"},
		{ID: v3, Priority: 30, Suffix: "/v3/"},
	}

	chosen, endpoint, err := utils.ChooseVersion(client, versions)
	if err != nil {
		return err
	}

	authOptions, isTokenAuthOptions := options.(golangsdk.AuthOptions)

	if isTokenAuthOptions {
		switch chosen.ID {
		case v2:
			return v2auth(client, endpoint, authOptions, golangsdk.EndpointOpts{})
		case v3:
			if authOptions.AgencyDomainName != "" && authOptions.AgencyName != "" {
				return v3authWithAgency(client, endpoint, &authOptions, golangsdk.EndpointOpts{})
			}
			return v3auth(client, endpoint, &authOptions, golangsdk.EndpointOpts{})
		default:
			// The switch statement must be out of date from the versions list.
			return fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
		}
	} else {
		akskAuthOptions, isAkSkOptions := options.(golangsdk.AKSKAuthOptions)

		if isAkSkOptions {
			if akskAuthOptions.AgencyDomainName != "" && akskAuthOptions.AgencyName != "" {
				return authWithAgencyByAKSK(client, endpoint, akskAuthOptions, golangsdk.EndpointOpts{})
			}
			return v3AKSKAuth(client, endpoint, akskAuthOptions, golangsdk.EndpointOpts{})

		}
		return fmt.Errorf("Unrecognized auth options provider: %s", reflect.TypeOf(options))
	}

}

// AuthenticateV2 explicitly authenticates against the identity v2 endpoint.
func AuthenticateV2(client *golangsdk.ProviderClient, options golangsdk.AuthOptions, eo golangsdk.EndpointOpts) error {
	return v2auth(client, "", options, eo)
}

func v2auth(client *golangsdk.ProviderClient, endpoint string, options golangsdk.AuthOptions, eo golangsdk.EndpointOpts) error {
	v2Client, err := NewIdentityV2(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

	v2Opts := tokens2.AuthOptions{
		IdentityEndpoint: options.IdentityEndpoint,
		Username:         options.Username,
		Password:         options.Password,
		TenantID:         options.TenantID,
		TenantName:       options.TenantName,
		AllowReauth:      options.AllowReauth,
		TokenID:          options.TokenID,
	}

	result := tokens2.Create(v2Client, v2Opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	if options.AllowReauth {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v2auth(client, endpoint, options, eo)
		}
	}
	client.TokenID = token.ID
	client.ProjectID = token.Tenant.ID
	client.EndpointLocator = func(opts golangsdk.EndpointOpts) (string, error) {
		return V2EndpointURL(catalog, opts)
	}

	return nil
}

// AuthenticateV3 explicitly authenticates against the identity v3 service.
func AuthenticateV3(client *golangsdk.ProviderClient, options tokens3.AuthOptionsBuilder, eo golangsdk.EndpointOpts) error {
	return v3auth(client, "", options, eo)
}

func v3auth(client *golangsdk.ProviderClient, endpoint string, opts tokens3.AuthOptionsBuilder, eo golangsdk.EndpointOpts) error {
	// Override the generated service endpoint with the one returned by the version endpoint.
	v3Client, err := NewIdentityV3(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

	result := tokens3.Create(v3Client, opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	project, err := result.ExtractProject()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	client.TokenID = token.ID
	if project != nil {
		client.ProjectID = project.ID
		client.DomainID = project.Domain.ID
	}

	if opts.CanReauth() {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v3auth(client, endpoint, opts, eo)
		}
	}
	client.EndpointLocator = func(opts golangsdk.EndpointOpts) (string, error) {
		return V3EndpointURL(catalog, opts)
	}

	return nil
}

func v3authWithAgency(client *golangsdk.ProviderClient, endpoint string, opts *golangsdk.AuthOptions, eo golangsdk.EndpointOpts) error {
	if opts.TokenID == "" {
		err := v3auth(client, endpoint, opts, eo)
		if err != nil {
			return err
		}
	} else {
		client.TokenID = opts.TokenID
	}

	opts1 := golangsdk.AgencyAuthOptions{
		AgencyName:       opts.AgencyName,
		AgencyDomainName: opts.AgencyDomainName,
		DelegatedProject: opts.DelegatedProject,
	}

	return v3auth(client, endpoint, &opts1, eo)
}

func getEntryByServiceId(entries []tokens3.CatalogEntry, serviceId string) *tokens3.CatalogEntry {
	if entries == nil {
		return nil
	}

	for idx, _ := range entries {
		if entries[idx].ID == serviceId {
			return &entries[idx]
		}
	}

	return nil
}

func getProjectID(client *golangsdk.ServiceClient, name string) (string, error) {
	opts := projects.ListOpts{
		Name: name,
	}
	allPages, err := projects.List(client, opts).AllPages()
	if err != nil {
		return "", err
	}

	projects, err := projects.ExtractProjects(allPages)

	if err != nil {
		return "", err
	}

	if len(projects) < 1 {
		return "", fmt.Errorf("[DEBUG] cannot find the tenant: %s", name)
	}

	return projects[0].ID, nil
}

func v3AKSKAuth(client *golangsdk.ProviderClient, endpoint string, options golangsdk.AKSKAuthOptions, eo golangsdk.EndpointOpts) error {
	v3Client, err := NewIdentityV3(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

	defer func() {
		v3Client.AKSKAuthOptions.ProjectId = options.ProjectId
		v3Client.AKSKAuthOptions.DomainID = options.DomainID
	}()
	v3Client.AKSKAuthOptions = options
	v3Client.AKSKAuthOptions.ProjectId = ""
	v3Client.AKSKAuthOptions.DomainID = ""

	if options.ProjectId == "" && options.ProjectName != "" {
		id, err := getProjectID(v3Client, options.ProjectName)
		if err != nil {
			return err
		}
		options.ProjectId = id
	}

	if options.DomainID == "" && options.Domain != "" {
		id, err := getDomainID(options.Domain, v3Client)
		if err != nil {
			options.DomainID = ""
		} else {
			options.DomainID = id
		}
	}

	if options.BssDomainID == "" && options.BssDomain != "" {
		id, err := getDomainID(options.BssDomain, v3Client)
		if err != nil {
			options.BssDomainID = ""
		} else {
			options.BssDomainID = id
		}
	}

	client.ProjectID = options.ProjectId
	client.DomainID = options.BssDomainID
	v3Client.ProjectID = options.ProjectId

	var entries = make([]tokens3.CatalogEntry, 0, 1)
	err = services.List(v3Client, services.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		serviceLst, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}

		for _, svc := range serviceLst {
			entry := tokens3.CatalogEntry{
				Type: svc.Type,
				//Name: svc.Name,
				ID: svc.ID,
			}
			entries = append(entries, entry)
		}

		return true, nil
	})

	if err != nil {
		return err
	}

	err = endpoints.List(v3Client, endpoints.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		endpoints, err := endpoints.ExtractEndpoints(page)
		if err != nil {
			return false, err
		}

		for _, endpoint := range endpoints {
			entry := getEntryByServiceId(entries, endpoint.ServiceID)

			if entry != nil {
				entry.Endpoints = append(entry.Endpoints, tokens3.Endpoint{
					URL:       strings.Replace(endpoint.URL, "$(tenant_id)s", options.ProjectId, -1),
					Region:    endpoint.Region,
					Interface: string(endpoint.Availability),
					ID:        endpoint.ID,
				})
			}
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	client.EndpointLocator = func(opts golangsdk.EndpointOpts) (string, error) {
		return V3EndpointURL(&tokens3.ServiceCatalog{
			Entries: entries,
		}, opts)
	}
	return nil
}

func authWithAgencyByAKSK(client *golangsdk.ProviderClient, endpoint string, opts golangsdk.AKSKAuthOptions, eo golangsdk.EndpointOpts) error {

	err := v3AKSKAuth(client, endpoint, opts, eo)
	if err != nil {
		return err
	}

	v3Client, err := NewIdentityV3(client, eo)
	if err != nil {
		return err
	}

	if v3Client.AKSKAuthOptions.DomainID == "" {
		return fmt.Errorf("Must config domain name")
	}

	opts2 := golangsdk.AgencyAuthOptions{
		AgencyName:       opts.AgencyName,
		AgencyDomainName: opts.AgencyDomainName,
		DelegatedProject: opts.DelegatedProject,
	}
	result := tokens3.Create(v3Client, &opts2)
	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	project, err := result.ExtractProject()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	client.TokenID = token.ID
	if project != nil {
		client.ProjectID = project.ID
	}

	client.ReauthFunc = func() error {
		client.TokenID = ""
		return authWithAgencyByAKSK(client, endpoint, opts, eo)
	}

	client.EndpointLocator = func(opts golangsdk.EndpointOpts) (string, error) {
		return V3EndpointURL(catalog, opts)
	}

	client.AKSKAuthOptions.AccessKey = ""
	return nil
}

func getDomainID(name string, client *golangsdk.ServiceClient) (string, error) {
	old := client.Endpoint
	defer func() { client.Endpoint = old }()

	client.Endpoint = old + "auth/"

	opts := domains.ListOpts{
		Name: name,
	}
	allPages, err := domains.List(client, &opts).AllPages()
	if err != nil {
		return "", fmt.Errorf("List domains failed, err=%s", err)
	}

	all, err := domains.ExtractDomains(allPages)
	if err != nil {
		return "", fmt.Errorf("Extract domains failed, err=%s", err)
	}

	count := len(all)
	switch count {
	case 0:
		err := &golangsdk.ErrResourceNotFound{}
		err.ResourceType = "iam"
		err.Name = name
		return "", err
	case 1:
		return all[0].ID, nil
	default:
		err := &golangsdk.ErrMultipleResourcesFound{}
		err.ResourceType = "iam"
		err.Name = name
		err.Count = count
		return "", err
	}
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the
// v2 identity service.
func NewIdentityV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	endpoint := client.IdentityBase + "v2.0/"
	clientType := "identity"
	var err error
	if !reflect.DeepEqual(eo, golangsdk.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	return &golangsdk.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

// NewIdentityV3 creates a ServiceClient that may be used to access the v3
// identity service.
func NewIdentityV3(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	endpoint := client.IdentityBase + "v3/"
	clientType := "identity"
	var err error
	if !reflect.DeepEqual(eo, golangsdk.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	// Ensure endpoint still has a suffix of v3.
	// This is because EndpointLocator might have found a versionless
	// endpoint and requests will fail unless targeted at /v3.
	if !strings.HasSuffix(endpoint, "v3/") {
		endpoint = endpoint + "v3/"
	}

	return &golangsdk.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

func initClientOpts(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts, clientType string) (*golangsdk.ServiceClient, error) {
	sc := new(golangsdk.ServiceClient)
	eo.ApplyDefaults(clientType)
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return sc, err
	}
	sc.ProviderClient = client
	sc.Endpoint = url
	sc.Type = clientType
	return sc, nil
}

// initcommonServiceClient create a ServiceClient which can not get from clientType directly.
// firstly, we initialize a service client by "volumev2" type, the endpoint likes https://evs.{region}.{xxx.com}/v2/{project_id}
// then we replace the endpoint with the specified srv and version.
func initcommonServiceClient(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts, srv string, version string) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "volumev2")
	if err != nil {
		return nil, err
	}

	e := strings.Replace(sc.Endpoint, "v2", version, 1)
	sc.Endpoint = strings.Replace(e, "evs", srv, 1)
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// TODO: Need to change to apig client type from apig once available
// ApiGateWayV1 creates a service client that is used for Huawei cloud for API gateway.
func ApiGateWayV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "apig", 1)
	sc.ResourceBase = sc.Endpoint + "v1.0/apigw/"
	return sc, err
}

// NewObjectStorageV1 creates a ServiceClient that may be used with the v1
// object storage package.
func NewObjectStorageV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "object-store")
}

// NewComputeV2 creates a ServiceClient that may be used with the v2 compute
// package.
func NewComputeV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "compute")
}

// NewNetworkV2 creates a ServiceClient that may be used with the v2 network
// package.
func NewNetworkV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, err
}

// NewBlockStorageV1 creates a ServiceClient that may be used to access the v1
// block storage service.
func NewBlockStorageV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "volume")
}

// NewBlockStorageV2 creates a ServiceClient that may be used to access the v2
// block storage service.
func NewBlockStorageV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "volumev2")
}

// NewBlockStorageV3 creates a ServiceClient that may be used to access the v3 block storage service.
func NewBlockStorageV3(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "volumev3")
}

// NewSharedFileSystemV2 creates a ServiceClient that may be used to access the v2 shared file system service.
func NewSharedFileSystemV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "sharev2")
}

// NewCDNV1 creates a ServiceClient that may be used to access the v1
// CDN service.
func NewCDNV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = "https://cdn.myhuaweicloud.com/"
	sc.ResourceBase = sc.Endpoint + "v1.0/"
	return sc, err
}

// NewOrchestrationV1 creates a ServiceClient that may be used to access the v1
// orchestration service.
func NewOrchestrationV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "orchestration")
}

// NewDBV1 creates a ServiceClient that may be used to access the v1 DB service.
func NewDBV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "database")
}

// NewDNSV2 creates a ServiceClient that may be used to access the v2 DNS
// service.
func NewDNSV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "dns")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

// NewImageServiceV1 creates a ServiceClient that may be used to access the v1
// image service.
func NewImageServiceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "image")
	sc.ResourceBase = sc.Endpoint + "v1/"
	return sc, err
}

// NewImageServiceV2 creates a ServiceClient that may be used to access the v2
// image service.
func NewImageServiceV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "image")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

// NewLoadBalancerV2 creates a ServiceClient that may be used to access the v2
// load balancer service.
func NewLoadBalancerV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "load-balancer")
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, err
}

// NewOtcV1 creates a ServiceClient that may be used with the v1 network package.
func NewElbV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts, otctype string) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "compute")
	//fmt.Printf("client=%+v.\n", sc)
	sc.Endpoint = strings.Replace(strings.Replace(sc.Endpoint, "ecs", otctype, 1), "/v2/", "/v1.0/", 1)
	//fmt.Printf("url=%s.\n", sc.Endpoint)
	sc.ResourceBase = sc.Endpoint
	sc.Type = otctype
	return sc, err
}

// NewSmnServiceV2 creates a ServiceClient that may be used to access the v2 Simple Message Notification service.
func NewSmnServiceV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {

	sc, err := initClientOpts(client, eo, "compute")
	sc.Endpoint = strings.Replace(sc.Endpoint, "ecs", "smn", 1)
	sc.ResourceBase = sc.Endpoint + "notifications/"
	sc.Type = "smn"
	return sc, err
}

//NewRdsServiceV1 creates the a ServiceClient that may be used to access the v1
//rds service which is a service of db instances management.
func NewRdsServiceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	newsc, err := initClientOpts(client, eo, "compute")
	rdsendpoint := strings.Replace(strings.Replace(newsc.Endpoint, "ecs", "rds", 1), "/v2/", "/rds/v1/", 1)
	newsc.Endpoint = rdsendpoint
	newsc.ResourceBase = rdsendpoint
	newsc.Type = "rds"
	return newsc, err
}

func NewCESClient(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "volumev2")
	if err != nil {
		return nil, err
	}
	e := strings.Replace(sc.Endpoint, "v2", "V1.0", 1)
	sc.Endpoint = strings.Replace(e, "evs", "ces", 1)
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// NewDRSServiceV2 creates a ServiceClient that may be used to access the v2 Data Replication Service.
func NewDRSServiceV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "volumev2")
	return sc, err
}

func NewComputeV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "ecs", 1)
	sc.Endpoint = sc.Endpoint + "v1/"
	sc.ResourceBase = sc.Endpoint + client.ProjectID + "/"
	return sc, err
}

func NewComputeV11(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "ecsv1.1")
	return sc, err
}

func NewEcsV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "ecs")
	return sc, err
}

func NewRdsTagV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "rds", 1)
	sc.Endpoint = sc.Endpoint + "v1/"
	sc.ResourceBase = sc.Endpoint + client.ProjectID + "/rds/"
	return sc, err
}

//NewAutoScalingService creates a ServiceClient that may be used to access the
//auto-scaling service of huawei public cloud
func NewAutoScalingService(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "volumev2")
	if err != nil {
		return nil, err
	}
	e := strings.Replace(sc.Endpoint, "v2", "autoscaling-api/v1", 1)
	sc.Endpoint = strings.Replace(e, "evs", "as", 1)
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// NewAutoScalingV1 creates a ServiceClient that may be used to access the AS service
func NewAutoScalingV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "asv1")
	return sc, err
}

// NewKmsKeyV1 creates a ServiceClient that may be used to access the v1
// kms key service.
func NewKmsKeyV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "compute")
	sc.Endpoint = strings.Replace(sc.Endpoint, "ecs", "kms", 1)
	sc.Endpoint = sc.Endpoint[:strings.LastIndex(sc.Endpoint, "v2")+3]
	sc.Endpoint = strings.Replace(sc.Endpoint, "v2", "v1.0", 1)
	sc.ResourceBase = sc.Endpoint
	sc.Type = "kms"
	return sc, err
}

func NewElasticLoadBalancer(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	if err != nil {
		return sc, err
	}
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "elb", 1)
	sc.Endpoint = strings.Replace(sc.Endpoint, "myhwclouds", "myhuaweicloud", 1)
	sc.ResourceBase = sc.Endpoint + "v1.0/"
	return sc, err
}

// NewNetworkV1 creates a ServiceClient that may be used with the v1 network
// package.
func NewNetworkV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.ResourceBase = sc.Endpoint + "v1/"
	return sc, err
}

// NewNatV2 creates a ServiceClient that may be used with the v2 nat package.
func NewNatV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "nat", 1)
	sc.Endpoint = strings.Replace(sc.Endpoint, "myhwclouds", "myhuaweicloud", 1)
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, err
}

// MapReduceV1 creates a ServiceClient that may be used with the v1 MapReduce service.
func MapReduceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "mrs", 1)
	sc.Endpoint = sc.Endpoint + "v1.1/"
	sc.ResourceBase = sc.Endpoint + client.ProjectID + "/"
	return sc, err
}

// NewMapReduceV1 creates a ServiceClient that may be used with the v1 MapReduce service.
func NewMapReduceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "mrs")
	sc.ResourceBase = sc.Endpoint + client.ProjectID + "/"
	return sc, err
}

// AntiDDoSV1 creates a ServiceClient that may be used with the v1 Anti DDoS service.
func AntiDDoSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "antiddos", 1)
	sc.Endpoint = sc.Endpoint + "v1/"
	sc.ResourceBase = sc.Endpoint + client.ProjectID + "/"
	return sc, err
}

// NewAntiDDoSV1 creates a ServiceClient that may be used with the v1 Anti DDoS Service
// package.
func NewAntiDDoSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "antiddos")
}

// NewAntiDDoSV2 creates a ServiceClient that may be used with the v2 Anti DDoS Service
// package.
func NewAntiDDoSV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "antiddos")
	sc.ResourceBase = sc.Endpoint + "v2/" + client.ProjectID + "/"
	return sc, err
}

func NewCCEV3(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "cce", 1)
	sc.Endpoint = strings.Replace(sc.Endpoint, "myhwclouds", "myhuaweicloud", 1)
	sc.ResourceBase = sc.Endpoint + "api/v3/projects/" + client.ProjectID + "/"
	return sc, err
}

// NewDMSServiceV1 creates a ServiceClient that may be used to access the v1 Distributed Message Service.
func NewDMSServiceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "dms", 1)
	sc.ResourceBase = sc.Endpoint + "v1.0/" + client.ProjectID + "/"
	return sc, err
}

// NewDCSServiceV1 creates a ServiceClient that may be used to access the v1 Distributed Cache Service.
func NewDCSServiceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "dcs", 1)
	sc.ResourceBase = sc.Endpoint + "v1.0/" + client.ProjectID + "/"
	return sc, err
}

// NewOBSService creates a ServiceClient that may be used to access the Object Storage Service.
func NewOBSService(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "object")
	return sc, err
}

//TODO: Need to change to sfs client type from evs once available
//NewSFSV2 creates a service client that is used for Huawei cloud  for SFS , it replaces the EVS type.
func NewHwSFSV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "sfs", 1)
	sc.ResourceBase = sc.Endpoint + "v2/" + client.ProjectID + "/"
	return sc, err
}

func NewBMSV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "compute")
	e := strings.Replace(sc.Endpoint, "v2", "v2.1", 1)
	sc.Endpoint = e
	sc.ResourceBase = e
	return sc, err
}

// NewDeHServiceV1 creates a ServiceClient that may be used to access the v1 Dedicated Hosts service.
func NewDeHServiceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "deh")
	return sc, err
}

// NewCSBSService creates a ServiceClient that can be used to access the Cloud Server Backup service.
func NewCSBSService(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "data-protect")
	return sc, err
}

// NewHwCSBSServiceV1 creates a ServiceClient that may be used to access the Huawei Cloud Server Backup service.
func NewHwCSBSServiceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "compute")
	sc.Endpoint = strings.Replace(sc.Endpoint, "ecs", "csbs", 1)
	e := strings.Replace(sc.Endpoint, "v2", "v1", 1)
	sc.Endpoint = e
	sc.ResourceBase = e
	return sc, err
}

func NewMLSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "mls", 1)
	sc.ResourceBase = sc.Endpoint + "v1.0/" + client.ProjectID + "/"
	return sc, err
}

func NewDWSClient(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "volumev2")
	if err != nil {
		return nil, err
	}
	e := strings.Replace(sc.Endpoint, "v2", "v1.0", 1)
	sc.Endpoint = strings.Replace(e, "evs", "dws", 1)
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// NewVBSV2 creates a ServiceClient that may be used to access the VBS service for Orange and Telefonica Cloud.
func NewVBSV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "vbsv2")
	return sc, err
}

// NewVBS creates a service client that is used for VBS.
func NewVBS(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "volumev2")
	if err != nil {
		return nil, err
	}
	sc.Endpoint = strings.Replace(sc.Endpoint, "evs", "vbs", 1)
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// NewMAASV1 creates a ServiceClient that may be used to access the MAAS service.
func NewMAASV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "maasv1")
	return sc, err
}

// MAASV1 creates a ServiceClient that may be used with the v1 MAAS service.
func MAASV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = "https://oms.myhuaweicloud.com/v1/"
	sc.ResourceBase = sc.Endpoint + client.ProjectID + "/"
	return sc, err
}

func NewHwAntiDDoSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "volumev2")
	if err != nil {
		return nil, err
	}
	e := strings.Replace(sc.Endpoint, "v2", "v1", 1)
	sc.Endpoint = strings.Replace(e, "evs", "antiddos", 1)
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// NewCTSService creates a ServiceClient that can be used to access the Cloud Trace service.
func NewCTSService(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "cts")
	return sc, err
}

// NewELBV1 creates a ServiceClient that may be used to access the ELB service.
func NewELBV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "elbv1")
	return sc, err
}

// NewRDSV1 creates a ServiceClient that may be used to access the RDS service.
func NewRDSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "rdsv1")
	return sc, err
}

// NewKMSV1 creates a ServiceClient that may be used to access the KMS service.
func NewKMSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "kms")
	return sc, err
}

// NewSMNV2 creates a ServiceClient that may be used to access the SMN service.
func NewSMNV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "smnv2")
	sc.ResourceBase = sc.Endpoint + "notifications/"
	return sc, err
}

// NewCCE creates a ServiceClient that may be used to access the CCE service.
func NewCCE(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "ccev2.0")
	sc.ResourceBase = sc.Endpoint + "api/v3/projects/" + client.ProjectID + "/"
	return sc, err
}

// NewWAF creates a ServiceClient that may be used to access the WAF service.
func NewWAFV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "waf")
	sc.ResourceBase = sc.Endpoint + "v1/" + client.ProjectID + "/waf/"
	return sc, err
}

// NewRDSV3 creates a ServiceClient that may be used to access the RDS service.
func NewRDSV3(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "rdsv3")
	return sc, err
}

// SDRSV1 creates a ServiceClient that may be used with the v1 SDRS service.
func SDRSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "sdrs", 1)
	sc.Endpoint = sc.Endpoint + "v1/" + client.ProjectID + "/"
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// TMSV1 creates a ServiceClient that may be used with the v1 TMS service.
func TMSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = "https://tms.myhuaweicloud.com/v1.0/"
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// NewSDRSV1 creates a ServiceClient that may be used to access the SDRS service.
func NewSDRSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "sdrs")
	return sc, err
}

// NewBSSV1 creates a ServiceClient that may be used to access the BSS service.
func NewBSSV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "bssv1")
	return sc, err
}

func NewSDKClient(c *golangsdk.ProviderClient, eo golangsdk.EndpointOpts, serviceType string) (*golangsdk.ServiceClient, error) {
	switch serviceType {
	case "mls":
		return NewMLSV1(c, eo)
	case "dws":
		return NewDWSClient(c, eo)
	case "nat":
		return NewNatV2(c, eo)
	}

	return initClientOpts(c, eo, serviceType)
}

// NewCESV1 creates a ServiceClient that may be used with the v1 CES service.
func NewCESV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "cesv1")
	return sc, err
}

// NewDDSV3 creates a ServiceClient that may be used to access the DDS service.
func NewDDSV3(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "ddsv3")
	return sc, err
}

// NewLTSV2 creates a ServiceClient that may be used to access the LTS service.
func NewLTSV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initcommonServiceClient(client, eo, "lts", "v2.0")
	return sc, err
}

// NewFGSV2 creates a ServiceClient that may be used with the v2 as
// package.
func NewFGSV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "fgsv2")
	return sc, err
}
