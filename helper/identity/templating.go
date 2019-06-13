package identity

import (
	"bytes"
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

type PopulateStringInput struct {
	AllowMissingSelectors bool
	ValidityCheckOnly     bool
	String                string
	Entity                *Entity
	Groups                []*Group
	Namespace             *namespace.Namespace
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
				tmplStr, err := performTemplating(p.Namespace, strings.TrimSpace(splitPiece[0]), p.Entity, p.Groups, p.AllowMissingSelectors)
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

func performTemplating(ns *namespace.Namespace, input string, entity *Entity, groups []*Group, allowMissingSelectors bool) (string, error) {
	performAliasTemplating := func(trimmed string, alias *Alias) (string, error) {
		switch {
		case trimmed == "id":
			return alias.ID, nil
		case trimmed == "name":
			if alias.Name == "" {
				return "", ErrTemplateValueNotFound
			}
			return alias.Name, nil
		case strings.HasPrefix(trimmed, "metadata.") || trimmed == "metadata":
			split := strings.SplitN(trimmed, ".", 2)

			switch len(split) {
			case 1:
				// all metadata as json, without outer braces
				jsonMetadata, err := json.Marshal(alias.Metadata)
				if err != nil {
					return "", err
				}
				jsonMetadata = bytes.Trim(jsonMetadata, "{}")
				return string(jsonMetadata), nil
			case 2:
				val, ok := alias.Metadata[split[1]]
				if !ok && !allowMissingSelectors {
					return "", ErrTemplateValueNotFound
				}
				return val, nil
			}
		}

		return "", ErrTemplateValueNotFound
	}

	performEntityTemplating := func(trimmed string) (string, error) {
		switch {
		case trimmed == "id":
			return entity.ID, nil
		case trimmed == "name":
			if entity.Name == "" && !allowMissingSelectors {
				return "", ErrTemplateValueNotFound
			}
			return entity.Name, nil
		case strings.HasPrefix(trimmed, "metadata.") || trimmed == "metadata":
			split := strings.SplitN(trimmed, ".", 2)

			switch len(split) {
			case 1:
				// all metadata as json, without outer braces
				jsonMetadata, err := json.Marshal(entity.Metadata)
				if err == nil {
					jsonMetadata = bytes.Trim(jsonMetadata, "{}")
					return string(jsonMetadata), nil
				}
				return "", nil
			case 2:
				val, ok := entity.Metadata[split[1]]
				if !ok && !allowMissingSelectors {
					return "", ErrTemplateValueNotFound
				}
				return val, nil
			}
		case trimmed == "group_names":
			return listGroups(groups, "name"), nil

		case trimmed == "group_ids":
			return listGroups(groups, "id"), nil

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
				if !allowMissingSelectors {
					return "", errors.New("alias not found")
				}
				return "", nil
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

	performTimeTemplating := func(trimmed string) (string, error) {
		opsSplit := strings.SplitN(trimmed, ".", 3)

		if opsSplit[0] != "now" {
			return "", fmt.Errorf("invalid time selector %q", opsSplit[0])
		}

		result := time.Now()
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
		if entity == nil {
			return "", ErrNoEntityAttachedToToken
		}
		return performEntityTemplating(strings.TrimPrefix(input, "identity.entity."))

	case strings.HasPrefix(input, "identity.groups."):
		if len(groups) == 0 {
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

	for _, g := range groups {
		var v string
		switch element {
		case "name":
			v = g.Name
		case "id":
			v = g.ID
		}
		out.WriteString(strconv.Quote(v))
		out.WriteString(",")
	}

	return strings.TrimSuffix(out.String(), ",")
}
