package rabbitmq

import (
	"testing"

	"github.com/hashicorp/vault/logical/framework"
)

type validateLeasesTestCase struct {
	Lease    int
	LeaseMax int
	Fail     bool
}

func TestConfigLease_validateLeases(t *testing.T) {
	cases := map[string]validateLeasesTestCase{
		"Both lease and lease max": {
			Lease:    60 * 60,
			LeaseMax: 60 * 60,
		},
		"Just lease": {
			Lease:    60 * 60,
			LeaseMax: 0,
		},
		"No lease nor lease max": {
			Lease:    0,
			LeaseMax: 0,
			Fail:     true,
		},
	}

	data := &framework.FieldData{
		Schema: configFields(),
	}
	for name, c := range cases {
		data.Raw = map[string]interface{}{
			leaseLabel:    c.Lease,
			leaseMaxLabel: c.LeaseMax,
		}

		_, _, err := validateLeases(data)
		if err != nil && c.Fail {
			// This was expected
			continue
		} else if err != nil {
			// This was unexpected
			t.Errorf("Failed: %s", name)
		} else if err == nil && c.Fail {
			// This was unexpected
			t.Errorf("Failed to fail: %s", name)
		}
	}
}
