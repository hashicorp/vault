// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestShouldIgnoreRotationError(t *testing.T) {
	tests := map[string]struct {
		err         error
		ignore      []string
		expectIgnore bool
	}{
		"nil error": {
			err:          nil,
			ignore:       []string{"UserNotFound"},
			expectIgnore: false,
		},
		"non-command error": {
			err:          mongo.ErrClientDisconnected,
			ignore:       []string{"UserNotFound"},
			expectIgnore: false,
		},
		"command error name matches": {
			err:          mongo.CommandError{Name: "UserNotFound"},
			ignore:       []string{"UserNotFound"},
			expectIgnore: true,
		},
		"command error name does not match": {
			err:          mongo.CommandError{Name: "UserNotFound"},
			ignore:       []string{"DuplicateKey"},
			expectIgnore: false,
		},
		"ignore list contains blanks": {
			err:          mongo.CommandError{Name: "UserAlreadyExists"},
			ignore:       []string{"", "UserAlreadyExists"},
			expectIgnore: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := shouldIgnoreRotationError(tc.err, tc.ignore)
			if got != tc.expectIgnore {
				t.Fatalf("got %v, expected %v", got, tc.expectIgnore)
			}
		})
	}
}

