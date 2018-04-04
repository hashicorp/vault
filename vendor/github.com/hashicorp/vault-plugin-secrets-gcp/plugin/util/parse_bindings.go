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
		if len(item.Keys) != 2 {
			merr = multierror.Append(merr, fmt.Errorf("invalid resource item does not have ID on line %d", item.Assign.Line))
			continue
		}

		key := item.Keys[0].Token.Value().(string)
		if key != "resource" {
			merr = multierror.Append(merr, fmt.Errorf("invalid key '%s' (line %d)", key, item.Assign.Line))
			continue
		}

		resourceName := item.Keys[1].Token.Value().(string)
		_, ok := bindings[resourceName]
		if !ok {
			bindings[resourceName] = make(StringSet)
		}

		resourceList := item.Val.(*ast.ObjectType).List
		for _, rolesItem := range resourceList.Items {
			key := rolesItem.Keys[0].Token.Text
			switch key {
			case "roles":
				parseRoles(rolesItem, bindings[resourceName], merr)
			default:
				merr = multierror.Append(merr, fmt.Errorf("invalid key '%s' in resource '%s' (line %d)", key, resourceName, item.Assign.Line))
				continue
			}
		}
	}
	err := merr.ErrorOrNil()
	if err != nil {
		return nil, err
	}
	return bindings, nil
}

func parseRoles(item *ast.ObjectItem, roleSet StringSet, merr *multierror.Error) {
	lst, ok := item.Val.(*ast.ListType)
	if !ok {
		merr = multierror.Append(merr, fmt.Errorf("roles must be a list (line %d)", item.Assign.Line))
		return
	}

	for _, roleItem := range lst.List {
		role := roleItem.(*ast.LiteralType).Token.Value().(string)

		tkns := strings.Split(role, "/")
		switch len(tkns) {
		case 2:
			// "roles/X"
			if tkns[0] == "roles" {
				roleSet.Add(role)
				continue
			}
		case 4:
			// "projects/X/roles/Y" or "organizations/X/roles/Y"
			if (tkns[0] == "projects" || tkns[0] == "organizations") && tkns[2] == "roles" {
				roleSet.Add(role)
				continue
			}
		}

		merr = multierror.Append(merr, fmt.Errorf("invalid role: %s (line %d): must be project-level, organization-level, or global role", role, roleItem.Pos().Line))
	}
}
