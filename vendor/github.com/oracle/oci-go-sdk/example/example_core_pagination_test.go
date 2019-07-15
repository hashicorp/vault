// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Example code for Core Services API
//
package example

import (
	"context"
	"fmt"
	"log"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/example/helpers"
)

// ExampleListShapes_Pagination demostrate how to use page parameter
func ExampleListShapes_Pagination() {
	c, err := core.NewComputeClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	request := core.ListShapesRequest{
		CompartmentId: helpers.CompartmentID(),
	}

	// to show how pagination works, reduce number of items to return in a paginated "List" call
	request.Limit = common.Int(2)

	listShapesFunc := func(request core.ListShapesRequest) (core.ListShapesResponse, error) {
		return c.ListShapes(context.Background(), request)
	}

	for r, err := listShapesFunc(request); ; r, err = listShapesFunc(request) {
		helpers.FatalIfError(err)

		log.Printf("list shapes returns: %v", r.Items)

		if r.OpcNextPage != nil {
			// if there are more items in next page, fetch items from next page
			request.Page = r.OpcNextPage
		} else {
			// no more result, break the loop
			break
		}
	}

	fmt.Println("list shapes completed")

	// Output:
	// list shapes completed
}
