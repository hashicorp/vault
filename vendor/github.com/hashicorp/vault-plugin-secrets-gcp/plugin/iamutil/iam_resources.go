//go:generate go run internal/generate_iam_resources.go

package iamutil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"google.golang.org/api/googleapi"
)

const (
	resourceParsingErrorTmpl     = `invalid resource "%s": %v`
	resourceMultipleServicesTmpl = `please provide a self-link or full resource name for non-service-unique resource type '%s' (supported services: %s)`
	resourceMultipleVersions     = `please provide a self-link with version instead; IAM support for this resource is for multiple non-preferred service versions`
)

type EnabledResources interface {
	Resource(resource string) (IamResource, error)
}

type iamResourceMap map[string]map[string]map[string]*IamResourceConfig

type generatedIamResources struct {
	resources iamResourceMap
}

func (apis *generatedIamResources) parseResource(name string) (*gcputil.RelativeResourceName, *IamResourceConfig, error) {
	rUrl, err := url.Parse(name)
	if err != nil {
		return nil, nil, fmt.Errorf(`resource "%s" is invalid URI`, name)
	}

	var relName *gcputil.RelativeResourceName
	var hasServiceVersion bool
	var serviceName string

	if rUrl.Scheme != "" {
		selfLink, err := gcputil.ParseProjectResourceSelfLink(name)
		if err != nil {
			return nil, nil, err
		}
		hasServiceVersion = true
		relName = selfLink.RelativeResourceName
	} else if rUrl.Host != "" {
		fullName, err := gcputil.ParseFullResourceName(name)
		if err != nil {
			return nil, nil, err
		}
		relName = fullName.RelativeResourceName
		serviceName = fullName.Service
	} else {
		relName, err = gcputil.ParseRelativeName(name)
		if err != nil {
			return nil, nil, err
		}
	}

	if relName == nil {
		return nil, nil, fmt.Errorf(resourceParsingErrorTmpl, name, "unable to parse relative name")
	}

	serviceMap, ok := apis.resources[relName.TypeKey]
	if !ok {
		return nil, nil, fmt.Errorf(resourceParsingErrorTmpl, name, fmt.Errorf("unsupported resource type: %s", relName.TypeKey))
	}

	var resConfig *IamResourceConfig
	if hasServiceVersion {
		resConfig, err = tryGetConfigForSelfLink(name, relName.TypeKey, serviceMap)
	} else {
		resConfig, err = tryGetUniqueVersion(serviceName, relName.TypeKey, serviceMap)
	}
	if err != nil {
		return nil, nil, err
	}
	if resConfig == nil {
		return nil, nil, fmt.Errorf(resourceParsingErrorTmpl, name, "unable to get IAM resource config")
	}

	return relName, resConfig, nil
}

type IamResource interface {
	GetIamPolicyRequest() (*http.Request, error)
	SetIamPolicyRequest(*Policy) (*http.Request, error)
}

func (apis *generatedIamResources) Resource(name string) (IamResource, error) {
	relName, cfg, err := apis.parseResource(name)
	if err != nil {
		return nil, err
	}

	return &iamResourceImpl{
		relativeId: relName,
		config:     cfg,
	}, nil
}

func tryGetConfigForSelfLink(link, typeKey string, resourceServices map[string]map[string]*IamResourceConfig) (*IamResourceConfig, error) {
	for _, verMap := range resourceServices {
		for _, resourceCfg := range verMap {
			prefix := resourceCfg.Service.RootUrl + resourceCfg.Service.ServicePath
			if strings.HasPrefix(link, prefix) {
				return resourceCfg, nil
			}
		}
	}
	return nil, fmt.Errorf("could not find service/version given in self-link for resource type %s", typeKey)
}

func tryGetUniqueVersion(serviceName, typeKey string, resourceServices map[string]map[string]*IamResourceConfig) (*IamResourceConfig, error) {
	if serviceName == "" {
		return tryGetUniqueServiceAndVersion(typeKey, resourceServices)
	}
	if resourceServices == nil {
		return nil, fmt.Errorf("no supported services for %s", typeKey)
	}
	verMap, hasService := resourceServices[serviceName]
	if !hasService {
		return nil, fmt.Errorf("unsupported service '%s' for resource type: %s", serviceName, typeKey)
	}
	return getResourceFromVersions(verMap)
}

func tryGetUniqueServiceAndVersion(typeKey string, resourceServices map[string]map[string]*IamResourceConfig) (*IamResourceConfig, error) {
	if resourceServices == nil || len(resourceServices) < 1 {
		return nil, fmt.Errorf("no supported services for %s", typeKey)
	}
	isUnique := len(resourceServices) == 1
	supported := ""
	for serviceName, verMap := range resourceServices {
		supported += serviceName + ", "
		if isUnique {
			return getResourceFromVersions(verMap)
		}
	}

	return nil, fmt.Errorf(resourceMultipleServicesTmpl, typeKey, strings.Trim(supported, ", "))
}

func getResourceFromVersions(versionsMap map[string]*IamResourceConfig) (*IamResourceConfig, error) {
	var preferredVer *IamResourceConfig
	var onlyCfg *IamResourceConfig

	for _, onlyCfg = range versionsMap {
		if onlyCfg.Service.IsPreferredVersion {
			preferredVer = onlyCfg
			break
		}
	}

	if preferredVer != nil {
		return preferredVer, nil
	} else if len(versionsMap) == 1 {
		return onlyCfg, nil
	} else {
		return nil, errors.New(resourceMultipleVersions)
	}
}

type iamResourceImpl struct {
	relativeId *gcputil.RelativeResourceName
	config     *IamResourceConfig
}

type IamResourceConfig struct {
	// Service this resource belongs to
	Service *ServiceConfig

	// Config for IAM Methods
	SetIamPolicy *HttpMethodCfg
	GetIamPolicy *HttpMethodCfg
}

type HttpMethodCfg struct {
	// HTTP method, e.g. GET/PUT/POST
	HttpMethod string `json:"httpMethod"`

	// Path is the API method's path with replacement keys, e.g.
	// v1/projects/{project}:getIamPolicy
	Path string `json:"flatPath"`

	// ReplacementKeys maps collectionIds in the expected resource format to the key for googleapis.Expand
	// For example, given input of:
	// 		Resource: "projects/my-project/zones/my-zone/instances/someInstance"
	// 		Method Path: "p/{projectId}/z/{zoneId}/i/{resource}"
	//
	// This would be:
	// 		map[string]string{
	//			"projects": "projectId" ,
	//			"zones": "zoneId",
	//			"instances": "resource"
	// 		}
	ReplacementKeys map[string]string
}

type ServiceConfig struct {
	// API service Name (e.g. "compute", "iam", "pubsub")
	Name string

	// API service Version.
	Version string

	// IsPreferredVersion is
	IsPreferredVersion bool

	// Root URL + Service Path is the prefix for all calls using this service.
	RootUrl     string
	ServicePath string
}

func (r *iamResourceImpl) SetIamPolicyRequest(p *Policy) (*http.Request, error) {
	data := struct {
		Policy *Policy `json:"policy,omitempty"`
	}{Policy: p}

	buf, err := googleapi.WithoutDataWrapper.JSONReader(data)
	if err != nil {
		return nil, err
	}

	return r.constructRequest(r.config.SetIamPolicy, buf)
}

func (r *iamResourceImpl) GetIamPolicyRequest() (*http.Request, error) {
	return r.constructRequest(r.config.GetIamPolicy, nil)
}

func (r *iamResourceImpl) constructRequest(httpMtd *HttpMethodCfg, data io.Reader) (*http.Request, error) {
	reqUrl := googleapi.ResolveRelative(r.config.Service.RootUrl+r.config.Service.ServicePath, httpMtd.Path)
	req, err := http.NewRequest(httpMtd.HttpMethod, reqUrl, data)
	if err != nil {
		return nil, err
	}

	replacementMap := make(map[string]string)
	for cId, replaceK := range httpMtd.ReplacementKeys {
		rId, ok := r.relativeId.IdTuples[cId]
		if !ok {
			return nil, fmt.Errorf("expected value for collection id %s", cId)
		}
		replacementMap[replaceK] = rId
	}

	googleapi.Expand(req.URL, replacementMap)
	return req, nil
}
