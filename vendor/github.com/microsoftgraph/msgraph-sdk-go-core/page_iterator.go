package msgraphgocore

import (
	"context"
	"errors"
	"net/url"
	"reflect"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

const PageIteratorErrorRegistryKey = "PAGE_ITERATOR_ERROR_REGISTRY_KEY"

// PageIterator represents an iterator object that can be used to get subsequent pages of a collection.
type PageIterator[T interface{}] struct {
	currentPage     PageResult[T]
	reqAdapter      abstractions.RequestAdapter
	pauseIndex      int
	constructorFunc serialization.ParsableFactory
	headers         *abstractions.RequestHeaders
	reqOptions      []abstractions.RequestOption
	errorMappings   abstractions.ErrorMappings
}

// PageResult represents a page object built from a graph response object
type PageResult[T interface{}] struct {
	oDataNextLink  *string
	oDataDeltaLink *string
	value          []T
}

func (p *PageResult[T]) getValue() []T {
	if p == nil {
		return nil
	}

	return p.value
}

func (p *PageResult[T]) getOdataNextLink() *string {
	if p == nil {
		return nil
	}

	return p.oDataNextLink
}

// NewPageIterator creates an iterator instance
//
// It has three parameters. res is the graph response from the initial request and represents the first page.
// reqAdapter is used for getting the next page and constructorFunc is used for serializing next page's response to the specified type.
func NewPageIterator[T interface{}](res interface{}, reqAdapter abstractions.RequestAdapter, constructorFunc serialization.ParsableFactory) (*PageIterator[T], error) {
	if reqAdapter == nil {
		return nil, errors.New("reqAdapter can't be nil")
	}

	page, err := convertToPage[T](res)
	if err != nil {
		return nil, err
	}

	errorMapping := getErrorMapper(PageIteratorErrorRegistryKey)

	return &PageIterator[T]{
		currentPage:     page,
		reqAdapter:      reqAdapter,
		pauseIndex:      0,
		constructorFunc: constructorFunc,
		headers:         abstractions.NewRequestHeaders(),
		errorMappings:   errorMapping,
	}, nil
}

// Iterate traverses all pages and enumerates all items in the current page and returns an error if something goes wrong.
//
// Iterate receives a callback function which is called with each item in the current page as an argument. The callback function
// returns a boolean. To traverse and enumerate all pages always return true and to pause traversal and enumeration
// return false from the callback.
//
// Example
//
//	pageIterator, err := NewPageIterator(resp, reqAdapter, parsableFactory)
//	callbackFunc := func (pageItem interface{}) bool {
//	    fmt.Println(page item.GetDisplayName())
//	    return true
//	}
//	err := pageIterator.Iterate(context.Background(), callbackFunc)
func (pI *PageIterator[T]) Iterate(context context.Context, callback func(pageItem T) bool) error {
	for {
		keepIterating := pI.enumerate(callback)

		if !keepIterating {
			// Callback returned false, stop iterating through pages.
			return nil
		}

		if pI.currentPage.getOdataNextLink() == nil || *pI.currentPage.getOdataNextLink() == "" {
			return nil
		}

		nextPage, err := pI.next(context)
		if err != nil {
			return err
		}

		pI.currentPage = nextPage
		pI.pauseIndex = 0 // when moving to the next page reset pauseIndex
	}
}

// SetHeaders provides headers for requests made to get subsequent pages
//
// Headers in the initial request -- request to get the first page -- are not included in subsequent page requests.
func (pI *PageIterator[T]) SetHeaders(headers *abstractions.RequestHeaders) {
	pI.headers = headers
}

// SetReqOptions provides configuration for handlers during requests for subsequent pages
func (pI *PageIterator[T]) SetReqOptions(reqOptions []abstractions.RequestOption) {
	pI.reqOptions = reqOptions
}

// GetOdataNextLink returns the @odata.nextLink value in the current page result.
func (pI *PageIterator[T]) GetOdataNextLink() *string {
	return pI.currentPage.oDataNextLink
}

// GetOdataDeltaLink returns the @odata.deltaLink value in current paged result.
func (pI *PageIterator[T]) GetOdataDeltaLink() *string {
	return pI.currentPage.oDataDeltaLink
}

func (pI *PageIterator[T]) next(context context.Context) (PageResult[T], error) {
	var page PageResult[T]

	resp, err := pI.fetchNextPage(context)
	if err != nil {
		return page, err
	}

	page, err = convertToPage[T](resp)
	if err != nil {
		return page, err
	}

	return page, nil
}

func (pI *PageIterator[T]) fetchNextPage(context context.Context) (serialization.Parsable, error) {
	var graphResponse serialization.Parsable
	var err error

	if pI.currentPage.getOdataNextLink() == nil {
		return graphResponse, nil
	}

	nextLink, err := url.Parse(*pI.currentPage.getOdataNextLink())
	if err != nil {
		return nil, errors.New("parsing nextLink url failed")
	}

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.GET
	requestInfo.SetUri(*nextLink)
	requestInfo.Headers.AddAll(pI.headers)
	requestInfo.AddRequestOptions(pI.reqOptions)

	graphResponse, err = pI.reqAdapter.Send(context, requestInfo, pI.constructorFunc, pI.errorMappings)
	if err != nil {
		return nil, err
	}

	return graphResponse, nil
}

func (pI *PageIterator[T]) enumerate(callback func(item T) bool) bool {
	keepIterating := true

	pageItems := pI.currentPage.getValue()
	if pageItems == nil {
		return false
	}

	// the current page has no items to enumerate
	if pI.currentPage.getValue() == nil {
		return false
	}

	// start/continue enumerating page items from  pauseIndex.
	// this makes it possible to resume iteration from where we paused iteration.
	for i := pI.pauseIndex; i < len(pageItems); i++ {
		keepIterating = callback(pageItems[i])

		// Set pauseIndex so that we know where to resume from.
		// Resumes from the next item
		pI.pauseIndex = i + 1

		if !keepIterating {
			break
		}
	}

	return keepIterating
}

// PageWithOdataNextLink represents a contract with the GetOdataNextLink() method
type PageWithOdataNextLink interface {
	GetOdataNextLink() *string
}

// PageWithOdataDeltaLink represents a contract with the GetOdataDeltaLink() method
type PageWithOdataDeltaLink interface {
	GetOdataDeltaLink() *string
}

func convertToPage[T interface{}](response interface{}) (PageResult[T], error) {
	var page PageResult[T]

	if response == nil {
		return page, errors.New("response cannot be nil")
	}

	method := reflect.ValueOf(response).MethodByName("GetValue")
	if method.Kind() == reflect.Invalid {
		return page, errors.New("value property missing in response object")
	}
	value := method.Call(nil)[0]

	// Collect all entities in the value slice.
	// This converts a graph slice ie []graph.User to a dynamic slice []interface{}
	collected := make([]T, 0)
	for i := 0; i < value.Len(); i++ {
		collected = append(collected, value.Index(i).Interface().(T))
	}

	parsablePage, ok := response.(PageWithOdataNextLink)
	if !ok {
		return page, errors.New("response does not have next link accessor")
	}

	deltablePage, ok := response.(PageWithOdataDeltaLink)
	if ok {
		page.oDataDeltaLink = deltablePage.GetOdataDeltaLink()
	}

	page.oDataNextLink = parsablePage.GetOdataNextLink()
	page.value = collected

	return page, nil
}
