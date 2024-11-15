package linodego

import (
	"context"
	"encoding/json"
	"errors"
)

// Tag represents a Tag object
type Tag struct {
	Label string `json:"label"`
}

// TaggedObject represents a Tagged Object object
type TaggedObject struct {
	Type    string          `json:"type"`
	RawData json.RawMessage `json:"data"`
	Data    any             `json:"-"`
}

// SortedObjects currently only includes Instances
type SortedObjects struct {
	Instances     []Instance
	LKEClusters   []LKECluster
	Domains       []Domain
	Volumes       []Volume
	NodeBalancers []NodeBalancer
	/*
		StackScripts  []Stackscript
	*/
}

// TaggedObjectList are a list of TaggedObjects, as returning by ListTaggedObjects
type TaggedObjectList []TaggedObject

// TagCreateOptions fields are those accepted by CreateTag
type TagCreateOptions struct {
	Label   string `json:"label"`
	Linodes []int  `json:"linodes,omitempty"`
	// @TODO is this implemented?
	LKEClusters   []int `json:"lke_clusters,omitempty"`
	Domains       []int `json:"domains,omitempty"`
	Volumes       []int `json:"volumes,omitempty"`
	NodeBalancers []int `json:"nodebalancers,omitempty"`
}

// GetCreateOptions converts a Tag to TagCreateOptions for use in CreateTag
func (i Tag) GetCreateOptions() (o TagCreateOptions) {
	o.Label = i.Label
	return
}

// ListTags lists Tags
func (c *Client) ListTags(ctx context.Context, opts *ListOptions) ([]Tag, error) {
	response, err := getPaginatedResults[Tag](ctx, c, "tags", opts)
	return response, err
}

// fixData stores an object of the type defined by Type in Data using RawData
func (i *TaggedObject) fixData() (*TaggedObject, error) {
	switch i.Type {
	case "linode":
		obj := Instance{}
		if err := json.Unmarshal(i.RawData, &obj); err != nil {
			return nil, err
		}
		i.Data = obj
	case "lke_cluster":
		obj := LKECluster{}
		if err := json.Unmarshal(i.RawData, &obj); err != nil {
			return nil, err
		}
		i.Data = obj
	case "nodebalancer":
		obj := NodeBalancer{}
		if err := json.Unmarshal(i.RawData, &obj); err != nil {
			return nil, err
		}
		i.Data = obj
	case "domain":
		obj := Domain{}
		if err := json.Unmarshal(i.RawData, &obj); err != nil {
			return nil, err
		}
		i.Data = obj
	case "volume":
		obj := Volume{}
		if err := json.Unmarshal(i.RawData, &obj); err != nil {
			return nil, err
		}
		i.Data = obj
	}

	return i, nil
}

// ListTaggedObjects lists Tagged Objects
func (c *Client) ListTaggedObjects(ctx context.Context, label string, opts *ListOptions) (TaggedObjectList, error) {
	response, err := getPaginatedResults[TaggedObject](ctx, c, formatAPIPath("tags/%s", label), opts)
	if err != nil {
		return nil, err
	}

	for i := range response {
		if _, err := response[i].fixData(); err != nil {
			return nil, err
		}
	}
	return response, nil
}

// SortedObjects converts a list of TaggedObjects into a Sorted Objects struct, for easier access
func (t TaggedObjectList) SortedObjects() (SortedObjects, error) {
	so := SortedObjects{}

	for _, o := range t {
		switch o.Type {
		case "linode":
			if instance, ok := o.Data.(Instance); ok {
				so.Instances = append(so.Instances, instance)
			} else {
				return so, errors.New("expected an Instance when Type was \"linode\"")
			}
		case "lke_cluster":
			if lkeCluster, ok := o.Data.(LKECluster); ok {
				so.LKEClusters = append(so.LKEClusters, lkeCluster)
			} else {
				return so, errors.New("expected an LKECluster when Type was \"lke_cluster\"")
			}
		case "domain":
			if domain, ok := o.Data.(Domain); ok {
				so.Domains = append(so.Domains, domain)
			} else {
				return so, errors.New("expected a Domain when Type was \"domain\"")
			}
		case "volume":
			if volume, ok := o.Data.(Volume); ok {
				so.Volumes = append(so.Volumes, volume)
			} else {
				return so, errors.New("expected an Volume when Type was \"volume\"")
			}
		case "nodebalancer":
			if nodebalancer, ok := o.Data.(NodeBalancer); ok {
				so.NodeBalancers = append(so.NodeBalancers, nodebalancer)
			} else {
				return so, errors.New("expected an NodeBalancer when Type was \"nodebalancer\"")
			}
		}
	}
	return so, nil
}

// CreateTag creates a Tag
func (c *Client) CreateTag(ctx context.Context, opts TagCreateOptions) (*Tag, error) {
	e := "tags"
	response, err := doPOSTRequest[Tag](ctx, c, e, opts)
	return response, err
}

// DeleteTag deletes the Tag with the specified id
func (c *Client) DeleteTag(ctx context.Context, label string) error {
	e := formatAPIPath("tags/%s", label)
	err := doDELETERequest(ctx, c, e)
	return err
}
