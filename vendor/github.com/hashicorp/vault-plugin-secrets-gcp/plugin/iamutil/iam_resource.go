package iamutil

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"google.golang.org/api/googleapi"
	"io"
	"net/http"
	"strings"
)

// IamResource handles constructing HTTP requests for getting and
// setting IAM policies.
type IamResource interface {
	GetIamPolicyRequest() (*http.Request, error)
	SetIamPolicyRequest(*Policy) (req *http.Request, err error)
}

// parsedIamResource implements IamResource.
type parsedIamResource struct {
	relativeId *gcputil.RelativeResourceName
	config     *IamRestResource
}

type IamRestResource struct {
	// Name is the base name of the resource
	// i.e. for a GCE instance: "instance"
	Name string

	// Type Key is the identifying path for the resource, or
	// the RESTful resource identifier without resource IDs
	// i.e. For a GCE instance: "projects/zones/instances"
	TypeKey string

	// Service Information
	// Service is the name of the service this resource belongs to.
	Service string

	// IsPreferredVersion is true if this version of the API/resource is preferred.
	IsPreferredVersion bool

	// IsPreferredVersion is true if this version of the API/resource is preferred.
	GetMethod RestMethod

	// IsPreferredVersion is true if this version of the API/resource is preferred.
	SetMethod RestMethod

	// Ordered parameters to be replaced in method paths
	Parameters []string

	// collection Id --> parameter to be replaced {} name
	CollectionReplacementKeys map[string]string
}

type RestMethod struct {
	HttpMethod    string
	BaseURL       string
	Path          string
	RequestFormat string
}

func (r *parsedIamResource) SetIamPolicyRequest(p *Policy) (req *http.Request, err error) {
	jsonP, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	reqJson := fmt.Sprintf(r.config.SetMethod.RequestFormat, jsonP)
	if !json.Valid([]byte(reqJson)) {
		return nil, fmt.Errorf("request format from generated IAM config invalid JSON: %s", reqJson)
	}

	return r.constructRequest(&r.config.SetMethod, strings.NewReader(reqJson))
}

func (r *parsedIamResource) GetIamPolicyRequest() (*http.Request, error) {
	return r.constructRequest(&r.config.GetMethod, nil)
}

func (r *parsedIamResource) constructRequest(restMethod *RestMethod, data io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(
		restMethod.HttpMethod,
		googleapi.ResolveRelative(restMethod.BaseURL, restMethod.Path),
		data)
	if err != nil {
		return nil, err
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	relId := r.relativeId
	replacementMap := make(map[string]string)

	if strings.Contains(restMethod.Path, "{+resource}") {
		// +resource is used to represent full relative resource name
		if len(r.config.Parameters) == 1 && r.config.Parameters[0] == "resource" {
			relName := ""
			tkns := strings.Split(r.config.TypeKey, "/")
			for _, colId := range tkns {
				relName += fmt.Sprintf("%s/%s/", colId, relId.IdTuples[colId])
			}
			replacementMap["resource"] = strings.Trim(relName, "/")
		}
	} else {
		for colId, resId := range relId.IdTuples {
			rId, ok := r.config.CollectionReplacementKeys[colId]
			if !ok {
				return nil, fmt.Errorf("expected value for collection id %s", colId)
			}
			replacementMap[rId] = resId
		}
	}

	googleapi.Expand(req.URL, replacementMap)
	return req, nil
}
