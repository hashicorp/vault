package identitytpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	ErrUnbalancedTemplatingCharacter = errors.New("unbalanced templating characters")
	ErrNoEntityAttachedToToken       = errors.New("string contains entity template directives but no entity was provided")
	ErrNoGroupsAttachedToToken       = errors.New("string contains groups template directives but no groups were provided")
	ErrTemplateValueNotFound         = errors.New("no value could be found for one of the template directives")
)

const (
	ACLTemplating = iota // must be the first value for backwards compatibility
	JSONTemplating
)

type PopulateStringInput struct {
	String            string
	ValidityCheckOnly bool
	Entity            *logical.Entity
	Groups            []*logical.Group
	NamespaceID       string
	Mode              int       // processing mode, ACLTemplate or JSONTemplating
	Now               time.Time // optional, defaults to current time

	templateHandler templateHandlerFunc
	groupIDs        []string
	groupNames      []string
}

// templateHandlerFunc allows generating string outputs based on data type, and
// different handlers can be used based on mode. For example in ACL mode, strings
// are emitted verbatim, but they're wrapped in double quotes for JSON mode. And
// some structures, like slices, might be rendered in one mode but prohibited in
// another.
type templateHandlerFunc func(interface{}, ...string) (string, error)

// aclTemplateHandler processes known parameter data types when operating
// in ACL mode.
func aclTemplateHandler(v interface{}, keys ...string) (string, error) {
	switch t := v.(type) {
	case string:
		if t == "" {
			return "", ErrTemplateValueNotFound
		}
		return t, nil
	case []string:
		return "", ErrTemplateValueNotFound
	case map[string]string:
		if len(keys) > 0 {
			val, ok := t[keys[0]]
			if ok {
				return val, nil
			}
		}
		return "", ErrTemplateValueNotFound
	}

	return "", fmt.Errorf("unknown type: %T", v)
}

// jsonTemplateHandler processes known parameter data types when operating
// in JSON mode.
func jsonTemplateHandler(v interface{}, keys ...string) (string, error) {
	jsonMarshaller := func(v interface{}) (string, error) {
		enc, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(enc), nil
	}

	switch t := v.(type) {
	case string:
		return strconv.Quote(t), nil
	case []string:
		return jsonMarshaller(t)
	case map[string]string:
		if len(keys) > 0 {
			return strconv.Quote(t[keys[0]]), nil
		}
		if t == nil {
			return "{}", nil
		}
		return jsonMarshaller(t)
	}

	return "", fmt.Errorf("unknown type: %T", v)
}

func PopulateString(p PopulateStringInput) (bool, string, error) {
	if p.String == "" {
		return false, "", nil
	}

	// preprocess groups
	for _, g := range p.Groups {
		p.groupNames = append(p.groupNames, g.Name)
		p.groupIDs = append(p.groupIDs, g.ID)
	}

	// set up mode-specific handler
	switch p.Mode {
	case ACLTemplating:
		p.templateHandler = aclTemplateHandler
	case JSONTemplating:
		p.templateHandler = jsonTemplateHandler
	default:
		return false, "", fmt.Errorf("unknown mode %q", p.Mode)
	}

	var subst bool
	splitStr := strings.Split(p.String, "{{")

	if len(splitStr) >= 1 {
		if strings.Contains(splitStr[0], "}}") {
			return false, "", ErrUnbalancedTemplatingCharacter
		}
		if len(splitStr) == 1 {
			return false, p.String, nil
		}
	}

	var b strings.Builder
	if !p.ValidityCheckOnly {
		b.Grow(2 * len(p.String))
	}

	for i, str := range splitStr {
		if i == 0 {
			if !p.ValidityCheckOnly {
				b.WriteString(str)
			}
			continue
		}
		splitPiece := strings.Split(str, "}}")
		switch len(splitPiece) {
		case 2:
			subst = true
			if !p.ValidityCheckOnly {
				tmplStr, err := performTemplating(strings.TrimSpace(splitPiece[0]), &p)
				if err != nil {
					return false, "", err
				}
				b.WriteString(tmplStr)
				b.WriteString(splitPiece[1])
			}
		default:
			return false, "", ErrUnbalancedTemplatingCharacter
		}
	}

	return subst, b.String(), nil
}

func performTemplating(input string, p *PopulateStringInput) (string, error) {

	performAliasTemplating := func(trimmed string, alias *logical.Alias) (string, error) {
		switch {
		case trimmed == "id":
			return p.templateHandler(alias.ID)

		case trimmed == "name":
			return p.templateHandler(alias.Name)

		case trimmed == "metadata":
			return p.templateHandler(alias.Metadata)

		case strings.HasPrefix(trimmed, "metadata."):
			split := strings.SplitN(trimmed, ".", 2)
			return p.templateHandler(alias.Metadata, split[1])
		}

		return "", ErrTemplateValueNotFound
	}

	performEntityTemplating := func(trimmed string) (string, error) {
		switch {
		case trimmed == "id":
			return p.templateHandler(p.Entity.ID)

		case trimmed == "name":
			return p.templateHandler(p.Entity.Name)

		case trimmed == "metadata":
			return p.templateHandler(p.Entity.Metadata)

		case strings.HasPrefix(trimmed, "metadata."):
			split := strings.SplitN(trimmed, ".", 2)
			return p.templateHandler(p.Entity.Metadata, split[1])

		case trimmed == "groups.names":
			return p.templateHandler(p.groupNames)

		case trimmed == "groups.ids":
			return p.templateHandler(p.groupIDs)

		case strings.HasPrefix(trimmed, "aliases."):
			split := strings.SplitN(strings.TrimPrefix(trimmed, "aliases."), ".", 2)
			if len(split) != 2 {
				return "", errors.New("invalid alias selector")
			}
			var alias *logical.Alias
			for _, a := range p.Entity.Aliases {
				if split[0] == a.MountAccessor {
					alias = a
					break
				}
			}
			if alias == nil {
				if p.Mode == ACLTemplating {
					return "", errors.New("alias not found")
				}

				// An empty alias is sufficient for generating defaults
				alias = &logical.Alias{Metadata: make(map[string]string)}
			}
			return performAliasTemplating(split[1], alias)
		}

		return "", ErrTemplateValueNotFound
	}

	performGroupsTemplating := func(trimmed string) (string, error) {
		var ids bool

		selectorSplit := strings.SplitN(trimmed, ".", 2)

		switch {
		case len(selectorSplit) != 2:
			return "", errors.New("invalid groups selector")

		case selectorSplit[0] == "ids":
			ids = true

		case selectorSplit[0] == "names":

		default:
			return "", errors.New("invalid groups selector")
		}
		trimmed = selectorSplit[1]

		accessorSplit := strings.SplitN(trimmed, ".", 2)
		if len(accessorSplit) != 2 {
			return "", errors.New("invalid groups accessor")
		}
		var found *logical.Group
		for _, group := range p.Groups {
			var compare string
			if ids {
				compare = group.ID
			} else {
				if p.NamespaceID != "" && group.NamespaceID != p.NamespaceID {
					continue
				}
				compare = group.Name
			}

			if compare == accessorSplit[0] {
				found = group
				break
			}
		}

		if found == nil {
			return "", fmt.Errorf("entity is not a member of group %q", accessorSplit[0])
		}

		trimmed = accessorSplit[1]

		switch {
		case trimmed == "id":
			return found.ID, nil

		case trimmed == "name":
			if found.Name == "" {
				return "", ErrTemplateValueNotFound
			}
			return found.Name, nil

		case strings.HasPrefix(trimmed, "metadata."):
			val, ok := found.Metadata[strings.TrimPrefix(trimmed, "metadata.")]
			if !ok {
				return "", ErrTemplateValueNotFound
			}
			return val, nil
		}

		return "", ErrTemplateValueNotFound
	}

	performTimeTemplating := func(trimmed string) (string, error) {
		now := p.Now
		if now.IsZero() {
			now = time.Now()
		}

		opsSplit := strings.SplitN(trimmed, ".", 3)

		if opsSplit[0] != "now" {
			return "", fmt.Errorf("invalid time selector %q", opsSplit[0])
		}

		result := now
		switch len(opsSplit) {
		case 1:
			// return current time
		case 2:
			return "", errors.New("missing time operand")

		case 3:
			duration, err := time.ParseDuration(opsSplit[2])
			if err != nil {
				return "", errwrap.Wrapf("invalid duration: {{err}}", err)
			}

			switch opsSplit[1] {
			case "plus":
				result = result.Add(duration)
			case "minus":
				result = result.Add(-duration)
			default:
				return "", fmt.Errorf("invalid time operator %q", opsSplit[1])
			}
		}

		return strconv.FormatInt(result.Unix(), 10), nil
	}

	switch {
	case strings.HasPrefix(input, "identity.entity."):
		if p.Entity == nil {
			return "", ErrNoEntityAttachedToToken
		}
		return performEntityTemplating(strings.TrimPrefix(input, "identity.entity."))

	case strings.HasPrefix(input, "identity.groups."):
		if len(p.Groups) == 0 {
			return "", ErrNoGroupsAttachedToToken
		}
		return performGroupsTemplating(strings.TrimPrefix(input, "identity.groups."))

	case strings.HasPrefix(input, "time."):
		return performTimeTemplating(strings.TrimPrefix(input, "time."))
	}

	return "", ErrTemplateValueNotFound
}
