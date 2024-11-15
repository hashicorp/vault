package models

import (
	"errors"
)

type ChangeType int

const (
	CREATED_CHANGETYPE ChangeType = iota
	UPDATED_CHANGETYPE
	DELETED_CHANGETYPE
)

func (i ChangeType) String() string {
	return []string{"created", "updated", "deleted"}[i]
}
func ParseChangeType(v string) (any, error) {
	result := CREATED_CHANGETYPE
	switch v {
	case "created":
		result = CREATED_CHANGETYPE
	case "updated":
		result = UPDATED_CHANGETYPE
	case "deleted":
		result = DELETED_CHANGETYPE
	default:
		return 0, errors.New("Unknown ChangeType value: " + v)
	}
	return &result, nil
}
func SerializeChangeType(values []ChangeType) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.String()
	}
	return result
}
