package framework

import (
	"errors"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/logical"
)

// PopulateIdentityTemplate takes a template string, an entity ID, and an
// instance of system view. It will query system view for information about the
// entity and use the resulting identity information to populate the template
// string.
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

	input := identitytpl.PopulateStringInput{
		String: tpl,
		Entity: entity,
		Groups: groups,
		Mode:   identitytpl.ACLTemplating,
	}

	_, out, err := identitytpl.PopulateString(input)
	if err != nil {
		return "", err
	}

	return out, nil
}

// ValidateIdentityTemplate takes a template string and returns if the string is
// a valid identity template.
func ValidateIdentityTemplate(tpl string) (bool, error) {
	hasTemplating, _, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
		Mode:              identitytpl.ACLTemplating,
		ValidityCheckOnly: true,
		String:            tpl,
	})
	if err != nil {
		return false, errwrap.Wrapf("failed to validate policy templating: {{err}}", err)
	}

	return hasTemplating, nil
}
