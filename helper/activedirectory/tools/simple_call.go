package tools

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/activedirectory"
)

var (
	// ex. "ldap://138.91.247.105"
	rawURL = os.Getenv("TEST_LDAP_URL")
	dn     = os.Getenv("TEST_DN")

	// these can be left blank if the operation performed doesn't require them
	username = os.Getenv("TEST_LDAP_USERNAME")
	password = os.Getenv("TEST_LDAP_PASSWORD")
)

// main executes one call using a simple client pointed at the given instance.
func main() {

	config := newInsecureConfig()
	client := activedirectory.NewClient(hclog.Default(), config)

	filters := map[*activedirectory.Field][]string{
		activedirectory.FieldRegistry.GivenName: {"Sara", "Sarah"},
	}

	entries, err := client.Search(filters)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("found %d entries:\n", len(entries))
	for _, entry := range entries {
		fmt.Printf("%s\n", entry)
	}
}

func newInsecureConfig() *activedirectory.Configuration {
	return &activedirectory.Configuration{
		RootDomainName: dn,
		Certificate:    "",
		InsecureTLS:    true,
		Password:       password,
		StartTLS:       false,
		TLSMinVersion:  771,
		TLSMaxVersion:  771,
		URLs:           []string{rawURL},
		Username:       username,
	}
}
