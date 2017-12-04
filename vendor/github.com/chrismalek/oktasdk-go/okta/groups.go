package okta

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

const (
	// GroupTypeOKTA - group type constant for an OKTA Mastered Group
	GroupTypeOKTA = "OKTA_GROUP"
	// GroupTypeBuiltIn - group type constant for a Built in OKTA groups
	GroupTypeBuiltIn = "BUILT_IN"
	// GroupTypeApp -- group type constant for app mastered group
	GroupTypeApp = "APP_GROUP"

	groupTypeFilter                  = "type"
	groupNameFilter                  = "q"
	groupLastMembershipUpdatedFilter = "lastMembershipUpdated"
	groupLastUpdatedFilter           = "lastUpdated"
)

// GroupsService handles communication with the Groups data related
// methods of the OKTA API.
type GroupsService service

// Group represents the Group Object from the OKTA API
type Group struct {
	ID                    string    `json:"id"`
	Created               time.Time `json:"created"`
	LastUpdated           time.Time `json:"lastUpdated"`
	LastMembershipUpdated time.Time `json:"lastMembershipUpdated"`
	ObjectClass           []string  `json:"objectClass"`
	Type                  string    `json:"type"`
	Profile               struct {
		Name                       string `json:"name"`
		Description                string `json:"description"`
		SamAccountName             string `json:"samAccountName"`
		Dn                         string `json:"dn"`
		WindowsDomainQualifiedName string `json:"windowsDomainQualifiedName"`
		ExternalID                 string `json:"externalId"`
	} `json:"profile"`
	Links struct {
		Logo []struct {
			Name string `json:"name"`
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"logo"`
		Users struct {
			Href string `json:"href"`
		} `json:"users"`
		Apps struct {
			Href string `json:"href"`
		} `json:"apps"`
	} `json:"_links"`
}

// GroupFilterOptions is used to generate a "Filter" to search for different groups
// The values here coorelate to API Search paramgters on the group API
type GroupFilterOptions struct {
	// This will be built by internal - may not need to export
	FilterString  string   `url:"filter,omitempty"`
	NextURL       *url.URL `url:"-"`
	GetAllPages   bool     `url:"-"`
	NumberOfPages int      `url:"-"`
	Limit         int      `url:"limit,omitempty"`

	NameStartsWith string `url:"q,omitempty"`
	GroupTypeEqual string `url:"-"`

	LastUpdated           dateFilter `url:"-"`
	LastMembershipUpdated dateFilter `url:"-"`
}

func (g Group) String() string {
	// return Stringify(g)
	return fmt.Sprintf("Group:(ID: {%v} - Type: {%v} - Group Name: {%v})\n", g.ID, g.Type, g.Profile.Name)
}

// ListWithFilter - Method to list groups with different filter options.
//  Pass in a GroupFilterOptions to specify filters. Values in that struct will turn into Query parameters
func (g *GroupsService) ListWithFilter(opt *GroupFilterOptions) ([]Group, *Response, error) {

	var u string
	var err error

	pagesRetreived := 0
	if opt.NextURL != nil {
		u = opt.NextURL.String()
	} else {
		if opt.GroupTypeEqual != "" {
			opt.FilterString = appendToFilterString(opt.FilterString, groupTypeFilter, FilterEqualOperator, opt.GroupTypeEqual)
		}

		// if opt.NameStartsWith != "" {
		// 	opt.FilterString = appendToFilterString(opt.FilterString, groupNameFilter, filterEqualOperator, opt.NameStartsWith)
		// }
		if (!opt.LastMembershipUpdated.Value.IsZero()) && (opt.LastMembershipUpdated.Operator != "") {
			opt.FilterString = appendToFilterString(opt.FilterString, groupLastMembershipUpdatedFilter, opt.LastMembershipUpdated.Operator, opt.LastMembershipUpdated.Value.UTC().Format(oktaFilterTimeFormat))
		}

		if (!opt.LastUpdated.Value.IsZero()) && (opt.LastUpdated.Operator != "") {
			opt.FilterString = appendToFilterString(opt.FilterString, groupLastUpdatedFilter, opt.LastUpdated.Operator, opt.LastUpdated.Value.UTC().Format(oktaFilterTimeFormat))
		}

		if opt.Limit == 0 {
			opt.Limit = defaultLimit
		}
		u, err = addOptions("groups", opt)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := g.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	groups := make([]Group, 1)
	resp, err := g.client.Do(req, &groups)
	if err != nil {
		return nil, resp, err
	}
	pagesRetreived++

	if (opt.NumberOfPages > 0 && pagesRetreived < opt.NumberOfPages) || opt.GetAllPages {

		for {

			if pagesRetreived == opt.NumberOfPages {
				break
			}
			if resp.NextURL != nil {
				var groupPage []Group
				pageOption := new(GroupFilterOptions)
				pageOption.NextURL = resp.NextURL
				pageOption.NumberOfPages = 1
				pageOption.Limit = opt.Limit

				groupPage, resp, err = g.ListWithFilter(pageOption)
				if err != nil {
					return groups, resp, err
				} else {
					groups = append(groups, groupPage...)
					pagesRetreived++
				}
			} else {
				break
			}
		}
	}
	return groups, resp, err
}

// GetByID gets a group from OKTA by the Gropu ID. An error is returned if the group is not found
func (g *GroupsService) GetByID(groupID string) (*Group, *Response, error) {

	u := fmt.Sprintf("groups/%v", groupID)
	req, err := g.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, err
	}

	group := new(Group)

	resp, err := g.client.Do(req, group)

	if err != nil {
		return nil, resp, err
	}

	return group, resp, err
}

// GetUsers returns the members in a group
//   Pass in an optional GroupFilterOptions struct to filter the results
//   The Users in the group are returned
func (g *GroupsService) GetUsers(groupID string, opt *GroupUserFilterOptions) (users []User, resp *Response, err error) {
	pagesRetreived := 0
	var u string
	if opt.NextURL != nil {
		u = opt.NextURL.String()
	} else {
		u = fmt.Sprintf("groups/%v/users", groupID)

		if opt.Limit == 0 {
			opt.Limit = defaultLimit
		}

		u, _ = addOptions(u, opt)
	}

	req, err := g.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, err
	}
	resp, err = g.client.Do(req, &users)

	if err != nil {
		return nil, resp, err
	}

	pagesRetreived++
	if (opt.NumberOfPages > 0 && pagesRetreived < opt.NumberOfPages) || opt.GetAllPages {

		for {

			if pagesRetreived == opt.NumberOfPages {
				break
			}
			if resp.NextURL != nil {

				var userPage []User
				pageOpts := new(GroupUserFilterOptions)
				pageOpts.NextURL = resp.NextURL
				pageOpts.Limit = opt.Limit
				pageOpts.NumberOfPages = 1

				userPage, resp, err = g.GetUsers(groupID, pageOpts)
				if err != nil {
					return users, resp, err
				} else {
					users = append(users, userPage...)
					pagesRetreived++
				}
			} else {
				break
			}

		}
	}

	return users, resp, err
}

// Add - Adds an OKTA Mastered Group with name and description. GroupName is required.
func (g *GroupsService) Add(groupName string, groupDescription string) (*Group, *Response, error) {

	if groupName == "" {
		return nil, nil, errors.New("groupName parameter is required for ADD")
	}

	newGroup := newGroup{}
	newGroup.Profile.Name = groupName
	newGroup.Profile.Description = groupDescription

	u := fmt.Sprintf("groups")

	req, err := g.client.NewRequest("POST", u, newGroup)

	if err != nil {
		return nil, nil, err
	}

	group := new(Group)

	resp, err := g.client.Do(req, group)

	if err != nil {
		return nil, resp, err
	}

	return group, resp, err
}

// Delete - Delets an OKTA Mastered Group with ID
func (g *GroupsService) Delete(groupID string) (*Response, error) {

	if groupID == "" {
		return nil, errors.New("groupID parameter is required for Delete")
	}
	u := fmt.Sprintf("groups/%v", groupID)

	req, err := g.client.NewRequest("DELETE", u, nil)

	if err != nil {
		return nil, err
	}

	resp, err := g.client.Do(req, nil)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// GroupUserFilterOptions is a struct that you populate which will limit or control group fetches and searches
//  The values here will coorelate to the search filtering allowed in the OKTA API. These values are turned into Query Parameters
type GroupUserFilterOptions struct {
	Limit         int      `url:"limit,omitempty"`
	NextURL       *url.URL `url:"-"`
	GetAllPages   bool     `url:"-"`
	NumberOfPages int      `url:"-"`
}

type newGroup struct {
	Profile struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"profile"`
}
