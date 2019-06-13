package identity

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
)

var testEntity = &Entity{
	ID:   "abc-123",
	Name: "Entity Name",
	Metadata: map[string]string{
		"color":         "green",
		"size":          "small",
		"non-printable": "\"\n\t",
	},
	Aliases: []*Alias{
		{
			Name: "aws_123",
			Metadata: map[string]string{
				"service": "ec2",
				"region":  "west",
			},
		},
	},
}

var testGroups = []*Group{
	{ID: "a08b0c02", Name: "g1"},
	{ID: "239bef91", Name: "g2"},
}

func Test_Compile(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expError bool
	}{
		{
			name: "valid parameters",
			template: `
			{
			    "id": "{{identity.entity.id}}",
			    "name": "{{identity.entity.name}}",
			    "all metadata": "{{identity.entity.metadata}}",
			    "one metadata key": "{{identity.entity.metadata.my_key}}",
			    "alias metadata": "{{identity.entity.aliases.aws_123.metadata}}",
			    "one alias metadata key": "{{identity.entity.aliases.aws_123.metadata.my_key}}",
			    "group names": "{{identity.entity.group_names}}",
			    "group ids": "{{identity.entity.group_ids}}",
			    "time now": "{{time.now}}",
			    "time plus": "{{time.now.plus.10s}}",
			    "time minus": "{{time.now.minus.1h}}"
			}`,
			expError: false,
		},
		{
			name:     "invalid json",
			template: `{"id": "{{identity.entity.id}}" }_`,
			expError: true,
		},
		{
			name: "unknown parameter",
			template: `
			{
			    "id": "{{identity.entity.whatisthis}}"
			}`,
			expError: true,
		},
		{
			name: "invalid time operation",
			template: `
			{
			    "time": "{{time.now.divide.4h}}"
			}`,
			expError: true,
		},
		{
			name: "invalid time duration",
			template: `
			{
			    "time": "{{time.now.plus.4th}}"
			}`,
			expError: true,
		},
	}

	for _, test := range tests {
		_, err := NewCompiledTemplate(test.template)
		if (err != nil) != test.expError {
			t.Fatalf("test %q: unexpected error result. got: %v", test.name, err)
		}
	}
}

func Test_Render(t *testing.T) {
	template := `
			{
			    "id": "{{identity.entity.id}}",
			    "name": "{{identity.entity.name}}",
			    "all metadata": "{{identity.entity.metadata}}",
			    "one metadata key": "{{identity.entity.metadata.color}}",
			    "one metadata key not found": "{{identity.entity.metadata.asldfk}}",
			    "alias metadata": "{{identity.entity.aliases.aws_123.metadata}}",
			    "alias not found metadata": "{{identity.entity.aliases.blahblah.metadata}}",
			    "one alias metadata key": "{{identity.entity.aliases.aws_123.metadata.service}}",
			    "one not found alias metadata key": "{{identity.entity.aliases.blahblah.metadata.service}}",
			    "group names": "{{identity.entity.group_names}}",
			    "group ids": "{{identity.entity.group_ids}}",
			    "time now": "{{time.now}}",
			    "time plus": "{{time.now.plus.10s}}",
			    "time minus": "{{time.now.minus.1h}}",
				"repeated and": {"nested element": "{{identity.entity.name}}"}
			}`

	expected := `
			{
			    "id": "abc-123",
			    "name": "Entity Name",
			    "all metadata": {"color":"green","non-printable":"\"\n\t","size":"small"},
			    "one metadata key": "green",
			    "one metadata key not found": "",
			    "alias metadata": {"region":"west","service":"ec2"},
			    "alias not found metadata": {},
			    "one alias metadata key": "ec2",
			    "one not found alias metadata key": "",
			    "group names": ["g1","g2"],
			    "group ids": ["a08b0c02","239bef91"],
			    "time now": <now>,
			    "time plus": <plus>,
			    "time minus": <minus>,
				"repeated and": {"nested element": "Entity Name"}
			}`

	now := time.Now().Unix()
	minus := now - (60 * 60)
	plus := now + 10

	expected = strings.ReplaceAll(expected, "<now>", strconv.FormatInt(now, 10))
	expected = strings.ReplaceAll(expected, "<plus>", strconv.FormatInt(plus, 10))
	expected = strings.ReplaceAll(expected, "<minus>", strconv.FormatInt(minus, 10))

	ct, err := NewCompiledTemplate(template)
	if err != nil {
		t.Fatal(err)
	}

	out, err := ct.Render(testEntity, testGroups)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(out, expected); diff != nil {
		t.Fatal(diff)
	}
}
