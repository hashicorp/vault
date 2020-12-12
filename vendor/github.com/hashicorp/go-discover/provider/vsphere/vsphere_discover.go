// Package vsphere provides node discovery for VMware vSphere.
//
// The package performs discovery by searching vCenter for all nodes matching a
// certain tag, it then discovers all known IP addresses through VMware tools
// that are not loopback or auto-configuration addresses.
//
// This package requires at least vSphere 6.0 in order to function.
package vsphere

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vic/pkg/vsphere/tags"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// providerLog is the local provider logger. This should be initialized from
// the provider entry point.
var logger *log.Logger

// setLog sets the logger.
func setLog(l *log.Logger) {
	if l != nil {
		logger = l
	} else {
		logger = log.New(ioutil.Discard, "", 0)
	}
}

// discoverErr prints out a friendly error heading for the top-level discovery
// errors. It should only be used in the Addrs method.
func discoverErr(format string, a ...interface{}) error {
	var s string
	if len(a) > 1 {
		s = fmt.Sprintf(format, a...)
	} else {
		s = format
	}
	return fmt.Errorf("discover-vsphere: %s", s)
}

// valueOrEnv provides a way of suppling configuration values through
// environment variables. Defined values always take priority.
func valueOrEnv(config map[string]string, key, env string) string {
	if v := config[key]; v != "" {
		return v
	}
	if v := os.Getenv(env); v != "" {
		logger.Printf("[DEBUG] Using value of %s for configuration of %s", env, key)
		return v
	}
	return ""
}

// vSphereClient is a client connection manager for the vSphere provider.
type vSphereClient struct {
	// The VIM/govmomi client.
	VimClient *govmomi.Client

	// The specialized tags client SDK imported from vmware/vic.
	TagsClient *tags.RestClient
}

// vimURL returns a URL to pass to the VIM SOAP client.
func vimURL(server, user, password string) (*url.URL, error) {
	u, err := url.Parse("https://" + server + "/sdk")
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %s", err)
	}

	u.User = url.UserPassword(user, password)

	return u, nil
}

// newVSphereClient returns a new vSphereClient after setting up the necessary
// connections.
func newVSphereClient(ctx context.Context, host, user, password string, insecure bool) (*vSphereClient, error) {
	logger.Println("[DEBUG] Connecting to vSphere client endpoints")

	client := new(vSphereClient)

	u, err := vimURL(host, user, password)
	if err != nil {
		return nil, fmt.Errorf("error generating SOAP endpoint url: %s", err)
	}

	// Set up the VIM/govmomi client connection
	client.VimClient, err = newVimSession(ctx, u, insecure)
	if err != nil {
		return nil, err
	}

	client.TagsClient, err = newRestSession(ctx, u, insecure)
	if err != nil {
		return nil, err
	}

	logger.Println("[DEBUG] All vSphere client endpoints connected successfully")
	return client, nil
}

// newVimSession connects the VIM SOAP API client connection.
func newVimSession(ctx context.Context, u *url.URL, insecure bool) (*govmomi.Client, error) {
	logger.Printf("[DEBUG] Creating new SOAP API session on endpoint %s", u.Host)
	client, err := govmomi.NewClient(ctx, u, insecure)
	if err != nil {
		return nil, fmt.Errorf("error setting up new vSphere SOAP client: %s", err)
	}

	logger.Println("[DEBUG] SOAP API session creation successful")
	return client, nil
}

// newRestSession connects to the vSphere REST API endpoint, necessary for
// tags.
func newRestSession(ctx context.Context, u *url.URL, insecure bool) (*tags.RestClient, error) {
	logger.Printf("[DEBUG] Creating new CIS REST API session on endpoint %s", u.Host)
	client := tags.NewClient(u, insecure, "")
	if err := client.Login(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to CIS REST endpoint: %s", err)
	}

	logger.Println("[DEBUG] CIS REST API session creation successful")
	return client, nil
}

// Provider defines the vSphere discovery provider.
type Provider struct{}

// Help implements the Provider interface for the vsphere package.
func (p *Provider) Help() string {
	return `VMware vSphere:

    provider:      "vsphere"
    tag_name:      The name of the tag to look up.
    category_name: The category of the tag to look up.
    host:          The host of the vSphere server to connect to.
    user:          The username to connect as.
    password:      The password of the user to connect to vSphere as.
    insecure_ssl:  Whether or not to skip SSL certificate validation.
    timeout:       Discovery context timeout (default: 10m)
`
}

// Addrs implements the Provider interface for the vsphere package.
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "vsphere" {
		return nil, discoverErr("invalid provider %s", args["provider"])
	}

	setLog(l)

	tagName := args["tag_name"]
	categoryName := args["category_name"]
	host := valueOrEnv(args, "host", "VSPHERE_SERVER")
	user := valueOrEnv(args, "user", "VSPHERE_USER")
	password := valueOrEnv(args, "password", "VSPHERE_PASSWORD")
	insecure, err := strconv.ParseBool(valueOrEnv(args, "insecure_ssl", "VSPHERE_ALLOW_UNVERIFIED_SSL"))
	if err != nil {
		logger.Println("[DEBUG] Non-truthy/falsey value for insecure_ssl, assuming false")
	}
	timeout, err := time.ParseDuration(args["timeout"])
	if err != nil {
		logger.Println("[DEBUG] Non-time value given for timeout, assuming 10m")
		timeout = time.Minute * 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := newVSphereClient(ctx, host, user, password, insecure)
	if err != nil {
		return nil, discoverErr(err.Error())
	}

	if tagName == "" || categoryName == "" {
		return nil, discoverErr("both tag_name and category_name must be specified")
	}

	logger.Printf("[INFO] Locating all virtual machine IP addresses with tag %q in category %q", tagName, categoryName)

	tagID, err := tagIDFromName(ctx, client.TagsClient, tagName, categoryName)
	if err != nil {
		return nil, discoverErr(err.Error())
	}

	addrs, err := virtualMachineIPsForTag(ctx, client, tagID)
	if err != nil {
		return nil, discoverErr(err.Error())
	}

	logger.Printf("[INFO] Final IP address list: %s", strings.Join(addrs, ","))
	return addrs, nil
}

// tagIDFromName helps convert the tag and category names into the final ID
// used for discovery.
func tagIDFromName(ctx context.Context, client *tags.RestClient, name, category string) (string, error) {
	logger.Printf("[DEBUG] Fetching tag ID for tag name %q and category %q", name, category)

	categoryID, err := tagCategoryByName(ctx, client, category)
	if err != nil {
		return "", err
	}

	return tagByName(ctx, client, name, categoryID)
}

// tagCategoryByName converts a tag category name into its ID.
func tagCategoryByName(ctx context.Context, client *tags.RestClient, name string) (string, error) {
	cats, err := client.GetCategoriesByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("could not get category for name %q: %s", name, err)
	}

	if len(cats) < 1 {
		return "", fmt.Errorf("category name %q not found", name)
	}
	if len(cats) > 1 {
		// Although GetCategoriesByName does not seem to think that tag categories
		// are unique, empirical observation via the console and API show that they
		// are. This error case is handled anyway.
		return "", fmt.Errorf("multiple categories with name %q found", name)
	}

	return cats[0].ID, nil
}

// tagByName converts a tag name into its ID.
func tagByName(ctx context.Context, client *tags.RestClient, name, categoryID string) (string, error) {
	tids, err := client.GetTagByNameForCategory(ctx, name, categoryID)
	if err != nil {
		return "", fmt.Errorf("could not get tag for name %q: %s", name, err)
	}

	if len(tids) < 1 {
		return "", fmt.Errorf("tag name %q not found in category ID %q", name, categoryID)
	}
	if len(tids) > 1 {
		// This situation is very similar to the one in tagCategoryByName. The API
		// docs even say that tags need to be unique in categories, yet
		// GetTagByNameForCategory still returns multiple results.
		return "", fmt.Errorf("multiple tags with name %q found", name)
	}

	logger.Printf("[DEBUG] Tag ID is %q", tids[0].ID)
	return tids[0].ID, nil
}

// virtualMachineIPsForTag is a higher-level wrapper that calls out to
// functions to fetch all of the virtual machines matching a certain tag ID,
// and then gets all of the IP addresses for those virtual machines.
func virtualMachineIPsForTag(ctx context.Context, client *vSphereClient, id string) ([]string, error) {
	vms, err := virtualMachinesForTag(ctx, client, id)
	if err != nil {
		return nil, err
	}

	return ipAddrsForVirtualMachines(ctx, client, vms)
}

// virtualMachinesForTag discovers all of the virtual machines that match a
// specific tag ID and returns their higher level helper objects.
func virtualMachinesForTag(ctx context.Context, client *vSphereClient, id string) ([]*object.VirtualMachine, error) {
	logger.Printf("[DEBUG] Locating all virtual machines under tag ID %q", id)

	var vms []*object.VirtualMachine

	objs, err := client.TagsClient.ListAttachedObjects(ctx, id)
	if err != nil {
		return nil, err
	}
	for i, obj := range objs {
		switch {
		case obj.Type == nil || obj.ID == nil:
			logger.Printf("[WARN] Discovered object at index %d has either no ID or type", i)
			continue
		case *obj.Type != "VirtualMachine":
			logger.Printf("[DEBUG] Discovered object ID %q is not a virutal machine", *obj.ID)
			continue
		}
		vm, err := virtualMachineFromMOID(ctx, client.VimClient, *obj.ID)
		if err != nil {
			return nil, fmt.Errorf("error locating virtual machine with ID %q: %s", *obj.ID, err)
		}
		vms = append(vms, vm)
	}

	logger.Printf("[DEBUG] Discovered virtual machines: %s", virtualMachineNames(vms))
	return vms, nil
}

// ipAddrsForVirtualMachines takes a set of virtual machines and returns a
// consolidated list of IP addresses for all of the VMs.
func ipAddrsForVirtualMachines(ctx context.Context, client *vSphereClient, vms []*object.VirtualMachine) ([]string, error) {
	var addrs []string
	for _, vm := range vms {
		as, err := buildAndSelectGuestIPs(ctx, vm)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, as...)
	}
	return addrs, nil
}

// virtualMachineFromMOID locates a virtual machine by its managed object
// reference ID.
func virtualMachineFromMOID(ctx context.Context, client *govmomi.Client, id string) (*object.VirtualMachine, error) {
	logger.Printf("[DEBUG] Locating VM with managed object ID %q", id)

	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "VirtualMachine",
		Value: id,
	}

	vm, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	// Should be safe to return here. If our reference returned here and is not a
	// VM, then we have bigger problems and to be honest we should be panicking
	// anyway.
	return vm.(*object.VirtualMachine), nil
}

// virtualMachineProperties is a convenience method that wraps fetching the
// VirtualMachine MO from its higher-level object.
//
// It takes a list of property keys to fetch. Keeping the property set small
// can sometimes result in significant performance improvements.
func virtualMachineProperties(ctx context.Context, vm *object.VirtualMachine, keys []string) (*mo.VirtualMachine, error) {
	logger.Printf("[DEBUG] Fetching properties for VM %q", vm.Name())
	var props mo.VirtualMachine
	if err := vm.Properties(ctx, vm.Reference(), keys, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

// buildAndSelectGuestIPs builds a list of IP addresses known to VMware tools,
// skipping local and auto-configuration addresses.
//
// The builder is non-discriminate and is only deterministic to the order that
// it discovers addresses in VMware tools.
func buildAndSelectGuestIPs(ctx context.Context, vm *object.VirtualMachine) ([]string, error) {
	logger.Printf("[DEBUG] Discovering addresses for virtual machine %q", vm.Name())
	var addrs []string

	props, err := virtualMachineProperties(ctx, vm, []string{"guest.net"})
	if err != nil {
		return nil, fmt.Errorf("cannot fetch properties for VM %q: %s", vm.Name(), err)
	}

	if props.Guest == nil || props.Guest.Net == nil {
		logger.Printf("[WARN] No networking stack information available for %q or VMware tools not running", vm.Name())
		return nil, nil
	}

	// Now fetch all IP addresses, checking at the same time to see if the IP
	// address is eligible to be a primary IP address.
	for _, n := range props.Guest.Net {
		if n.IpConfig != nil {
			for _, addr := range n.IpConfig.IpAddress {
				if skipIPAddr(net.ParseIP(addr.IpAddress)) {
					continue
				}
				addrs = append(addrs, addr.IpAddress)
			}
		}
	}

	logger.Printf("[INFO] Discovered IP addresses for virtual machine %q: %s", vm.Name(), strings.Join(addrs, ","))
	return addrs, nil
}

// skipIPAddr defines the set of criteria that buildAndSelectGuestIPs uses to
// check to see if it needs to skip an IP address.
func skipIPAddr(ip net.IP) bool {
	switch {
	case ip.IsLinkLocalMulticast():
		fallthrough
	case ip.IsLinkLocalUnicast():
		fallthrough
	case ip.IsLoopback():
		fallthrough
	case ip.IsMulticast():
		return true
	}
	return false
}

// virtualMachineNames is a helper method that returns all the names for a list
// of virtual machines, comma separated.
func virtualMachineNames(vms []*object.VirtualMachine) string {
	var s []string
	for _, vm := range vms {
		s = append(s, vm.Name())
	}
	return strings.Join(s, ",")
}
