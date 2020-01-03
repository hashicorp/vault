package framework

import (
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestIdentityTemplating(t *testing.T) {
	sysView := &logical.StaticSystemView{
		EntityVal: &logical.Entity{
			ID:   "test-id",
			Name: "test",
			Aliases: []*logical.Alias{
				{
					ID:            "alias-id",
					Name:          "test alias",
					MountAccessor: "test_mount",
					MountType:     "secret",
					Metadata: map[string]string{
						"alias-metadata": "alias-metadata-value",
					},
				},
			},
			Metadata: map[string]string{
				"entity-metadata": "entity-metadata-value",
			},
		},
		GroupsVal: []*logical.Group{
			{
				ID:   "group1-id",
				Name: "group1",
				Metadata: map[string]string{
					"group-metadata": "group-metadata-value",
				},
			},
		},
	}

	tCases := []struct {
		tpl      string
		expected string
	}{
		{
			tpl:      "{{identity.entity.id}}",
			expected: "test-id",
		},
		{
			tpl:      "{{identity.entity.name}}",
			expected: "test",
		},
		{
			tpl:      "{{identity.entity.metadata.entity-metadata}}",
			expected: "entity-metadata-value",
		},
		{
			tpl:      "{{identity.entity.aliases.test_mount.id}}",
			expected: "alias-id",
		},
		{
			tpl:      "{{identity.entity.aliases.test_mount.id}}",
			expected: "alias-id",
		},
		{
			tpl:      "{{identity.entity.aliases.test_mount.name}}",
			expected: "test alias",
		},
		{
			tpl:      "{{identity.entity.aliases.test_mount.metadata.alias-metadata}}",
			expected: "alias-metadata-value",
		},
		{
			tpl:      "{{identity.groups.ids.group1-id.name}}",
			expected: "group1",
		},
		{
			tpl:      "{{identity.groups.names.group1.id}}",
			expected: "group1-id",
		},
		{
			tpl:      "{{identity.groups.names.group1.metadata.group-metadata}}",
			expected: "group-metadata-value",
		},
		{
			tpl:      "{{identity.groups.ids.group1-id.metadata.group-metadata}}",
			expected: "group-metadata-value",
		},
	}

	for _, tCase := range tCases {
		out, err := PopulateIdentityTemplate(tCase.tpl, "test", sysView)
		if err != nil {
			t.Fatal(err)
		}

		if out != tCase.expected {
			t.Fatalf("got %q, expected %q", out, tCase.expected)
		}
	}
}
