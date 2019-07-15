package example

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcemanager"
)

// ExampleResourceManager for how to do CRUD for Resource Manager Stack
// The comparement id is read from the environment variable OCI_COMPARTMENT_ID
func ExampleResourceManager() {
	provider := common.DefaultConfigProvider()
	client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
	helpers.FatalIfError(err)

	ctx := context.Background()

	stackID := createStack(ctx, provider, client)
	defer deleteStack(ctx, stackID, client)
	listStacks(ctx, client)
	updateStack(ctx, stackID, client)
	getStack(ctx, stackID, client)

	// Output:
	// create stack completed
	// list stacks completed
	// update stack completed
	// get stack completed
	// delete stack completed
}

func createStack(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient) string {
	stackName := fmt.Sprintf("test-%s", helpers.GetRandomString(8))
	region, _ := provider.Region()
	tenancyOcid, _ := provider.TenancyOCID()

	// create resource manager stack with type ZIP_UPLOAD by passing a base64 encoded Terraform zip string
	// user has multiple ways to create stack, details check https://docs.cloud.oracle.com/iaas/api/#/en/resourcemanager/20180917/datatypes/CreateConfigSourceDetails
	req := resourcemanager.CreateStackRequest{
		CreateStackDetails: resourcemanager.CreateStackDetails{
			CompartmentId: helpers.CompartmentID(),
			ConfigSource: resourcemanager.CreateZipUploadConfigSourceDetails{
				WorkingDirectory:     common.String("vcn"),
				ZipFileBase64Encoded: common.String("[pls use your base64 encoded TF template]"),
			},
			DisplayName: common.String(stackName),
			Description: common.String(fmt.Sprintf("%s-description", stackName)),
			Variables: map[string]string{
				"compartment_ocid": *helpers.CompartmentID(),
				"region":           region,
				"tenancy_ocid":     tenancyOcid,
			},
		},
	}

	stackResp, err := client.CreateStack(ctx, req)
	helpers.FatalIfError(err)

	fmt.Println("create stack completed")
	return *stackResp.Stack.Id
}

func updateStack(ctx context.Context, stackID string, client resourcemanager.ResourceManagerClient) {
	stackName := fmt.Sprintf("test-v1-%s", helpers.GetRandomString(8))

	// update displayName and description of resource manager stack
	req := resourcemanager.UpdateStackRequest{
		StackId: common.String(stackID),
		UpdateStackDetails: resourcemanager.UpdateStackDetails{
			DisplayName: common.String(stackName),
			Description: common.String(fmt.Sprintf("%s-description", stackName)),
		},
	}

	_, err := client.UpdateStack(ctx, req)
	helpers.FatalIfError(err)

	fmt.Println("update stack completed")
}

func listStacks(ctx context.Context, client resourcemanager.ResourceManagerClient) {
	req := resourcemanager.ListStacksRequest{
		CompartmentId: helpers.CompartmentID(),
	}

	// list resource manager stack
	_, err := client.ListStacks(ctx, req)
	helpers.FatalIfError(err)

	fmt.Println("list stacks completed")
}

func getStack(ctx context.Context, stackID string, client resourcemanager.ResourceManagerClient) {
	req := resourcemanager.GetStackRequest{
		StackId: common.String(stackID),
	}

	// get details a particular resource manager stack
	_, err := client.GetStack(ctx, req)
	helpers.FatalIfError(err)

	fmt.Println("get stack completed")
}

func deleteStack(ctx context.Context, stackID string, client resourcemanager.ResourceManagerClient) {
	req := resourcemanager.DeleteStackRequest{
		StackId: common.String(stackID),
	}

	// delete a resource manager stack
	_, err := client.DeleteStack(ctx, req)
	helpers.FatalIfError(err)

	fmt.Println("delete stack completed")
}
