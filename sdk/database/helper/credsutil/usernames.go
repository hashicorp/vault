package credsutil

import (
	"fmt"
	"strings"
	"time"
)

type CaseOp int

const (
	KeepCase CaseOp = iota
	Uppercase
	Lowercase
)

type usernameBuilder struct {
	displayName string
	roleName    string
	separator   string

	maxLen        int
	caseOperation CaseOp
}

func (ub usernameBuilder) makeUsername() (string, error) {
	userUUID, err := RandomAlphaNumeric(20, false)
	if err != nil {
		return "", err
	}

	now := fmt.Sprint(time.Now().Unix())

	parts := []string{
		"v",
		ub.displayName,
		ub.roleName,
		userUUID,
		now,
	}
	username := joinNonEmpty(ub.separator, parts...)
	if ub.maxLen > 0 {
		username = trunc(username, ub.maxLen)
	}
	switch ub.caseOperation {
	case Lowercase:
		username = strings.ToLower(username)
	case Uppercase:
		username = strings.ToUpper(username)
	}

	return username, nil
}

type UsernameOpt func(*usernameBuilder)

func DisplayName(dispName string, maxLength int) UsernameOpt {
	return func(b *usernameBuilder) {
		if maxLength > 0 {
			dispName = trunc(dispName, maxLength)
		}
		b.displayName = dispName
	}
}

func RoleName(roleName string, maxLength int) UsernameOpt {
	return func(b *usernameBuilder) {
		if maxLength > 0 {
			roleName = trunc(roleName, maxLength)
		}
		b.roleName = roleName
	}
}

func Separator(sep string) UsernameOpt {
	return func(b *usernameBuilder) {
		b.separator = sep
	}
}

func MaxLength(maxLen int) UsernameOpt {
	return func(b *usernameBuilder) {
		b.maxLen = maxLen
	}
}

func Case(c CaseOp) UsernameOpt {
	return func(b *usernameBuilder) {
		b.caseOperation = c
	}
}

func ToLower() UsernameOpt {
	return Case(Lowercase)
}

func ToUpper() UsernameOpt {
	return Case(Uppercase)
}

func GenerateUsername(opts ...UsernameOpt) (string, error) {
	b := usernameBuilder{
		separator:     "_",
		maxLen:        100,
		caseOperation: KeepCase,
	}

	for _, opt := range opts {
		opt(&b)
	}

	return b.makeUsername()
}

func trunc(str string, l int) string {
	if len(str) < l {
		return str
	}
	return str[:l]
}

func joinNonEmpty(sep string, vals ...string) string {
	if sep == "" {
		return strings.Join(vals, sep)
	}
	switch len(vals) {
	case 0:
		return ""
	case 1:
		return vals[0]
	}
	builder := &strings.Builder{}
	for _, val := range vals {
		if val == "" {
			continue
		}
		if builder.Len() > 0 {
			builder.WriteString(sep)
		}
		builder.WriteString(val)
	}
	return builder.String()
}
