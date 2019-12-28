package framework

import (
	"errors"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/identity"
	"github.com/hashicorp/vault/sdk/logical"
)

func PopulateIdentityTemplate(tpl string, entityID string, sysView logical.SystemView) (string, error) {
	entity, err := sysView.EntityInfo(entityID)
	if err != nil {
		return "", err
	}
	if entity == nil {
		return "", errors.New("no entity found")
	}

	logicalGroups, err := sysView.GroupsForEntity(entityID)
	if err != nil {
		return "", err
	}

	groups := make([]*identity.Group, len(logicalGroups))
	for i, g := range logicalGroups {
		groups[i] = identity.FromLogicalGroup(g)
	}

	// TODO: Namespace bound?
	input := identity.PopulateStringInput{
		String: tpl,
		Entity: identity.FromLogicalEntity(entity),
		Groups: groups,
		Mode:   identity.ACLTemplating,
	}

	_, out, err := identity.PopulateString(input)
	if err != nil {
		return "", err
	}

	return out, nil
}

func ValidateIdentityTemplate(tpl string) (bool, error) {
	hasTemplating, _, err := identity.PopulateString(identity.PopulateStringInput{
		Mode:              identity.ACLTemplating,
		ValidityCheckOnly: true,
		String:            tpl,
	})
	if err != nil {
		return false, errwrap.Wrapf("failed to validate policy templating: {{err}}", err)
	}

	return hasTemplating, nil
}
