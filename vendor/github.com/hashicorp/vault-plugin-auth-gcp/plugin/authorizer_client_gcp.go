package gcpauth

import (
	"context"
	"fmt"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
)

var _ client = (*gcpClient)(nil)

// gcpClient implements client and communicates with the GCP API. It is
// abstracted as an interface for stubbing during testing. See stubbedClient for
// more details.
type gcpClient struct {
	logger     log.Logger
	computeSvc *compute.Service
	iamSvc     *iam.Service
}

func (c *gcpClient) InstanceGroups(ctx context.Context, project string, boundInstanceGroups []string) (map[string][]string, error) {
	// map of zone names to a slice of instance group names in that zone.
	igz := make(map[string][]string)

	if err := c.computeSvc.InstanceGroups.
		AggregatedList(project).
		Fields("items/*/instanceGroups/name").
		Pages(ctx, func(l *compute.InstanceGroupAggregatedList) error {
			for k, v := range l.Items {
				// Some groups returned are regional
				// TODO(emilymye, #73): Support regions?
				if strings.Contains(k, "/regions/") {
					c.logger.Debug("ignoring instance groups under region in instance group aggregated list", "key", k)
					continue
				}

				zone, err := zoneFromSelfLink(k)
				if err != nil {
					return err
				}

				for _, g := range v.InstanceGroups {
					if strutil.StrListContains(boundInstanceGroups, g.Name) {
						igz[zone] = append(igz[zone], g.Name)
					}
				}
			}
			return nil
		}); err != nil {
		return nil, err
	}

	return igz, nil
}

func (c *gcpClient) InstanceGroupContainsInstance(ctx context.Context, project, zone, group, instanceSelfLink string) (bool, error) {
	var req compute.InstanceGroupsListInstancesRequest
	resp, err := c.computeSvc.InstanceGroups.
		ListInstances(project, zone, group, &req).
		Filter(fmt.Sprintf("instance eq %s", instanceSelfLink)).
		Context(ctx).
		Do()
	if err != nil {
		return false, err
	}

	if resp != nil && len(resp.Items) > 0 {
		return true, nil
	}
	return false, nil
}

func (c *gcpClient) ServiceAccount(ctx context.Context, name string) (string, string, error) {
	account, err := c.iamSvc.Projects.ServiceAccounts.
		Get(name).
		Fields("uniqueId", "email").
		Context(ctx).
		Do()
	if err != nil {
		return "", "", err
	}

	return account.UniqueId, account.Email, nil
}
