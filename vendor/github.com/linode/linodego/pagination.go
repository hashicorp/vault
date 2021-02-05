package linodego

/**
 * Pagination and Filtering types and helpers
 */

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/linode/linodego/pkg/errors"
)

// PageOptions are the pagination parameters for List endpoints
type PageOptions struct {
	Page    int `url:"page,omitempty" json:"page"`
	Pages   int `url:"pages,omitempty" json:"pages"`
	Results int `url:"results,omitempty" json:"results"`
}

// ListOptions are the pagination and filtering (TODO) parameters for endpoints
type ListOptions struct {
	*PageOptions
	Filter string
}

// NewListOptions simplified construction of ListOptions using only
// the two writable properties, Page and Filter
func NewListOptions(page int, filter string) *ListOptions {
	return &ListOptions{PageOptions: &PageOptions{Page: page}, Filter: filter}
}

// listHelper abstracts fetching and pagination for GET endpoints that
// do not require any Ids (top level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
// nolint
func (c *Client) listHelper(ctx context.Context, i interface{}, opts *ListOptions) error {
	req := c.R(ctx)
	if opts != nil && opts.PageOptions != nil && opts.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(opts.Page))
	}

	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	if opts != nil && len(opts.Filter) > 0 {
		req.SetHeader("X-Filter", opts.Filter)
	}

	switch v := i.(type) {
	case *LinodeKernelsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(LinodeKernelsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*LinodeKernelsPagedResponse).Pages
			results = r.Result().(*LinodeKernelsPagedResponse).Results
			v.appendData(r.Result().(*LinodeKernelsPagedResponse))
		}
	case *LinodeTypesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(LinodeTypesPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*LinodeTypesPagedResponse).Pages
			results = r.Result().(*LinodeTypesPagedResponse).Results
			v.appendData(r.Result().(*LinodeTypesPagedResponse))
		}
	case *ImagesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(ImagesPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*ImagesPagedResponse).Pages
			results = r.Result().(*ImagesPagedResponse).Results
			v.appendData(r.Result().(*ImagesPagedResponse))
		}
	case *StackscriptsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(StackscriptsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*StackscriptsPagedResponse).Pages
			results = r.Result().(*StackscriptsPagedResponse).Results
			v.appendData(r.Result().(*StackscriptsPagedResponse))
		}
	case *InstancesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(InstancesPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*InstancesPagedResponse).Pages
			results = r.Result().(*InstancesPagedResponse).Results
			v.appendData(r.Result().(*InstancesPagedResponse))
		}
	case *RegionsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(RegionsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*RegionsPagedResponse).Pages
			results = r.Result().(*RegionsPagedResponse).Results
			v.appendData(r.Result().(*RegionsPagedResponse))
		}
	case *VolumesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(VolumesPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*VolumesPagedResponse).Pages
			results = r.Result().(*VolumesPagedResponse).Results
			v.appendData(r.Result().(*VolumesPagedResponse))
		}
	case *DomainsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(DomainsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			response, ok := r.Result().(*DomainsPagedResponse)
			if !ok {
				return fmt.Errorf("response is not a *DomainsPagedResponse")
			}
			pages = response.Pages
			results = response.Results
			v.appendData(response)
		}
	case *EventsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(EventsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*EventsPagedResponse).Pages
			results = r.Result().(*EventsPagedResponse).Results
			v.appendData(r.Result().(*EventsPagedResponse))
		}
	case *FirewallsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(FirewallsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*FirewallsPagedResponse).Pages
			results = r.Result().(*FirewallsPagedResponse).Results
			v.appendData(r.Result().(*FirewallsPagedResponse))
		}
	case *LKEClustersPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(LKEClustersPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*LKEClustersPagedResponse).Pages
			results = r.Result().(*LKEClustersPagedResponse).Results
			v.appendData(r.Result().(*LKEClustersPagedResponse))
		}
	case *LKEVersionsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(LKEVersionsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*LKEVersionsPagedResponse).Pages
			results = r.Result().(*LKEVersionsPagedResponse).Results
			v.appendData(r.Result().(*LKEVersionsPagedResponse))
		}
	case *LongviewSubscriptionsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(LongviewSubscriptionsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*LongviewSubscriptionsPagedResponse).Pages
			results = r.Result().(*LongviewSubscriptionsPagedResponse).Results
			v.appendData(r.Result().(*LongviewSubscriptionsPagedResponse))
		}
	case *LongviewClientsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(LongviewClientsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*LongviewClientsPagedResponse).Pages
			results = r.Result().(*LongviewClientsPagedResponse).Results
			v.appendData(r.Result().(*LongviewClientsPagedResponse))
		}
	case *IPAddressesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(IPAddressesPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*IPAddressesPagedResponse).Pages
			results = r.Result().(*IPAddressesPagedResponse).Results
			v.appendData(r.Result().(*IPAddressesPagedResponse))
		}
	case *IPv6PoolsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(IPv6PoolsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*IPv6PoolsPagedResponse).Pages
			results = r.Result().(*IPv6PoolsPagedResponse).Results
			v.appendData(r.Result().(*IPv6PoolsPagedResponse))
		}
	case *IPv6RangesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(IPv6RangesPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*IPv6RangesPagedResponse).Pages
			results = r.Result().(*IPv6RangesPagedResponse).Results
			v.appendData(r.Result().(*IPv6RangesPagedResponse))
			// @TODO consolidate this type with IPv6PoolsPagedResponse?
		}
	case *SSHKeysPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(SSHKeysPagedResponse{}).Get(v.endpoint(c))); err == nil {
			response, ok := r.Result().(*SSHKeysPagedResponse)
			if !ok {
				return fmt.Errorf("response is not a *SSHKeysPagedResponse")
			}
			pages = response.Pages
			results = response.Results
			v.appendData(response)
		}
	case *TicketsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(TicketsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*TicketsPagedResponse).Pages
			results = r.Result().(*TicketsPagedResponse).Results
			v.appendData(r.Result().(*TicketsPagedResponse))
		}
	case *InvoicesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(InvoicesPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*InvoicesPagedResponse).Pages
			results = r.Result().(*InvoicesPagedResponse).Results
			v.appendData(r.Result().(*InvoicesPagedResponse))
		}
	case *NotificationsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(NotificationsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*NotificationsPagedResponse).Pages
			results = r.Result().(*NotificationsPagedResponse).Results
			v.appendData(r.Result().(*NotificationsPagedResponse))
		}
	case *OAuthClientsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(OAuthClientsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*OAuthClientsPagedResponse).Pages
			results = r.Result().(*OAuthClientsPagedResponse).Results
			v.appendData(r.Result().(*OAuthClientsPagedResponse))
		}
	case *PaymentsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(PaymentsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*PaymentsPagedResponse).Pages
			results = r.Result().(*PaymentsPagedResponse).Results
			v.appendData(r.Result().(*PaymentsPagedResponse))
		}
	case *NodeBalancersPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(NodeBalancersPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*NodeBalancersPagedResponse).Pages
			results = r.Result().(*NodeBalancersPagedResponse).Results
			v.appendData(r.Result().(*NodeBalancersPagedResponse))
		}
	case *TagsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(TagsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*TagsPagedResponse).Pages
			results = r.Result().(*TagsPagedResponse).Results
			v.appendData(r.Result().(*TagsPagedResponse))
		}
	case *TokensPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(TokensPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*TokensPagedResponse).Pages
			results = r.Result().(*TokensPagedResponse).Results
			v.appendData(r.Result().(*TokensPagedResponse))
		}
	case *UsersPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(UsersPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*UsersPagedResponse).Pages
			results = r.Result().(*UsersPagedResponse).Results
			v.appendData(r.Result().(*UsersPagedResponse))
		}
	case *ObjectStorageBucketsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(ObjectStorageBucketsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*ObjectStorageBucketsPagedResponse).Pages
			results = r.Result().(*ObjectStorageBucketsPagedResponse).Results
			v.appendData(r.Result().(*ObjectStorageBucketsPagedResponse))
		}
	case *ObjectStorageClustersPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(ObjectStorageClustersPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*ObjectStorageClustersPagedResponse).Pages
			results = r.Result().(*ObjectStorageClustersPagedResponse).Results
			v.appendData(r.Result().(*ObjectStorageClustersPagedResponse))
		}
	case *ObjectStorageKeysPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(ObjectStorageKeysPagedResponse{}).Get(v.endpoint(c))); err == nil {
			pages = r.Result().(*ObjectStorageKeysPagedResponse).Pages
			results = r.Result().(*ObjectStorageKeysPagedResponse).Results
			v.appendData(r.Result().(*ObjectStorageKeysPagedResponse))
		}
	/**
	case ProfileAppsPagedResponse:
	case ProfileWhitelistPagedResponse:
	case ManagedContactsPagedResponse:
	case ManagedCredentialsPagedResponse:
	case ManagedIssuesPagedResponse:
	case ManagedLinodeSettingsPagedResponse:
	case ManagedServicesPagedResponse:
	**/
	default:
		log.Fatalf("listHelper interface{} %+v used", i)
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page++ {
			if err := c.listHelper(ctx, i, &ListOptions{PageOptions: &PageOptions{Page: page}}); err != nil {
				return err
			}
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}

		if opts.Page == 0 {
			for page := 2; page <= pages; page++ {
				opts.Page = page
				if err := c.listHelper(ctx, i, opts); err != nil {
					return err
				}
			}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}

// listHelperWithID abstracts fetching and pagination for GET endpoints that
// require an Id (second level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
// nolint
func (c *Client) listHelperWithID(ctx context.Context, i interface{}, idRaw interface{}, opts *ListOptions) error {
	req := c.R(ctx)
	if opts != nil && opts.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(opts.Page))
	}

	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	id, _ := idRaw.(int)

	if opts != nil && len(opts.Filter) > 0 {
		req.SetHeader("X-Filter", opts.Filter)
	}

	switch v := i.(type) {
	case *DomainRecordsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(DomainRecordsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			response, ok := r.Result().(*DomainRecordsPagedResponse)
			if !ok {
				return fmt.Errorf("response is not a *DomainRecordsPagedResponse")
			}
			pages = response.Pages
			results = response.Results
			v.appendData(response)
		}
	case *FirewallDevicesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(FirewallDevicesPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*FirewallDevicesPagedResponse).Pages
			results = r.Result().(*FirewallDevicesPagedResponse).Results
			v.appendData(r.Result().(*FirewallDevicesPagedResponse))
		}
	case *InstanceConfigsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(InstanceConfigsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*InstanceConfigsPagedResponse).Pages
			results = r.Result().(*InstanceConfigsPagedResponse).Results
			v.appendData(r.Result().(*InstanceConfigsPagedResponse))
		}
	case *InstanceDisksPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(InstanceDisksPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*InstanceDisksPagedResponse).Pages
			results = r.Result().(*InstanceDisksPagedResponse).Results
			v.appendData(r.Result().(*InstanceDisksPagedResponse))
		}
	case *InstanceVolumesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(InstanceVolumesPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*InstanceVolumesPagedResponse).Pages
			results = r.Result().(*InstanceVolumesPagedResponse).Results
			v.appendData(r.Result().(*InstanceVolumesPagedResponse))
		}
	case *InvoiceItemsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(InvoiceItemsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*InvoiceItemsPagedResponse).Pages
			results = r.Result().(*InvoiceItemsPagedResponse).Results
			v.appendData(r.Result().(*InvoiceItemsPagedResponse))
		}
	case *LKEClusterAPIEndpointsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(LKEClusterAPIEndpointsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*LKEClusterAPIEndpointsPagedResponse).Pages
			results = r.Result().(*LKEClusterAPIEndpointsPagedResponse).Results
			v.appendData(r.Result().(*LKEClusterAPIEndpointsPagedResponse))
		}
	case *LKEClusterPoolsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(LKEClusterPoolsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*LKEClusterPoolsPagedResponse).Pages
			results = r.Result().(*LKEClusterPoolsPagedResponse).Results
			v.appendData(r.Result().(*LKEClusterPoolsPagedResponse))
		}
	case *NodeBalancerConfigsPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(NodeBalancerConfigsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			pages = r.Result().(*NodeBalancerConfigsPagedResponse).Pages
			results = r.Result().(*NodeBalancerConfigsPagedResponse).Results
			v.appendData(r.Result().(*NodeBalancerConfigsPagedResponse))
		}
	case *TaggedObjectsPagedResponse:
		idStr := idRaw.(string)

		if r, err = errors.CoupleAPIErrors(req.SetResult(TaggedObjectsPagedResponse{}).Get(v.endpointWithID(c, idStr))); err == nil {
			pages = r.Result().(*TaggedObjectsPagedResponse).Pages
			results = r.Result().(*TaggedObjectsPagedResponse).Results
			v.appendData(r.Result().(*TaggedObjectsPagedResponse))
		}
	/**
	case TicketAttachmentsPagedResponse:
		if r, err = req.SetResult(v).Get(v.endpoint(c)); r.Error() != nil {
			return errors.New(r)
		} else if err == nil {
			pages = r.Result().(*TicketAttachmentsPagedResponse).Pages
			results = r.Result().(*TicketAttachmentsPagedResponse).Results
			v.appendData(r.Result().(*TicketAttachmentsPagedResponse))
		}
	case TicketRepliesPagedResponse:
		if r, err = req.SetResult(v).Get(v.endpoint(c)); r.Error() != nil {
			return errors.New(r)
		} else if err == nil {
			pages = r.Result().(*TicketRepliesPagedResponse).Pages
			results = r.Result().(*TicketRepliesPagedResponse).Results
			v.appendData(r.Result().(*TicketRepliesPagedResponse))
		}
	**/
	default:
		log.Fatalf("Unknown listHelperWithID interface{} %T used", i)
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page++ {
			if err := c.listHelperWithID(ctx, i, id, &ListOptions{PageOptions: &PageOptions{Page: page}}); err != nil {
				return err
			}
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}
		if opts.Page == 0 {
			for page := 2; page <= pages; page++ {
				opts.Page = page
				if err := c.listHelperWithID(ctx, i, id, opts); err != nil {
					return err
				}
			}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}

// listHelperWithTwoIDs abstracts fetching and pagination for GET endpoints that
// require twos IDs (third level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
func (c *Client) listHelperWithTwoIDs(ctx context.Context, i interface{}, firstID, secondID int, opts *ListOptions) error {
	req := c.R(ctx)

	if opts != nil && opts.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(opts.Page))
	}

	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	if opts != nil && len(opts.Filter) > 0 {
		req.SetHeader("X-Filter", opts.Filter)
	}

	switch v := i.(type) {
	case *NodeBalancerNodesPagedResponse:
		if r, err = errors.CoupleAPIErrors(req.SetResult(NodeBalancerNodesPagedResponse{}).Get(v.endpointWithTwoIDs(c, firstID, secondID))); err == nil {
			pages = r.Result().(*NodeBalancerNodesPagedResponse).Pages
			results = r.Result().(*NodeBalancerNodesPagedResponse).Results
			v.appendData(r.Result().(*NodeBalancerNodesPagedResponse))
		}
	default:
		log.Fatalf("Unknown listHelperWithTwoIDs interface{} %T used", i)
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page++ {
			if err := c.listHelper(ctx, i, &ListOptions{PageOptions: &PageOptions{Page: page}}); err != nil {
				return err
			}
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}
		if opts.Page == 0 {
			for page := 2; page <= pages; page++ {
				opts.Page = page
				if err := c.listHelperWithTwoIDs(ctx, i, firstID, secondID, opts); err != nil {
					return err
				}
			}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}
