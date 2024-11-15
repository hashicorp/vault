// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/base62"
)

var (
	// Per https://en.wikipedia.org/wiki/Password_strength#Guidelines_for_strong_passwords
	minimumLengthOfComplexString = 8

	passwordComplexityPrefix = "?@09AZ"
	pwdFieldTmpl             = "{{PASSWORD}}"
)

type passwordGenerator interface {
	GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error)
}

// GeneratePassword from the password configuration. This will either generate based on a password policy
// or from the provided formatter. The formatter/length options are deprecated.
func GeneratePassword(ctx context.Context, passConf passwordConf, generator passwordGenerator) (password string, err error) {
	err = passConf.validate()
	if err != nil {
		return "", err
	}

	if passConf.PasswordPolicy != "" {
		return generator.GeneratePasswordFromPolicy(ctx, passConf.PasswordPolicy)
	}
	return generateDeprecatedPassword(passConf.Formatter, passConf.Length)
}

func generateDeprecatedPassword(formatter string, totalLength int) (string, error) {
	// Has formatter
	if formatter != "" {
		passLen := lengthOfPassword(formatter, totalLength)
		pwd, err := base62.Random(passLen)
		if err != nil {
			return "", err
		}
		return strings.Replace(formatter, pwdFieldTmpl, pwd, 1), nil
	}

	// Doesn't have formatter
	pwd, err := base62.Random(totalLength - len(passwordComplexityPrefix))
	if err != nil {
		return "", err
	}
	return passwordComplexityPrefix + pwd, nil
}

func lengthOfPassword(formatter string, totalLength int) int {
	lengthOfText := len(formatter) - len(pwdFieldTmpl)
	return totalLength - lengthOfText
}
