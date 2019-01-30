package identity

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/namespace"
)

var (
	ErrUnbalancedTemplatingCharacter = errors.New("unbalanced templating characters")
	ErrNoEntityAttachedToToken       = errors.New("string contains entity template directives but no entity was provided")
	ErrNoGroupsAttachedToToken       = errors.New("string contains groups template directives but no groups were provided")
	ErrTemplateValueNotFound         = errors.New("no value could be found for one of the template directives")
)

type PopulateStringInput struct {
	ValidityCheckOnly bool
	String            string
	Entity            *Entity
	Groups            []*Group
	Namespace         *namespace.Namespace
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
				tmplStr, err := performTemplating(p.Namespace, strings.TrimSpace(splitPiece[0]), p.Entity, p.Groups)
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

func performTemplating(ns *namespace.Namespace, input string, entity *Entity, groups []*Group) (string, error) {
	performAliasTemplating := func(trimmed string, alias *Alias) (string, error) {
		switch {
		case trimmed == "id":
			return alias.ID, nil
		case trimmed == "name":
			if alias.Name == "" {
				return "", ErrTemplateValueNotFound
			}
			return alias.Name, nil
		case strings.HasPrefix(trimmed, "metadata."):
			val, ok := alias.Metadata[strings.TrimPrefix(trimmed, "metadata.")]
			if !ok {
				return "", ErrTemplateValueNotFound
			}
			return val, nil
		}

		return "", ErrTemplateValueNotFound
	}

	performEntityTemplating := func(trimmed string) (string, error) {
		switch {
		case trimmed == "id":
			return entity.ID, nil
		case trimmed == "name":
			if entity.Name == "" {
				return "", ErrTemplateValueNotFound
			}
			return entity.Name, nil
		case strings.HasPrefix(trimmed, "metadata."):
			val, ok := entity.Metadata[strings.TrimPrefix(trimmed, "metadata.")]
			if !ok {
				return "", ErrTemplateValueNotFound
			}
			return val, nil
		case strings.HasPrefix(trimmed, "aliases."):
			split := strings.SplitN(strings.TrimPrefix(trimmed, "aliases."), ".", 2)
			if len(split) != 2 {
				return "", errors.New("invalid alias selector")
			}
			var found *Alias
			for _, alias := range entity.Aliases {
				if split[0] == alias.MountAccessor {
					found = alias
					break
				}
			}
			if found == nil {
				return "", errors.New("alias not found")
			}
			return performAliasTemplating(split[1], found)
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
		for _, group := range groups {
			var compare string
			if ids {
				compare = group.ID
			} else {
				if ns != nil && group.NamespaceID == ns.ID {
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

	switch {
	case strings.HasPrefix(input, "identity.entity."):
		if entity == nil {
			return "", ErrNoEntityAttachedToToken
		}
		return performEntityTemplating(strings.TrimPrefix(input, "identity.entity."))

	case strings.HasPrefix(input, "identity.groups."):
		if len(groups) == 0 {
			return "", ErrNoGroupsAttachedToToken
		}
		return performGroupsTemplating(strings.TrimPrefix(input, "identity.groups."))
	}

	return "", ErrTemplateValueNotFound
}
