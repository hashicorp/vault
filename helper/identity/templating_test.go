package identity

import (
	"errors"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
)

func TestPopulate_Basic(t *testing.T) {
	var tests = []struct {
		name              string
		input             string
		output            string
		err               error
		entityName        string
		metadata          map[string]string
		aliasAccessor     string
		aliasID           string
		aliasName         string
		nilEntity         bool
		validityCheckOnly bool
		aliasMetadata     map[string]string
		groupName         string
		groupMetadata     map[string]string
	}{
		{
			name:   "no_templating",
			input:  "path foobar {",
			output: "path foobar {",
		},
		{
			name:  "only_closing",
			input: "path foobar}} {",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:  "closing_in_front",
			input: "path }} {{foobar}} {",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:  "closing_in_back",
			input: "path {{foobar}} }}",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:   "basic",
			input:  "path /{{identity.entity.id}}/ {",
			output: "path /entityID/ {",
		},
		{
			name:       "multiple",
			input:      "path {{identity.entity.name}} {\n\tval = {{identity.entity.metadata.foo}}\n}",
			entityName: "entityName",
			metadata:   map[string]string{"foo": "bar"},
			output:     "path entityName {\n\tval = bar\n}",
		},
		{
			name:     "multiple_bad_name",
			input:    "path {{identity.entity.name}} {\n\tval = {{identity.entity.metadata.foo}}\n}",
			metadata: map[string]string{"foo": "bar"},
			err:      ErrTemplateValueNotFound,
		},
		{
			name:  "unbalanced_close",
			input: "path {{identity.entity.id}} {\n\tval = {{ent}}ity.metadata.foo}}\n}",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:  "unbalanced_open",
			input: "path {{identity.entity.id}} {\n\tval = {{ent{{ity.metadata.foo}}\n}",
			err:   ErrUnbalancedTemplatingCharacter,
		},
		{
			name:      "no_entity_no_directives",
			input:     "path {{identity.entity.id}} {\n\tval = {{ent{{ity.metadata.foo}}\n}",
			err:       ErrNoEntityAttachedToToken,
			nilEntity: true,
		},
		{
			name:      "no_entity_no_diretives",
			input:     "path name {\n\tval = foo\n}",
			output:    "path name {\n\tval = foo\n}",
			nilEntity: true,
		},
		{
			name:          "alias_id_name",
			input:         "path {{ identity.entity.name}} {\n\tval = {{identity.entity.aliases.foomount.id}}\n}",
			entityName:    "entityName",
			aliasAccessor: "foomount",
			aliasID:       "aliasID",
			metadata:      map[string]string{"foo": "bar"},
			output:        "path entityName {\n\tval = aliasID\n}",
		},
		{
			name:          "alias_id_name_bad_selector",
			input:         "path foobar {\n\tval = {{identity.entity.aliases.foomount}}\n}",
			aliasAccessor: "foomount",
			err:           errors.New("invalid alias selector"),
		},
		{
			name:          "alias_id_name_bad_accessor",
			input:         "path \"foobar\" {\n\tval = {{identity.entity.aliases.barmount.id}}\n}",
			aliasAccessor: "foomount",
			err:           errors.New("alias not found"),
		},
		{
			name:          "alias_id_name",
			input:         "path \"{{identity.entity.name}}\" {\n\tval = {{identity.entity.aliases.foomount.metadata.zip}}\n}",
			entityName:    "entityName",
			aliasAccessor: "foomount",
			aliasID:       "aliasID",
			metadata:      map[string]string{"foo": "bar"},
			aliasMetadata: map[string]string{"zip": "zap"},
			output:        "path \"entityName\" {\n\tval = zap\n}",
		},
		{
			name:       "group_name",
			input:      "path \"{{identity.groups.ids.groupID.name}}\" {\n\tval = {{identity.entity.name}}\n}",
			entityName: "entityName",
			groupName:  "groupName",
			output:     "path \"groupName\" {\n\tval = entityName\n}",
		},
		{
			name:       "group_bad_id",
			input:      "path \"{{identity.groups.ids.hroupID.name}}\" {\n\tval = {{identity.entity.name}}\n}",
			entityName: "entityName",
			groupName:  "groupName",
			err:        errors.New("entity is not a member of group \"hroupID\""),
		},
		{
			name:       "group_id",
			input:      "path \"{{identity.groups.names.groupName.id}}\" {\n\tval = {{identity.entity.name}}\n}",
			entityName: "entityName",
			groupName:  "groupName",
			output:     "path \"groupID\" {\n\tval = entityName\n}",
		},
		{
			name:       "group_bad_name",
			input:      "path \"{{identity.groups.names.hroupName.id}}\" {\n\tval = {{identity.entity.name}}\n}",
			entityName: "entityName",
			groupName:  "groupName",
			err:        errors.New("entity is not a member of group \"hroupName\""),
		},
	}

	for _, test := range tests {
		var entity *Entity
		if !test.nilEntity {
			entity = &Entity{
				ID:       "entityID",
				Name:     test.entityName,
				Metadata: test.metadata,
			}
		}
		if test.aliasAccessor != "" {
			entity.Aliases = []*Alias{
				&Alias{
					MountAccessor: test.aliasAccessor,
					ID:            test.aliasID,
					Name:          test.aliasName,
					Metadata:      test.aliasMetadata,
				},
			}
		}
		var groups []*Group
		if test.groupName != "" {
			groups = append(groups, &Group{
				ID:          "groupID",
				Name:        test.groupName,
				Metadata:    test.groupMetadata,
				NamespaceID: namespace.RootNamespace.ID,
			})
		}
		subst, out, err := PopulateString(&PopulateStringInput{
			ValidityCheckOnly: test.validityCheckOnly,
			String:            test.input,
			Entity:            entity,
			Groups:            groups,
			Namespace:         namespace.RootNamespace,
		})
		if err != nil {
			if test.err == nil {
				t.Fatalf("%s: expected success, got error: %v", test.name, err)
			}
			if err.Error() != test.err.Error() {
				t.Fatalf("%s: got error: %v", test.name, err)
			}
		}
		if out != test.output {
			t.Fatalf("%s: bad output: %s", test.name, out)
		}
		if err == nil && !subst && out != test.input {
			t.Fatalf("%s: bad subst flag", test.name)
		}
	}
}
