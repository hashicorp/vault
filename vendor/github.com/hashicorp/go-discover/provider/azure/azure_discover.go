// Package azure provides node discovery for Microsoft Azure.
package azure

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2015-06-15/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Microsoft Azure:

   provider:          "azure"
   tenant_id:         The id of the tenant
   client_id:         The id of the client
   subscription_id:   The id of the subscription
   secret_access_key: The authentication credential
    **NOTE** The secret_access_key value often may have an equals sign in it's value,
    especially if generated from the Azure Portal. So is important to wrap in single quotes
    eg. secret_acccess_key='fpOfcHQJAQBczjAxiVpeyLmX1M0M0KPBST+GU2GvEN4='

   Variables can also be provided by environmental variables:
    export ARM_SUBSCRIPTION_ID for subscription
    export ARM_TENANT_ID for tenant
    export ARM_CLIENT_ID for client
    export ARM_CLIENT_SECRET for secret access key

   If none of those options are given, the Azure SDK is using the default  environment based authentication outlined
   here https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization#use-environment-based-authentication
   This will fallback to MSI if nothing is explicitly specified.

   Use these configuration parameters when using tags:

   tag_name:          The name of the tag to filter on
   tag_value:         The value of the tag to filter on

   Use these configuration parameters when using Virtual Machine Scale Sets:

   resource_group:    The name of the resource group to filter on
   vm_scale_set:      The name of the virtual machine scale set to filter on

   When using tags the only permission needed is Microsoft.Network/networkInterfaces/*

   When using Virtual Machine Scale Sets the only role action needed is Microsoft.Compute/virtualMachineScaleSets/*/read.

   It is recommended you make a dedicated key used only for auto-joining.
`
}

// argsOrEnv allows you to pick an environmental variable for a setting if the arg is not set
func argsOrEnv(args map[string]string, key, env string) string {
	if value, ok := args[key]; ok {
		return value
	}
	return os.Getenv(env)
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	var authorizer autorest.Authorizer

	if args["provider"] != "azure" {
		return nil, fmt.Errorf("discover-azure: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	// check for environmental variables, and use if the argument hasn't been set in config
	tenantID := argsOrEnv(args, "tenant_id", "ARM_TENANT_ID")
	clientID := argsOrEnv(args, "client_id", "ARM_CLIENT_ID")
	subscriptionID := argsOrEnv(args, "subscription_id", "ARM_SUBSCRIPTION_ID")
	secretKey := argsOrEnv(args, "secret_access_key", "ARM_CLIENT_SECRET")

	// Try to use the argument and environment provided arguments first, if this fails fall back to the Azure
	// SDK provided methods
	if tenantID != "" && clientID != "" && secretKey != "" {
		var err error
		authorizer, err = auth.NewClientCredentialsConfig(clientID, secretKey, tenantID).Authorizer()
		if err != nil {
			return nil, fmt.Errorf("discover-azure (ClientCredentials): %s", err)
		}
	} else {
		var err error
		authorizer, err = auth.NewAuthorizerFromEnvironment()
		if err != nil {
			return nil, fmt.Errorf("discover-azure (EnvironmentCredentials): %s", err)
		}
	}

	// Use tags if using network interfaces
	tagName := args["tag_name"]
	tagValue := args["tag_value"]

	// Use resourceGroup and vmScaleSet if using vm scale sets
	resourceGroup := args["resource_group"]
	vmScaleSet := args["vm_scale_set"]

	// Setup the client using autorest; followed the structure from Terraform
	vmnet := network.NewInterfacesClient(subscriptionID)
	vmnet.Sender = autorest.CreateSender(autorest.WithLogging(l))
	vmnet.Authorizer = authorizer

	if p.userAgent != "" {
		vmnet.Client.UserAgent = p.userAgent
	}

	if tagName != "" && tagValue != "" && resourceGroup == "" && vmScaleSet == "" {
		l.Printf("[DEBUG] discover-azure: using tag method. tag_name: %s, tag_value: %s", tagName, tagValue)
		return fetchAddrsWithTags(tagName, tagValue, vmnet, l)
	} else if resourceGroup != "" && vmScaleSet != "" && tagName == "" && tagValue == "" {
		l.Printf("[DEBUG] discover-azure: using vm scale set method. resource_group: %s, vm_scale_set: %s", resourceGroup, vmScaleSet)
		return fetchAddrsWithVmScaleSet(resourceGroup, vmScaleSet, vmnet, l)
	} else {
		l.Printf("[ERROR] discover-azure: tag_name: %s, tag_value: %s", tagName, tagValue)
		l.Printf("[ERROR] discover-azure: resource_group %s, vm_scale_set %s", resourceGroup, vmScaleSet)
		return nil, fmt.Errorf("discover-azure: unclear configuration. use (tag name and value) or (resouce_group and vm_scale_set)")
	}

}

func fetchAddrsWithTags(tagName string, tagValue string, vmnet network.InterfacesClient, l *log.Logger) ([]string, error) {
	// Get all network interfaces across resource groups
	// unless there is a compelling reason to restrict

	ctx := context.Background()
	netres, err := vmnet.ListAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("discover-azure: %s", err)
	}

	if len(netres.Values()) == 0 {
		return nil, fmt.Errorf("discover-azure: no interfaces")
	}

	// Choose any PrivateIPAddress with the matching tag
	var addrs []string
	for _, v := range netres.Values() {
		var id string
		if v.ID != nil {
			id = *v.ID
		} else {
			id = "ip address id not found"
		}
		if v.Tags == nil {
			l.Printf("[DEBUG] discover-azure: Interface %s has no tags", id)
			continue
		}
		tv := v.Tags[tagName] // *string
		if tv == nil {
			l.Printf("[DEBUG] discover-azure: Interface %s did not have tag: %s", id, tagName)
			continue
		}
		if *tv != tagValue {
			l.Printf("[DEBUG] discover-azure: Interface %s tag value was: %s which did not match: %s", id, *tv, tagValue)
			continue
		}
		if v.IPConfigurations == nil {
			l.Printf("[DEBUG] discover-azure: Interface %s had no ip configuration", id)
			continue
		}
		for _, x := range *v.IPConfigurations {
			if x.PrivateIPAddress == nil {
				l.Printf("[DEBUG] discover-azure: Interface %s had no private ip", id)
				continue
			}
			iAddr := *x.PrivateIPAddress
			l.Printf("[DEBUG] discover-azure: Interface %s has private ip: %s", id, iAddr)
			addrs = append(addrs, iAddr)
		}
	}
	l.Printf("[DEBUG] discover-azure: Found ip addresses: %v", addrs)
	return addrs, nil
}

func fetchAddrsWithVmScaleSet(resourceGroup string, vmScaleSet string, vmnet network.InterfacesClient, l *log.Logger) ([]string, error) {
	// Get all network interfaces for a specific virtual machine scale set
	ctx := context.Background()
	netres, err := vmnet.ListVirtualMachineScaleSetNetworkInterfaces(ctx, resourceGroup, vmScaleSet)
	if err != nil {
		return nil, fmt.Errorf("discover-azure: %s", err)
	}

	if len(netres.Values()) == 0 {
		return nil, fmt.Errorf("discover-azure: no interfaces")
	}

	// Get all of PrivateIPAddresses we can.
	var addrs []string
	for _, v := range netres.Values() {
		var id string
		if v.ID != nil {
			id = *v.ID
		} else {
			id = "ip address id not found"
		}
		if v.IPConfigurations == nil {
			l.Printf("[DEBUG] discover-azure: Interface %s had no ip configuration", id)
			continue
		}
		for _, x := range *v.IPConfigurations {
			if x.PrivateIPAddress == nil {
				l.Printf("[DEBUG] discover-azure: Interface %s had no private ip", id)
				continue
			}
			iAddr := *x.PrivateIPAddress
			l.Printf("[DEBUG] discover-azure: Interface %s has private ip: %s", id, iAddr)
			addrs = append(addrs, iAddr)
		}
	}
	l.Printf("[DEBUG] discover-azure: Found ip addresses: %v", addrs)
	return addrs, nil
}
