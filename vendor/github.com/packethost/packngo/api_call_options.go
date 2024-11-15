package packngo

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ListSortDirection string

const (
	SortDirectionAsc  ListSortDirection = "asc"
	SortDirectionDesc ListSortDirection = "desc"
)

// GetOptions are options common to Equinix Metal API GET requests
type GetOptions struct {
	// Includes are a list of fields to expand in the request results.
	//
	// For resources that contain collections of other resources, the Equinix Metal API
	// will only return the `Href` value of these resources by default. In
	// nested API Go types, this will result in objects that have zero values in
	// all fiends except their `Href` field. When an object's associated field
	// name is "included", the returned fields will be Uumarshalled into the
	// nested object. Field specifiers can use a dotted notation up to three
	// references deep. (For example, "memberships.projects" can be used in
	// ListUsers.)
	Includes []string `url:"include,omitempty,comma"`

	// Excludes reduce the size of the API response by removing nested objects
	// that may be returned.
	//
	// The default behavior of the Equinix Metal API is to "exclude" fields, but some
	// API endpoints have an "include" behavior on certain fields. Nested Go
	// types unmarshalled into an "excluded" field will only have a values in
	// their `Href` field.
	Excludes []string `url:"exclude,omitempty,comma"`

	// QueryParams for API URL, used for arbitrary filters
	QueryParams map[string]string `url:"-"`

	// Page is the page of results to retrieve for paginated result sets
	Page int `url:"page,omitempty"`

	// PerPage is the number of results to return per page for paginated result
	// sets,
	PerPage int `url:"per_page,omitempty"`

	// Search is a special API query parameter that, for resources that support
	// it, will filter results to those with any one of various fields matching
	// the supplied keyword.  For example, a resource may have a defined search
	// behavior matches either a name or a fingerprint field, while another
	// resource may match entirely different fields.  Search is currently
	// implemented for SSHKeys and uses an exact match.
	Search string `url:"search,omitempty"`

	SortBy        string            `url:"sort_by,omitempty"`
	SortDirection ListSortDirection `url:"sort_direction,omitempty"`

	Meta meta `url:"-"`
}

type ListOptions = GetOptions
type SearchOptions = GetOptions

type QueryAppender interface {
	WithQuery(path string) string // we use this in all List functions (urlQuery)
	GetPage() int                 // we use this in List
	Including(...string)          // we use this in Device List to add facility
	Excluding(...string)
}

// GetOptions returns GetOptions from GetOptions (and is nil-receiver safe)
func (g *GetOptions) GetOptions() *GetOptions {
	getOpts := GetOptions{}
	if g != nil {
		getOpts.Includes = g.Includes
		getOpts.Excludes = g.Excludes
	}
	return &getOpts
}

func (g *GetOptions) WithQuery(apiPath string) string {
	params := g.Encode()
	if params != "" {
		// parse path, take existing vars
		return fmt.Sprintf("%s?%s", apiPath, params)
	}
	return apiPath
}

// OptionsGetter provides GetOptions
type OptionsGetter interface {
	GetOptions() *GetOptions
}

func (g *GetOptions) GetPage() int { // guaranteed int
	if g == nil {
		return 0
	}
	return g.Page
}

func (g *GetOptions) CopyOrNew() *GetOptions {
	if g == nil {
		return &GetOptions{}
	}
	ret := *g
	return &ret
}

func (g *GetOptions) Filter(key, value string) *GetOptions {
	return g.AddParam(key, value)
}

// AddParam adds key=value to URL path
func (g *GetOptions) AddParam(key, value string) *GetOptions {
	ret := g.CopyOrNew()
	if ret.QueryParams == nil {
		ret.QueryParams = map[string]string{}
	}
	ret.QueryParams[key] = value
	return ret
}

// Including ensures that the variadic refs are included in a copy of the
// options, resulting in expansion of the the referred sub-resources. Unknown
// values within refs will be silently ignore by the API.
func (g *GetOptions) Including(refs ...string) *GetOptions {
	ret := g.CopyOrNew()
	for _, v := range refs {
		if !contains(ret.Includes, v) {
			ret.Includes = append(ret.Includes, v)
		}
	}
	return ret
}

// Excluding ensures that the variadic refs are included in the "Excluded" param
// in a copy of the options.
// Unknown values within refs will be silently ignore by the API.
func (g *GetOptions) Excluding(refs ...string) *GetOptions {
	ret := g.CopyOrNew()
	for _, v := range refs {
		if !contains(ret.Excludes, v) {
			ret.Excludes = append(ret.Excludes, v)
		}
	}
	return ret
}

func stripQuery(inURL string) string {
	u, _ := url.Parse(inURL)
	u.RawQuery = ""
	return u.String()
}

// nextPage is common and extracted from all List functions
func nextPage(meta meta, opts *GetOptions) (path string) {
	if meta.Next != nil && (opts.GetPage() == 0) {
		optsCopy := opts.CopyOrNew()
		optsCopy.Page = meta.CurrentPageNum + 1
		return optsCopy.WithQuery(stripQuery(meta.Next.Href))
	}
	if opts != nil {
		opts.Meta = meta
	}
	return ""
}

const (
	IncludeParam       = "include"
	ExcludeParam       = "exclude"
	PageParam          = "page"
	PerPageParam       = "per_page"
	SearchParam        = "search"
	SortByParam        = "sort_by"
	SortDirectionParam = "sort_direction"
)

// Encode generates a URL query string ("?foo=bar")
func (g *GetOptions) Encode() string {
	if g == nil {
		return ""
	}
	v := url.Values{}
	for k, val := range g.QueryParams {
		v.Add(k, val)
	}
	// the names parameters will rewrite arbitrary options
	if g.Includes != nil && len(g.Includes) > 0 {
		v.Add(IncludeParam, strings.Join(g.Includes, ","))
	}
	if g.Excludes != nil && len(g.Excludes) > 0 {
		v.Add(ExcludeParam, strings.Join(g.Excludes, ","))
	}
	if g.Page != 0 {
		v.Add(PageParam, strconv.Itoa(g.Page))
	}
	if g.PerPage != 0 {
		v.Add(PerPageParam, strconv.Itoa(g.PerPage))
	}
	if g.Search != "" {
		v.Add(SearchParam, g.Search)
	}
	if g.SortBy != "" {
		v.Add(SortByParam, g.SortBy)
	}
	if g.SortDirection != "" {
		v.Add(SortDirectionParam, string(g.SortDirection))
	}
	return v.Encode()
}

/*
// Encode generates a URL query string ("?foo=bar")
func (g *GetOptions) Encode() string {
	urlValues, _ := query.Values(g)
	return urlValues.Encode()
}
*/
