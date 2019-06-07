package vault

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/identity"
)

func Test_OIDC_TestABC(t *testing.T) {
	ent1 := identity.Entity{
		ID:   "abc-123",
		Name: "Entity Name",
		Metadata: map[string]string{
			"color": "green",
			"size":  "small",
		},
		Aliases: []*identity.Alias{
			{
				Name: "aws_123",
				Metadata: map[string]string{
					"region": "west",
				},
			},
		},
	}
	tpl := `
{
  "alias_claims": {
    "id": "{{identity.entity.id}}",
    "region": "west",
    "stuff": "{{identity.entity.metadata}}"
  }
}
`
	out, err := ABC(tpl, ent1)
	if err != nil {
		t.Fatal(err)
	}

	exp := `
{
  "alias_claims": {
    "id": "abc-123",
    "region": "west",
    "stuff": {"color":"green","size":"small"}
  }
}
`

	if diff := deep.Equal(out, exp); diff != nil {
		t.Fatal(diff)
	}
}

func Test_OIDC_classifyParameters(t *testing.T) {
	ent1 := identity.Entity{
		ID:   "abc-123",
		Name: "Entity Name",
		Metadata: map[string]string{
			"color": "green",
			"size":  "small",
		},
		Aliases: []*identity.Alias{
			{
				Name: "aws_123",
				Metadata: map[string]string{
					"region": "west",
				},
			},
		},
	}
	tests := []struct {
		s         string
		entity    identity.Entity
		expResult string
		expErr    bool
	}{
		{
			s:         "identity.entity.id",
			entity:    ent1,
			expResult: `"abc-123"`,
		},
		{
			s:         "identity.entity.name",
			entity:    ent1,
			expResult: `"Entity Name"`,
		},
		{
			s:         "identity.entity.metadata",
			entity:    ent1,
			expResult: `{"color":"green","size":"small"}`,
		},
		{
			s:         "identity.entity.metadata.size",
			entity:    ent1,
			expResult: `"small"`,
		},
		{
			s:         "identity.entity.aliases.aws_123.metadata.region",
			entity:    ent1,
			expResult: `"west"`,
		},
	}
	for _, tt := range tests {
		gotResult, err := classifyParameters(tt.entity, tt.s)
		if (err != nil) != tt.expErr {
			t.Errorf("classifyParameters() error = %v, expErr %v", err, tt.expErr)
			return
		}
		if gotResult != tt.expResult {
			t.Errorf("classifyParameters() = %v, want %v", gotResult, tt.expResult)
		}
	}
}
