// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Example code for Audit API
//
package example

import (
	"context"
	"fmt"
	"time"

	"github.com/oracle/oci-go-sdk/audit"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
)

func ExampleListEvents() {
	c, clerr := audit.NewAuditClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	// list events for last 5 hour
	req := audit.ListEventsRequest{
		CompartmentId: helpers.CompartmentID(),
		StartTime:     &common.SDKTime{time.Now().Add(time.Hour * -5)},
		EndTime:       &common.SDKTime{time.Now()},
	}

	_, err := c.ListEvents(context.Background(), req)
	helpers.FatalIfError(err)

	//log.Printf("events returned back: %v", resp.Items)
	fmt.Println("list events completed")

	// Output:
	// list events completed
}
