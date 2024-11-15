// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

const bindingTemplate = "util/bindings_template"

func BindingsHCL(bindings map[string]StringSet) (string, error) {
	tpl, err := template.ParseFiles(bindingTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "bindings", bindings); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ParseBindings(bindingsStr string) (map[string]StringSet, error) {
	// Try to base64 decode
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(bindingsStr))
	decoded, b64err := ioutil.ReadAll(decoder)

	var bindsString string
	if b64err != nil {
		bindsString = bindingsStr
	} else {
		bindsString = string(decoded)
	}

	root, err := hcl.Parse(bindsString)
	if err != nil {
		if b64err == nil {
			return nil, errwrap.Wrapf("unable to parse base64-encoded bindings as valid HCL: {{err}}", err)
		} else {
			return nil, errwrap.Wrapf("unable to parse raw string bindings as valid HCL: {{err}}", err)
		}
	}

	bindingLst, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, errors.New("unable to parse bindings: does not contain a root object")
	}

	bindingsMap, err := parseBindingObjList(bindingLst)
	if err != nil {
		return nil, errwrap.Wrapf("unable to parse bindings: {{err}}", err)
	}
	return bindingsMap, nil
}

func parseBindingObjList(topList *ast.ObjectList) (map[string]StringSet, error) {
	var merr *multierror.Error

	bindings := make(map[string]StringSet)

	for _, item := range topList.Items {
		err := parseResourceObject(item, bindings)
		if err != nil {
			merr = multierror.Append(merr, fmt.Errorf("(line %d) %v", item.Assign.Line, err))
		}
	}
	err := merr.ErrorOrNil()
	if err != nil {
		return nil, err
	}
	return bindings, nil
}

func parseResourceObject(item *ast.ObjectItem, bindings map[string]StringSet) error {
	if len(item.Keys) != 2 || item.Keys[0] == nil || item.Keys[1] == nil {
		return fmt.Errorf(`top-level items must have format "resource" "$resource_name"`)
	}

	k, err := parseStringFromObjectKey(item, item.Keys[0])
	if err != nil {
		return err
	}
	if k != "resource" {
		return fmt.Errorf(`invalid item %q, expected "resource"`, k)
	}

	resourceName, err := parseStringFromObjectKey(item, item.Keys[1])
	if err != nil {
		return err
	}

	_, ok := bindings[resourceName]
	if !ok {
		bindings[resourceName] = make(StringSet)
	}
	boundRoles := bindings[resourceName]

	resourceItemList := item.Val.(*ast.ObjectType).List
	if resourceItemList == nil {
		return fmt.Errorf("invalid empty roles list for item (line %d)", item.Assign.Line)
	}

	var merr *multierror.Error
	for _, rolesObj := range resourceItemList.Items {
		err := parseRolesObject(rolesObj, boundRoles)
		if err != nil {
			merr = multierror.Append(merr, fmt.Errorf("role list (line %d): %v", rolesObj.Assign.Line, err))
		}
	}
	return merr.ErrorOrNil()
}

func parseRolesObject(rolesObj *ast.ObjectItem, parsedRoles StringSet) error {
	if rolesObj == nil || len(rolesObj.Keys) != 1 || rolesObj.Keys[0] == nil {
		return fmt.Errorf(`expected "roles" list, got nil object item`)
	}
	k, err := parseStringFromObjectKey(rolesObj, rolesObj.Keys[0])
	if err != nil {
		return err
	}
	if k != "roles" {
		return fmt.Errorf(`invalid key %q in resource, expected "roles"`, k)
	}

	if rolesObj.Val == nil {
		return fmt.Errorf(`expected "roles" list, got nil value`)
	}
	roleList, ok := rolesObj.Val.(*ast.ListType)
	if !ok {
		return fmt.Errorf("parsing error, expected list of roles for key 'roles'")
	}
	var merr *multierror.Error
	for _, singleRoleObj := range roleList.List {
		role, err := parseRole(rolesObj, singleRoleObj)
		if err != nil {
			merr = multierror.Append(merr, err)
		} else {
			parsedRoles.Add(role)
		}
	}
	return merr.ErrorOrNil()
}

func parseRole(parent *ast.ObjectItem, roleNode ast.Node) (string, error) {
	if roleNode == nil {
		return "", fmt.Errorf(`unexpected empty role item (line %d)`, parent.Assign.Line)
	}

	roleLitType, ok := roleNode.(*ast.LiteralType)
	if !ok || roleLitType == nil {
		return "", fmt.Errorf(`unexpected nil item in roles list (line %d)`, parent.Assign.Line)
	}

	roleRaw := roleLitType.Token.Value()
	role, ok := roleRaw.(string)
	if !ok {
		return "", fmt.Errorf(`unexpected item %v in roles list is not a string (line %d)`, roleRaw, parent.Assign.Line)
	}

	tkns := strings.Split(role, "/")
	if len(tkns) == 2 && tkns[0] == "roles" {
		return role, nil
	}
	if len(tkns) == 4 && tkns[2] == "roles" {
		// "projects/X/roles/Y" or "organizations/X/roles/Y"
		if tkns[0] == "projects" || tkns[0] == "organizations" {
			return role, nil
		}
	}
	return "", fmt.Errorf(`invalid role %q (line %d) must be one of following formats: "projects/X/roles/Y", "organizations/X/roles/Y", "roles/X"`, role, parent.Assign.Line)
}

func parseStringFromObjectKey(parent *ast.ObjectItem, k *ast.ObjectKey) (string, error) {
	if k == nil || k.Token.Value() == nil {
		return "", fmt.Errorf("expected string, got nil value (Llne %d)", parent.Assign.Line)
	}
	vRaw := k.Token.Value()
	v, ok := vRaw.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %v (Llne %d)", parent.Assign.Line, vRaw)
	}

	return v, nil
}
