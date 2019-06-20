package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"

	"github.com/hashicorp/vault/helper/namespace"
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
	Mode              int
	ValidityCheckOnly bool
	String            string
	Entity            *Entity
	Groups            []*Group
	Namespace         *namespace.Namespace

	// Optional time to use during templating. In unset, current time will be used
	Now time.Time
}

func PopulateString(p *PopulateStringInput) (bool, string, error) {
	if p == nil {
		return false, "", errors.New("nil input")
	}

	if p.String == "" {
		return false, "", nil
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
				tmplStr, err := performTemplating(strings.TrimSpace(splitPiece[0]), p)
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

	// Handle quote wrapping uniformly. For ACL templating, this is a no-op.
	quote := func(s string) string { return s }
	if p.Mode == JSONTemplating {
		quote = strconv.Quote
	}

	performAliasTemplating := func(trimmed string, alias *Alias) (string, error) {
		switch {
		case trimmed == "id":
			return quote(alias.ID), nil
		case trimmed == "name":
			if alias.Name == "" {
				return "", ErrTemplateValueNotFound
			}
			return alias.Name, nil
		case trimmed == "metadata":
			if p.Mode == ACLTemplating {
				return "", ErrTemplateValueNotFound
			}

			jsonMetadata, err := marshalMetadata(alias.Metadata)
			if err != nil {
				return "", err
			}
			return jsonMetadata, nil

		case strings.HasPrefix(trimmed, "metadata."):
			split := strings.SplitN(trimmed, ".", 2)

			switch len(split) {
			case 2:
				val, ok := alias.Metadata[split[1]]
				if !ok && p.Mode == ACLTemplating {
					return "", ErrTemplateValueNotFound
				}
				return quote(val), nil
			}
		}

		return "", ErrTemplateValueNotFound
	}

	performEntityTemplating := func(trimmed string) (string, error) {
		switch {
		case trimmed == "id":
			return quote(p.Entity.ID), nil
		case trimmed == "name":
			if p.Entity.Name == "" && p.Mode == ACLTemplating {
				return "", ErrTemplateValueNotFound
			}
			return quote(p.Entity.Name), nil
		case trimmed == "metadata":
			if p.Mode == ACLTemplating {
				return "", ErrTemplateValueNotFound
			}

			jsonMetadata, err := marshalMetadata(p.Entity.Metadata)
			if err != nil {
				return "", err
			}
			return jsonMetadata, nil
		case strings.HasPrefix(trimmed, "metadata."):
			split := strings.SplitN(trimmed, ".", 2)

			switch len(split) {
			case 2:
				val, ok := p.Entity.Metadata[split[1]]
				if !ok && p.Mode == ACLTemplating {
					return "", ErrTemplateValueNotFound
				}
				return quote(val), nil
			}
		case trimmed == "group_names":
			if p.Mode == ACLTemplating {
				return "", ErrTemplateValueNotFound
			}
			return listGroups(p.Groups, "name"), nil

		case trimmed == "group_ids":
			if p.Mode == ACLTemplating {
				return "", ErrTemplateValueNotFound
			}
			return listGroups(p.Groups, "id"), nil

		case strings.HasPrefix(trimmed, "aliases."):
			split := strings.SplitN(strings.TrimPrefix(trimmed, "aliases."), ".", 2)
			if len(split) != 2 {
				return "", errors.New("invalid alias selector")
			}
			var alias *Alias
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
				alias = &Alias{Metadata: make(map[string]string)}
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
		var found *Group
		for _, group := range p.Groups {
			var compare string
			if ids {
				compare = group.ID
			} else {
				if p.Namespace != nil && group.NamespaceID == p.Namespace.ID {
					compare = group.Name
				} else {
					continue
				}
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

func listGroups(groups []*Group, element string) string {
	var out strings.Builder

	out.WriteString("[")
	for i, g := range groups {
		var v string
		switch element {
		case "name":
			v = g.Name
		case "id":
			v = g.ID
		}
		if i > 0 {
			out.WriteString(",")
		}
		out.WriteString(strconv.Quote(v))
	}
	out.WriteString("]")

	return out.String()
}

// marshalMetadata converts a metadata object into JSON, with special handling
// for nil objects which are rendered as {} instead of null.
func marshalMetadata(metadata map[string]string) (string, error) {
	if metadata == nil {
		return "{}", nil
	}

	enc, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}

	return string(enc), nil
}
