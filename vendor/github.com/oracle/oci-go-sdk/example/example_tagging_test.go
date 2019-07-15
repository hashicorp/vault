// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Example code for Tagging Service API
//

package example

import (
	"context"
	"fmt"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/identity"
)

// ExampleTagging shows the sample for tag and tagNamespace operations: create, update, get, list etc...
func ExampleTagging() {
	c, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	ctx := context.Background()
	tagNamespaceID := createTagNamespace(ctx, c, common.String("GOSDKSampleTagNamespaceName"))
	fmt.Println("tag namespace created")

	tagName := common.String("GOSDKSampleTagName")
	createTag(ctx, c, tagNamespaceID, tagName)
	fmt.Println("tag created")

	// get tag
	getTagReq := identity.GetTagRequest{
		TagNamespaceId: tagNamespaceID,
		TagName:        tagName,
	}
	_, err = c.GetTag(ctx, getTagReq)
	helpers.FatalIfError(err)
	fmt.Println("get tag")

	// list tags, list operations are paginated and take a "page" parameter
	// to allow you to get the next batch of items from the server
	// for pagination sample, please refer to 'example_core_pagination_test.go'
	listTagReq := identity.ListTagsRequest{
		TagNamespaceId: tagNamespaceID,
	}
	_, err = c.ListTags(ctx, listTagReq)
	helpers.FatalIfError(err)
	fmt.Println("list tag")

	// get tag namespace
	getTagNamespaceReq := identity.GetTagNamespaceRequest{
		TagNamespaceId: tagNamespaceID,
	}
	_, err = c.GetTagNamespace(ctx, getTagNamespaceReq)
	helpers.FatalIfError(err)
	fmt.Println("get tag namespace")

	// list tag namespaces
	listTagNamespaceReq := identity.ListTagNamespacesRequest{
		CompartmentId: helpers.CompartmentID(),
	}
	_, err = c.ListTagNamespaces(ctx, listTagNamespaceReq)
	helpers.FatalIfError(err)
	fmt.Println("list tag namespace")

	// retire a tag namespace by using the update tag namespace operation
	updateTagNamespaceReq := identity.UpdateTagNamespaceRequest{
		TagNamespaceId: tagNamespaceID,
		UpdateTagNamespaceDetails: identity.UpdateTagNamespaceDetails{
			IsRetired: common.Bool(true),
		},
	}

	_, err = c.UpdateTagNamespace(ctx, updateTagNamespaceReq)
	helpers.FatalIfError(err)
	fmt.Println("tag namespace retired")

	// retire a tag by using the update tag operation
	updateTagReq := identity.UpdateTagRequest{
		TagNamespaceId: tagNamespaceID,
		TagName:        tagName,
		UpdateTagDetails: identity.UpdateTagDetails{
			IsRetired: common.Bool(true),
		},
	}
	_, err = c.UpdateTag(ctx, updateTagReq)
	helpers.FatalIfError(err)
	fmt.Println("tag retired")

	// reactivate a tag namespace
	updateTagNamespaceReq = identity.UpdateTagNamespaceRequest{
		TagNamespaceId: tagNamespaceID,
		UpdateTagNamespaceDetails: identity.UpdateTagNamespaceDetails{
			// reactivate a tag namespace by using the update tag namespace operation
			IsRetired: common.Bool(false),
		},
	}

	_, err = c.UpdateTagNamespace(ctx, updateTagNamespaceReq)
	helpers.FatalIfError(err)
	fmt.Println("tag namespace reactivated")

	// Output:
	// tag namespace created
	// tag created
	// get tag
	// list tag
	// get tag namespace
	// list tag namespace
	// tag namespace retired
	// tag retired
	// tag namespace reactivated
}

// ExampleFreeformAndDefinedTag shows how to use freeform and defined tags
func ExampleFreeformAndDefinedTag() {
	// create a tag namespace and two tags
	identityClient, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	ctx := context.Background()

	tagNamespaceName := "GOSDKSampleTagNamespaceName_1"
	tagNamespaceID := createTagNamespace(ctx, identityClient, common.String(tagNamespaceName))
	fmt.Println("tag namespace created")

	tagName := "GOSDKSampleTagName_1"
	createTag(ctx, identityClient, tagNamespaceID, common.String(tagName))
	fmt.Println("tag1 created")

	tagName2 := "GOSDKSampleTagName_2"
	createTag(ctx, identityClient, tagNamespaceID, common.String(tagName2))
	fmt.Println("tag2 created")

	// We can assign freeform and defined tags at resource creation time. Freeform tags are a dictionary of
	// string-to-string, where the key is the tag name and the value is the tag value.
	//
	// Defined tags are a dictionary where the key is the tag namespace (string) and the value is another dictionary. In
	// this second dictionary, the key is the tag name (string) and the value is the tag value. The tag names have to
	// correspond to the name of a tag within the specified namespace (and the namespace must exist).
	freeformTags := map[string]string{"free": "form", "another": "item"}
	definedTags := map[string]map[string]interface{}{
		tagNamespaceName: {
			tagName:  "hello",
			tagName2: "world",
		},
	}

	coreClient, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	// create a new VCN with tags
	createVCNReq := core.CreateVcnRequest{
		CreateVcnDetails: core.CreateVcnDetails{
			CidrBlock:     common.String("10.0.0.0/16"),
			CompartmentId: helpers.CompartmentID(),
			DisplayName:   common.String("GOSDKSampleVCNName"),
			DnsLabel:      common.String("vcndns"),
			FreeformTags:  freeformTags,
			DefinedTags:   definedTags,
		},
	}

	resp, err := coreClient.CreateVcn(ctx, createVCNReq)

	if err != nil && resp.RawResponse.StatusCode == 404 {
		// You may get a 404 if you create/reactivate a tag and try and use it straight away. If you have a delay/sleep between
		// creating the tag and then using it (or alternatively retry the 404) that may resolve the issue.
		time.Sleep(time.Second * 10)
		resp, err = coreClient.CreateVcn(ctx, createVCNReq)
	}

	helpers.FatalIfError(err)
	fmt.Println("VCN created with tags")

	// replace the tag
	freeformTags = map[string]string{"total": "replaced"}

	// update the tag value
	definedTags[tagNamespaceName][tagName2] = "replaced"

	// update the VCN with different tag values
	updateVCNReq := core.UpdateVcnRequest{
		VcnId: resp.Id,
		UpdateVcnDetails: core.UpdateVcnDetails{
			FreeformTags: freeformTags,
			DefinedTags:  definedTags,
		},
	}
	_, err = coreClient.UpdateVcn(ctx, updateVCNReq)
	helpers.FatalIfError(err)
	fmt.Println("VCN tag updated")

	// remove the tag from VCN
	updateVCNReq.FreeformTags = nil
	updateVCNReq.DefinedTags = nil
	_, err = coreClient.UpdateVcn(ctx, updateVCNReq)
	helpers.FatalIfError(err)
	fmt.Println("VCN tag removed")

	defer func() {
		request := core.DeleteVcnRequest{
			VcnId: resp.Id,
		}

		_, err = coreClient.DeleteVcn(ctx, request)
		helpers.FatalIfError(err)
		fmt.Println("VCN deleted")
	}()

	// Output:
	// tag namespace created
	// tag1 created
	// tag2 created
	// VCN created with tags
	// VCN tag updated
	// VCN tag removed
	// VCN deleted
}

func createTagNamespace(ctx context.Context, client identity.IdentityClient, name *string) *string {
	req := identity.CreateTagNamespaceRequest{}
	req.CompartmentId = helpers.CompartmentID()
	req.Name = name
	req.Description = common.String("GOSDK Sample TagNamespace Description")

	resp, err := client.CreateTagNamespace(context.Background(), req)
	helpers.FatalIfError(err)

	return resp.Id
}

func createTag(ctx context.Context, client identity.IdentityClient, tagNamespaceID *string, tagName *string) *string {
	req := identity.CreateTagRequest{
		TagNamespaceId: tagNamespaceID,
	}

	req.Name = tagName
	req.Description = common.String("GOSDK Sample Tag Description")
	req.FreeformTags = map[string]string{"GOSDKSampleTagKey": "GOSDKSampleTagValue"}

	resp, err := client.CreateTag(context.Background(), req)
	helpers.FatalIfError(err)

	return resp.Id
}
