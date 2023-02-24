package server

import (
	"os"
	"testing"

	"github.com/go-test/deep"
	sdkResource "github.com/hashicorp/hcp-sdk-go/resource"
	"github.com/hashicorp/vault/internalshared/configutil"
)

func TestHCPLinkConfig(t *testing.T) {
	os.Unsetenv("HCP_CLIENT_ID")
	os.Unsetenv("HCP_CLIENT_SECRET")
	os.Unsetenv("HCP_RESOURCE_ID")

	config, err := LoadConfigFile("./test-fixtures/hcp_link_config.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resIDRaw := "organization/bc58b3d0-2eab-4ab8-abf4-f61d3c9975ff/project/1c78e888-2142-4000-8918-f933bbbc7690/hashicorp.example.resource/example"
	res, _ := sdkResource.FromString(resIDRaw)

	expected := &Config{
		Storage: &Storage{
			Type:   "inmem",
			Config: map[string]string{},
		},
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:8200",
					TLSDisable:            true,
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},
			HCPLinkConf: &configutil.HCPLinkConfig{
				ResourceIDRaw: resIDRaw,
				Resource:      &res,
				ClientID:      "J2TtcSYOyPUkPV2z0mSyDtvitxLVjJmu",
				ClientSecret:  "N9JtHZyOnHrIvJZs82pqa54vd4jnkyU3xCcqhFXuQKJZZuxqxxbP1xCfBZVB82vY",
			},
			DisableMlock: true,
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}
