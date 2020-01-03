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

	groups, err := sysView.GroupsForEntity(entityID)
	if err != nil {
		return "", err
	}

	// TODO: Namespace bound?
	input := identity.PopulateStringInput{
		String: tpl,
		Entity: entity,
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
