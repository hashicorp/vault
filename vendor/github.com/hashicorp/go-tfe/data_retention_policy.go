// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import "regexp"

// DataRetentionPolicyChoice is a choice type struct that represents the possible types
// of a drp returned by a polymorphic relationship. If a value is available, exactly one field
// will be non-nil.
type DataRetentionPolicyChoice struct {
	DataRetentionPolicy            *DataRetentionPolicy
	DataRetentionPolicyDeleteOlder *DataRetentionPolicyDeleteOlder
	DataRetentionPolicyDontDelete  *DataRetentionPolicyDontDelete
}

// Returns whether one of the choices is populated
func (d DataRetentionPolicyChoice) IsPopulated() bool {
	return d.DataRetentionPolicy != nil ||
		d.DataRetentionPolicyDeleteOlder != nil ||
		d.DataRetentionPolicyDontDelete != nil
}

// Convert the DataRetentionPolicyChoice to the legacy DataRetentionPolicy struct
// Returns nil if the policy cannot be represented by a legacy DataRetentionPolicy
func (d *DataRetentionPolicyChoice) ConvertToLegacyStruct() *DataRetentionPolicy {
	if d == nil {
		return nil
	}
	if d.DataRetentionPolicy != nil {
		// TFE v202311-1 and v202312-1 will return a deprecated DataRetentionPolicy in the DataRetentionPolicyChoice struct
		return d.DataRetentionPolicy
	} else if d.DataRetentionPolicyDeleteOlder != nil {
		// DataRetentionPolicy was functionally replaced by DataRetentionPolicyDeleteOlder in TFE v202401
		return &DataRetentionPolicy{
			ID:                   d.DataRetentionPolicyDeleteOlder.ID,
			DeleteOlderThanNDays: d.DataRetentionPolicyDeleteOlder.DeleteOlderThanNDays,
		}
	}
	return nil
}

// DataRetentionPolicy describes the retention policy of deleting records older than the specified number of days.
//
// Deprecated: Use DataRetentionPolicyDeleteOlder instead. This is the original representation of a
// data retention policy, only present in TFE v202311-1 and v202312-1
type DataRetentionPolicy struct {
	ID                   string `jsonapi:"primary,data-retention-policies"`
	DeleteOlderThanNDays int    `jsonapi:"attr,delete-older-than-n-days"`
}

// DataRetentionPolicySetOptions is the options for a creating a DataRetentionPolicy.
//
// Deprecated: Use DataRetentionPolicyDeleteOlder variations instead
type DataRetentionPolicySetOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,data-retention-policies"`

	// DeleteOlderThanNDays is the number of days to retain records for.
	DeleteOlderThanNDays int `jsonapi:"attr,delete-older-than-n-days"`
}

// DataRetentionPolicyDeleteOlder describes the retention policy of deleting records older than the specified number of days.
type DataRetentionPolicyDeleteOlder struct {
	ID string `jsonapi:"primary,data-retention-policy-delete-olders"`

	// DeleteOlderThanNDays is the number of days to retain records for.
	DeleteOlderThanNDays int `jsonapi:"attr,delete-older-than-n-days"`
}

// DataRetentionPolicyDontDelete describes the retention policy of never deleting records.
type DataRetentionPolicyDontDelete struct {
	ID string `jsonapi:"primary,data-retention-policy-dont-deletes"`
}

// DataRetentionPolicyDeleteOlderSetOptions describes the options for a creating a DataRetentionPolicyDeleteOlder.
type DataRetentionPolicyDeleteOlderSetOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,data-retention-policy-delete-olders"`

	// DeleteOlderThanNDays is the number of days records will be retained for after their creation.
	DeleteOlderThanNDays int `jsonapi:"attr,delete-older-than-n-days"`
}

// DataRetentionPolicyDontDeleteSetOptions describes the options for a creating a DataRetentionPolicyDontDelete.
type DataRetentionPolicyDontDeleteSetOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,data-retention-policy-dont-deletes"`
}

// error we get when trying to unmarshal a data retention policy from TFE v202401+ into the deprecated DataRetentionPolicy struct
var drpUnmarshalEr = regexp.MustCompile(`Trying to Unmarshal an object of type \".+\", but \"data-retention-policies\" does not match`)
