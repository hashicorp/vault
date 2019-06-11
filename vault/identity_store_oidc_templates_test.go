package vault

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/identity"
)

var testEntity = &identity.Entity{
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

var testGroups = []*identity.Group{
	{ID: "a08b0c02", Name: "g1"},
	{ID: "239bef91", Name: "g2"},
}

func Test_OIDC_TestABC(t *testing.T) {
	tpl := `
{
  "alias_claims": {
    "id": "{{identity.entity.id}}",
    "region": "west",
    "stuff": "{{identity.entity.metadata}}",
    "groups": "{{identity.entity.group_names}}",
    "group_ids": "{{identity.entity.group_ids}}",
    "notamatch": "{{identity.entity.nope}}"
  }
}
`
	pt, err := CompileTemplate(tpl)
	if err != nil {
		t.Fatal(err)
	}

	out, err := pt.Render(testEntity, testGroups)
	if err != nil {
		t.Fatal(err)
	}

	exp := `
{
  "alias_claims": {
    "id": "abc-123",
    "region": "west",
    "stuff": {"color":"green","size":"small"},
    "groups": ["g1","g2"],
    "group_ids": ["a08b0c02","239bef91"],
    "notamatch": "{{identity.entity.nope}}"
  }
}
`

	if diff := deep.Equal(out, exp); diff != nil {
		t.Fatal(diff)
	}
}

/*
func Test_OIDC_classifyParameters(t *testing.T) {
	tests := []struct {
		s         string
		entity    *identity.Entity
		groups    []*identity.Group
		expResult string
		expErr    bool
	}{
		{
			s:         "identity.entity.id",
			entity:    testEntity,
			expResult: `"abc-123"`,
		},
		{
			s:         "identity.entity.name",
			entity:    testEntity,
			expResult: `"Entity Name"`,
		},
		{
			s:         "identity.entity.metadata",
			entity:    testEntity,
			expResult: `{"color":"green","size":"small"}`,
		},
		{
			s:         "identity.entity.metadata.size",
			entity:    testEntity,
			expResult: `"small"`,
		},
		{
			s:         "identity.entity.aliases.aws_123.metadata.region",
			entity:    testEntity,
			expResult: `"west"`,
		},
		{
			s:         "identity.entity.aliases.aws_123.metadata.region",
			entity:    testEntity,
			expResult: `"west"`,
		},
		{
			s:         "identity.entity.group_names",
			entity:    testEntity,
			groups:    testGroups,
			expResult: `["g1","g2"]`,
		},
		{
			s:         "identity.entity.group_ids",
			entity:    testEntity,
			groups:    testGroups,
			expResult: `["a08b0c02","239bef91"]`,
		},
	}
	for _, tt := range tests {
		gotResult, err := classifyParameters(tt.entity, tt.groups, tt.s)
		if (err != nil) != tt.expErr {
			t.Errorf("classifyParameters() error = %v, expErr %v", err, tt.expErr)
			return
		}
		if gotResult != tt.expResult {
			t.Errorf("classifyParameters() = %v, want %v", gotResult, tt.expResult)
		}
	}
}
*/
