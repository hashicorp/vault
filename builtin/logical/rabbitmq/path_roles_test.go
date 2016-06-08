package rabbitmq

import (
	"testing"

	"github.com/hashicorp/vault/logical/framework"
)

type validateNameTestCase struct {
	Name string
	Fail bool
}

func TestRoles_validateName(t *testing.T) {
	cases := map[string]validateNameTestCase{
		"test name": {
			Name: "test",
		},
		"empty name": {
			Name: "",
			Fail: true,
		},
	}

	data := &framework.FieldData{
		Schema: rolesFields(),
	}
	for name, c := range cases {
		data.Raw = map[string]interface{}{
			"name": c.Name,
		}

		actual, err := validateName(data)
		if err != nil && !c.Fail {
			t.Error(err)
		}

		if c.Name != actual {
			t.Errorf("Fail: %s: expected %s, got %s", name, c.Name, actual)
		}
	}
}
