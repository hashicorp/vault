package models

import (
	"errors"
)

type LifecycleEventType int

const (
	MISSED_LIFECYCLEEVENTTYPE LifecycleEventType = iota
	SUBSCRIPTIONREMOVED_LIFECYCLEEVENTTYPE
	REAUTHORIZATIONREQUIRED_LIFECYCLEEVENTTYPE
)

func (i LifecycleEventType) String() string {
	return []string{"missed", "subscriptionRemoved", "reauthorizationRequired"}[i]
}
func ParseLifecycleEventType(v string) (any, error) {
	result := MISSED_LIFECYCLEEVENTTYPE
	switch v {
	case "missed":
		result = MISSED_LIFECYCLEEVENTTYPE
	case "subscriptionRemoved":
		result = SUBSCRIPTIONREMOVED_LIFECYCLEEVENTTYPE
	case "reauthorizationRequired":
		result = REAUTHORIZATIONREQUIRED_LIFECYCLEEVENTTYPE
	default:
		return 0, errors.New("Unknown LifecycleEventType value: " + v)
	}
	return &result, nil
}
func SerializeLifecycleEventType(values []LifecycleEventType) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.String()
	}
	return result
}
