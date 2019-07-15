// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Example code for Database API
//
package example

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/database"
	"github.com/oracle/oci-go-sdk/example/helpers"
)

func ListPreviewDBVersion() {
	c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	listAdbPreviewVersionsReq := database.ListAutonomousDbPreviewVersionsRequest{
		CompartmentId: helpers.CompartmentID(),
	}

	_, err := c.ListAutonomousDbPreviewVersions(context.Background(), listAdbPreviewVersionsReq)
	helpers.FatalIfError(err)

	fmt.Println("List Preview Versions successful")
	// Output:
	// List Preview Versions successful
}

func ExampleCreateAdb() {
	c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	createDbDetails := database.CreateAutonomousDatabaseDetails{
		CompartmentId:        helpers.CompartmentID(),
		DbName:               common.String("gosdkdb"),
		CpuCoreCount:         common.Int(1),
		DataStorageSizeInTBs: common.Int(1),
		AdminPassword:        common.String("DBaaS12345_#"),
		IsAutoScalingEnabled: common.Bool(true),
	}

	createadbReq := database.CreateAutonomousDatabaseRequest{
		CreateAutonomousDatabaseDetails: createDbDetails,
	}

	_, err := c.CreateAutonomousDatabase(context.Background(), createadbReq)
	helpers.FatalIfError(err)

	fmt.Println("create adb successful")

	// Output:
	// create adb successful
}

func ExampleCreateAdbPreview() {
	c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	createDbDetails := database.CreateAutonomousDatabaseDetails{
		CompartmentId:                            helpers.CompartmentID(),
		DbName:                                   common.String("gosdkPreview"),
		CpuCoreCount:                             common.Int(1),
		DataStorageSizeInTBs:                     common.Int(1),
		AdminPassword:                            common.String("DBaaS12345_#"),
		IsAutoScalingEnabled:                     common.Bool(false),
		IsPreviewVersionWithServiceTermsAccepted: common.Bool(true),
	}

	createadbReq := database.CreateAutonomousDatabaseRequest{
		CreateAutonomousDatabaseDetails: createDbDetails,
	}

	_, err := c.CreateAutonomousDatabase(context.Background(), createadbReq)
	helpers.FatalIfError(err)

	fmt.Println("create adb Preview successful")

	// Output:
	// create adb successful
}

func ExampleUpdateAdb() {
	c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	updateDbDetails := database.UpdateAutonomousDatabaseDetails{
		CpuCoreCount:         common.Int(2),
		DataStorageSizeInTBs: common.Int(2),
		IsAutoScalingEnabled: common.Bool(false),
	}

	updateReq := database.UpdateAutonomousDatabaseRequest{
		AutonomousDatabaseId:            common.String("replacewithvalidocid"),
		UpdateAutonomousDatabaseDetails: updateDbDetails,
	}
	_, err := c.UpdateAutonomousDatabase(context.Background(), updateReq)
	helpers.FatalIfError(err)

	fmt.Println("update adb successful")

	// Output:
	// update adb successful
}

func ExampleUpdateAdbAcl() {
	c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	updateDbDetails := database.UpdateAutonomousDatabaseDetails{
		WhitelistedIps: []string{"1.1.1.1/28", "3.3.3.3"},
	}

	updateReq := database.UpdateAutonomousDatabaseRequest{
		AutonomousDatabaseId:            common.String("replacewithvalidocid"),
		UpdateAutonomousDatabaseDetails: updateDbDetails,
	}
	_, err := c.UpdateAutonomousDatabase(context.Background(), updateReq)
	helpers.FatalIfError(err)

	fmt.Println("update adb acl successful")

	// Output:
	// update adb acl successful
}

func ExampleUpdateAdbLisenceType() {
	c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	updateDbDetails := database.UpdateAutonomousDatabaseDetails{
		LicenseModel: database.UpdateAutonomousDatabaseDetailsLicenseModelLicenseIncluded,
	}

	updateReq := database.UpdateAutonomousDatabaseRequest{
		AutonomousDatabaseId:            common.String("replacewithvalidocid"),
		UpdateAutonomousDatabaseDetails: updateDbDetails,
	}
	_, err := c.UpdateAutonomousDatabase(context.Background(), updateReq)
	helpers.FatalIfError(err)

	fmt.Println("update adb license type successful")

	// Output:
	// update adb license type successful
}
