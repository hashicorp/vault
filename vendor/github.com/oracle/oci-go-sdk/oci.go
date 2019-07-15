/*
This is the official Go SDK for Oracle Cloud Infrastructure

Installation

Refer to https://github.com/oracle/oci-go-sdk/blob/master/README.md#installing for installation instructions.

Configuration

Refer to https://github.com/oracle/oci-go-sdk/blob/master/README.md#configuring for configuration instructions.

Quickstart

The following example shows how to get started with the SDK. The example belows creates an identityClient
struct with the default configuration. It then utilizes the identityClient to list availability domains and prints
them out to stdout

	import (
		"context"
		"fmt"

		"github.com/oracle/oci-go-sdk/common"
		"github.com/oracle/oci-go-sdk/identity"
	)

	func main() {
		c, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// The OCID of the tenancy containing the compartment.
		tenancyID, err := common.DefaultConfigProvider().TenancyOCID()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		request := identity.ListAvailabilityDomainsRequest{
			CompartmentId: &tenancyID,
		}

		r, err := c.ListAvailabilityDomains(context.Background(), request)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("List of available domains: %v", r.Items)
		return
	}

More examples can be found in the SDK Github repo: https://github.com/oracle/oci-go-sdk/tree/master/example

Optional Fields in the SDK

Optional fields are represented with the `mandatory:"false"` tag on input structs. The SDK will omit all optional fields that are nil when making requests.
In the case of enum-type fields, the SDK will omit fields whose value is an empty string.

Helper Functions

The SDK uses pointers for primitive types in many input structs. To aid in the construction of such structs, the SDK provides
functions that return a pointer for a given value. For example:

	// Given the struct
	type CreateVcnDetails struct {

		// Example: `172.16.0.0/16`
		CidrBlock *string `mandatory:"true" json:"cidrBlock"`

		CompartmentId *string `mandatory:"true" json:"compartmentId"`

		DisplayName *string `mandatory:"false" json:"displayName"`

	}

	// We can use the helper functions to build the struct
	details := core.CreateVcnDetails{
		CidrBlock:     common.String("172.16.0.0/16"),
		CompartmentId: common.String("someOcid"),
		DisplayName:   common.String("myVcn"),
	}


Customizing Requests

The SDK exposes functionality that allows the user to customize any http request before is sent to the service.

You can do so by setting the `Interceptor` field in any of the `Client` structs. For example:

	client, err := audit.NewAuditClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		panic(err)
	}

	// This will add a header called "X-CustomHeader" to all request
	// performed with client
	client.Interceptor = func(request *http.Request) error {
		request.Header.Set("X-CustomHeader", "CustomValue")
		return nil
	}

The Interceptor closure gets called before the signing process, thus any changes done to the request will be properly
signed and submitted to the service.


Signing Custom Requests

The SDK exposes a stand-alone signer that can be used to signing custom requests. Related code can be found here:
https://github.com/oracle/oci-go-sdk/blob/master/common/http_signer.go.

The example below shows how to create a default signer.

	client := http.Client{}
	var request http.Request
	request = ... // some custom request

	// Set the Date header
	request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	// And a provider of cryptographic keys
	provider := common.DefaultConfigProvider()

	// Build the signer
	signer := common.DefaultSigner(provider)

	// Sign the request
	signer.Sign(&request)

	// Execute the request
	client.Do(request)



The signer also allows more granular control on the headers used for signing. For example:

	client := http.Client{}
	var request http.Request
	request = ... // some custom request

	// Set the Date header
	request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	// Mandatory headers to be used in the sign process
	defaultGenericHeaders    = []string{"date", "(request-target)", "host"}

	// Optional headers
	optionalHeaders = []string{"content-length", "content-type", "x-content-sha256"}

	// A predicate that specifies when to use the optional signing headers
	optionalHeadersPredicate := func (r *http.Request) bool {
		return r.Method == http.MethodPost
	}

	// And a provider of cryptographic keys
	provider := common.DefaultConfigProvider()

	// Build the signer
	signer := common.RequestSigner(provider, defaultGenericHeaders, optionalHeaders, optionalHeadersPredicate)

	// Sign the request
	signer.Sign(&request)

	// Execute the request
	c.Do(request)

You can combine a custom signer with the exposed clients in the SDK.
This allows you to add custom signed headers to the request. Following is an example:

	//Create a provider of cryptographic keys
	provider := common.DefaultConfigProvider()

	//Create a client for the service you interested in
	c, _ := identity.NewIdentityClientWithConfigurationProvider(provider)

	//Define a custom header to be signed, and add it to the list of default headers
	customHeader := "opc-my-token"
	allHeaders := append(common.DefaultGenericHeaders(), customHeader)

	//Overwrite the signer in your client to sign the new slice of headers
	c.Signer = common.RequestSigner(provider, allHeaders, common.DefaultBodyHeaders())

	//Set the value of the header. This can be done with an Interceptor
	c.Interceptor = func(request *http.Request) error {
		request.Header.Add(customHeader, "customvalue")
		return nil
	}

	//Execute your operation as before
	c.ListGroups(..)

Bear in mind that some services have a white list of headers that it expects to be signed.
Therefore, adding an arbitrary header can result in authentications errors.
To see a runnable example, see https://github.com/oracle/oci-go-sdk/blob/master/example/example_identity_test.go


For more information on the signing algorithm refer to: https://docs.cloud.oracle.com/Content/API/Concepts/signingrequests.htm

Polymorphic JSON Requests and Responses

Some operations accept or return polymorphic JSON objects. The SDK models such objects as interfaces. Further the SDK provides
structs that implement such interfaces. Thus, for all operations that expect interfaces as input, pass the struct in the SDK that satisfies
such interface. For example:

	c, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		panic(err)
	}

	// The CreateIdentityProviderRequest takes a CreateIdentityProviderDetails interface as input
	rCreate := identity.CreateIdentityProviderRequest{}

	// The CreateSaml2IdentityProviderDetails struct implements the CreateIdentityProviderDetails interface
	details := identity.CreateSaml2IdentityProviderDetails{}
	details.CompartmentId = common.String(getTenancyID())
	details.Name = common.String("someName")
	//... more setup if needed
	// Use the above struct
	rCreate.CreateIdentityProviderDetails = details

	// Make the call
	rspCreate, createErr := c.CreateIdentityProvider(context.Background(), rCreate)

In the case of a polymorphic response you can type assert the interface to the expected type. For example:

	rRead := identity.GetIdentityProviderRequest{}
	rRead.IdentityProviderId = common.String("aValidId")
	response, err := c.GetIdentityProvider(context.Background(), rRead)

	provider := response.IdentityProvider.(identity.Saml2IdentityProvider)

An example of polymorphic JSON request handling can be found here: https://github.com/oracle/oci-go-sdk/blob/master/example/example_core_test.go#L63


Pagination

When calling a list operation, the operation will retrieve a page of results. To retrieve more data, call the list operation again,
passing in the value of the most recent response's OpcNextPage as the value of Page in the next list operation call.
When there is no more data the OpcNextPage field will be nil. An example of pagination using this logic can be found here: https://github.com/oracle/oci-go-sdk/blob/master/example/example_core_pagination_test.go

Logging and Debugging

The SDK has a built-in logging mechanism used internally. The internal logging logic is used to record the raw http
requests, responses and potential errors when (un)marshalling request and responses.

Built-in logging in the SDK is controlled via the environment variable "OCI_GO_SDK_DEBUG" and its contents. The below are possible values for the "OCI_GO_SDK_DEBUG" variable

1. "info" or "i" enables all info logging messages

2. "debug" or "d"  enables all debug and info logging messages

3. "verbose" or "v" or "1" enables all verbose, debug and info logging messages

4. "null" turns all logging messages off.

If the value of the environment variable does not match any of the above then default logging level is "info".
If the environment variable is not present then no logging messages are emitted.


Retry

Sometimes you may need to wait until an attribute of a resource, such as an instance or a VCN, reaches a certain state.
An example of this would be launching an instance and then waiting for the instance to become available, or waiting until a subnet in a VCN has been terminated.
You might also want to retry the same operation again if there's network issue etc...
This can be accomplished by using the RequestMetadata.RetryPolicy. You can find the examples here: https://github.com/oracle/oci-go-sdk/blob/master/example/example_retry_test.go

Using the SDK with a Proxy Server

The GO SDK uses the net/http package to make calls to OCI services. If your environment requires you to use a proxy server for outgoing HTTP requests
then you can set this up in the following ways:

1. Configuring environment variable as described here https://golang.org/pkg/net/http/#ProxyFromEnvironment
2. Modifying the underlying Transport struct for a service client

In order to modify the underlying Transport struct in HttpClient, you can do something similar to (sample code for audit service client):
	// create audit service client
	client, clerr := audit.NewAuditClientWithConfigurationProvider(common.DefaultConfigProvider())

	// create a proxy url
	proxyURL, err := url.Parse("http(s)://[username]:[password]@[ip address]:[port]")

	client.HTTPClient = &http.Client{
		// adding the proxy settings to the http.Transport
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}


Uploading Large Objects

The Object Storage service supports multipart uploads to make large object uploads easier by splitting the large object into parts. The Go SDK supports raw multipart upload operations for advanced use cases, as well as a higher level upload class that uses the multipart upload APIs. For links to the APIs used for multipart upload operations, see Managing Multipart Uploads (https://docs.cloud.oracle.com/iaas/Content/Object/Tasks/usingmultipartuploads.htm). Higher level multipart uploads are implemented using the UploadManager, which will: split a large object into parts for you, upload the parts in parallel, and then recombine and commit the parts as a single object in storage.

This code sample shows how to use the UploadManager to automatically split an object into parts for upload to simplify interaction with the Object Storage service: https://github.com/oracle/oci-go-sdk/blob/master/example/example_objectstorage_test.go


Forward Compatibility

Some response fields are enum-typed. In the future, individual services may return values not covered by existing enums
for that field. To address this possibility, every enum-type response field is a modeled as a type that supports any string.
Thus if a service returns a value that is not recognized by your version of the SDK, then the response field will be set to this value.

When individual services return a polymorphic JSON response not available as a concrete struct, the SDK will return an implementation that only satisfies
the interface modeling the polymorphic JSON response.


New Region Support

If you are using a version of the SDK released prior to the announcement of a new region, you may need to use a workaround to reach it, depending on whether the region is in the oraclecloud.com realm.

A region is a localized geographic area. For more information on regions and how to identify them, see Regions and Availability Domains(https://docs.cloud.oracle.com/iaas/Content/General/Concepts/regions.htm).

A realm is a set of regions that share entities. You can identify your realm by looking at the domain name at the end of the network address. For example, the realm for xyz.abc.123.oraclecloud.com is oraclecloud.com.

oraclecloud.com Realm: For regions in the oraclecloud.com realm, even if common.Region does not contain the new region, the forward compatibility of the SDK can automatically handle it. You can pass new region names just as you would pass ones that are already defined. For more information on passing region names in the configuration, see Configuring (https://github.com/oracle/oci-go-sdk/blob/master/README.md#configuring). For details on common.Region, see (https://github.com/oracle/oci-go-sdk/blob/master/common/common.go).

Other Realms: For regions in realms other than oraclecloud.com, you can use the following workarounds to reach new regions with earlier versions of the SDK.

NOTE: Be sure to supply the appropriate endpoints for your region.

You can overwrite the target host with client.Host:
	client.Host = 'https://identity.us-gov-phoenix-1.oraclegovcloud.com'

If you are authenticating via instance principals, you can set the authentication endpoint in an environment variable:
	export OCI_SDK_AUTH_CLIENT_REGION_URL="https://identity.us-gov-phoenix-1.oraclegovcloud.com"


Contributions

Got a fix for a bug, or a new feature you'd like to contribute? The SDK is open source and accepting pull requests on GitHub
https://github.com/oracle/oci-go-sdk

License

Licensing information available at: https://github.com/oracle/oci-go-sdk/blob/master/LICENSE.txt

Notifications

To be notified when a new version of the Go SDK is released, subscribe to the following feed: https://github.com/oracle/oci-go-sdk/releases.atom

Questions or Feedback

Please refer to this link: https://github.com/oracle/oci-go-sdk#help




*/
package oci

//go:generate go run cmd/genver/main.go cmd/genver/version_template.go --output common/version.go
