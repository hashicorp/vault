package linodego

/**
 * Pagination and Filtering types and helpers
 */

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// PageOptions are the pagination parameters for List endpoints
type PageOptions struct {
	Page    int `json:"page"    url:"page,omitempty"`
	Pages   int `json:"pages"   url:"pages,omitempty"`
	Results int `json:"results" url:"results,omitempty"`
}

// ListOptions are the pagination and filtering (TODO) parameters for endpoints
// nolint
type ListOptions struct {
	*PageOptions
	PageSize int    `json:"page_size"`
	Filter   string `json:"filter"`

	// QueryParams allows for specifying custom query parameters on list endpoint
	// calls. QueryParams should be an instance of a struct containing fields with
	// the `query` tag.
	QueryParams any
}

// NewListOptions simplified construction of ListOptions using only
// the two writable properties, Page and Filter
func NewListOptions(page int, filter string) *ListOptions {
	return &ListOptions{PageOptions: &PageOptions{Page: page}, Filter: filter}
}

// Hash returns the sha256 hash of the provided ListOptions.
// This is necessary for caching purposes.
func (l ListOptions) Hash() (string, error) {
	data, err := json.Marshal(l)
	if err != nil {
		return "", fmt.Errorf("failed to cache ListOptions: %w", err)
	}

	h := sha256.New()

	h.Write(data)

	return hex.EncodeToString(h.Sum(nil)), nil
}

func applyListOptionsToRequest(opts *ListOptions, req *resty.Request) error {
	if opts == nil {
		return nil
	}

	if opts.QueryParams != nil {
		params, err := flattenQueryStruct(opts.QueryParams)
		if err != nil {
			return fmt.Errorf("failed to apply list options: %w", err)
		}

		req.SetQueryParams(params)
	}

	if opts.PageOptions != nil && opts.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(opts.Page))
	}

	if opts.PageSize > 0 {
		req.SetQueryParam("page_size", strconv.Itoa(opts.PageSize))
	}

	if len(opts.Filter) > 0 {
		req.SetHeader("X-Filter", opts.Filter)
	}

	return nil
}

type PagedResponse interface {
	endpoint(...any) string
	castResult(*resty.Request, string) (int, int, error)
}

// flattenQueryStruct flattens a structure into a Resty-compatible query param map.
// Fields are mapped using the `query` struct tag.
func flattenQueryStruct(val any) (map[string]string, error) {
	result := make(map[string]string)

	reflectVal := reflect.ValueOf(val)

	// Deref pointer if necessary
	if reflectVal.Kind() == reflect.Pointer {
		if reflectVal.IsNil() {
			return nil, fmt.Errorf("QueryParams is a nil pointer")
		}
		reflectVal = reflect.Indirect(reflectVal)
	}

	if reflectVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf(
			"expected struct type for the QueryParams but got: %s",
			reflectVal.Kind().String(),
		)
	}

	valType := reflectVal.Type()

	for i := range valType.NumField() {
		currentField := valType.Field(i)

		queryTag, ok := currentField.Tag.Lookup("query")
		// Skip untagged fields
		if !ok {
			continue
		}

		valField := reflectVal.FieldByName(currentField.Name)
		if !valField.IsValid() {
			return nil, fmt.Errorf("invalid query param tag: %s", currentField.Name)
		}

		// Skip if it's a zero value
		if valField.IsZero() {
			continue
		}

		// Deref the pointer is necessary
		if valField.Kind() == reflect.Pointer {
			valField = reflect.Indirect(valField)
		}

		fieldString, err := queryFieldToString(valField)
		if err != nil {
			return nil, err
		}

		result[queryTag] = fieldString
	}

	return result, nil
}

func queryFieldToString(value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.String:
		return value.String(), nil
	case reflect.Int64, reflect.Int32, reflect.Int:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil
	default:
		return "", fmt.Errorf("unsupported query param type: %s", value.Type().Name())
	}
}

type legacyPagedResponse[T any] struct {
	*PageOptions
	Data []T `json:"data"`
}
