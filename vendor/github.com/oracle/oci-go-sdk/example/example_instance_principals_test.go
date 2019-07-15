package example

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/identity"
)

// ExampleInstancePrincipals lists the availability domains in your tenancy.
// Make sure you run this example from a instance with the right permissions. In this example
// the root compartment is read from the OCI_ROOT_COMPARTMENT_ID environment variable.
// More information on instance principals can be found here: https://docs.cloud.oracle.com/Content/Identity/Tasks/callingservicesfrominstances.htm
func ExampleInstancePrincipals() {

	provider, err := auth.InstancePrincipalConfigurationProvider()
	helpers.FatalIfError(err)

	tenancyID := helpers.RootCompartmentID()
	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: tenancyID,
	}

	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)
	// Override the region, this is an optional step.
	// the InstancePrincipalsConfigurationProvider defaults to the region
	// in which the compute instance is currently running
	client.SetRegion(string(common.RegionLHR))

	r, err := client.ListAvailabilityDomains(context.Background(), request)
	helpers.FatalIfError(err)

	log.Printf("list of available domains: %v", r.Items)
	fmt.Println("Done")

	// Output:
	// Done
}

// ExampleInstancePrincipalsWithCustomClient lists the availability domains in your tenancy.
// Similar to the example above, this example shows how to customize the client.
func ExampleInstancePrincipalsWithCustomClient() {

	// Just load the system cert pool for demonstration purposes.
	rootCaPool, err := x509.SystemCertPool()

	helpers.FatalIfError(err)

	provider, err := auth.InstancePrincipalConfigurationProviderWithCustomClient(func(dispatcher common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error) {
		client := dispatcher.(*http.Client)
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCaPool,
			},
		}
		return client, nil
	})

	tenancyID := helpers.RootCompartmentID()
	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: tenancyID,
	}

	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)

	r, err := client.ListAvailabilityDomains(context.Background(), request)
	helpers.FatalIfError(err)

	log.Printf("list of available domains: %v", r.Items)
	fmt.Println("Done")
}
