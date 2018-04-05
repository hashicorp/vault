package iamutil

import (
	"fmt"

	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/util"
)

const (
	ServiceAccountMemberTmpl = "serviceAccount:%s"
)

type Policy struct {
	Bindings []*Binding `json:"bindings,omitempty"`
	Etag     string     `json:"etag,omitempty"`
}

type Binding struct {
	Members []string `json:"members,omitempty"`
	Role    string   `json:"role,omitempty"`
}

type PolicyDelta struct {
	Roles util.StringSet
	Email string
}

func (p *Policy) AddBindings(toAdd *PolicyDelta) (changed bool, updated *Policy) {
	return p.ChangedBindings(toAdd, nil)
}

func (p *Policy) RemoveBindings(toRemove *PolicyDelta) (changed bool, updated *Policy) {
	return p.ChangedBindings(nil, toRemove)
}

func (p *Policy) ChangedBindings(toAdd *PolicyDelta, toRemove *PolicyDelta) (changed bool, updated *Policy) {
	if toAdd == nil && toRemove == nil {
		return false, p
	}

	var toAddMem, toRemoveMem string
	if toAdd != nil {
		toAddMem = fmt.Sprintf(ServiceAccountMemberTmpl, toAdd.Email)
	}
	if toRemove != nil {
		toRemoveMem = fmt.Sprintf(ServiceAccountMemberTmpl, toRemove.Email)
	}

	changed = false

	newBindings := make([]*Binding, 0, len(p.Bindings))
	alreadyAdded := make(util.StringSet)

	for _, bind := range p.Bindings {
		memberSet := util.ToSet(bind.Members)

		if toAdd != nil {
			if toAdd.Roles.Includes(bind.Role) {
				changed = true
				alreadyAdded.Add(bind.Role)
				memberSet.Add(toAddMem)
			}
		}

		if toRemove != nil {
			if toRemove.Roles.Includes(bind.Role) {
				if memberSet.Includes(toRemoveMem) {
					changed = true
					delete(memberSet, toRemoveMem)
				}
			}
		}

		if len(memberSet) > 0 {
			newBindings = append(newBindings, &Binding{
				Role:    bind.Role,
				Members: memberSet.ToSlice(),
			})
		}
	}

	if toAdd != nil {
		for r := range toAdd.Roles {
			if !alreadyAdded.Includes(r) {
				changed = true
				newBindings = append(newBindings, &Binding{
					Role:    r,
					Members: []string{toAddMem},
				})
			}
		}
	}

	if changed {
		return true, &Policy{
			Bindings: newBindings,
			Etag:     p.Etag,
		}
	}
	return false, p
}
