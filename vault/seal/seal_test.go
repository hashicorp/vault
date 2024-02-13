// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal

import (
	"context"
	"fmt"
	"testing"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/stretchr/testify/require"
)

func Test_keyIdSet(t *testing.T) {
	type args struct {
		value *MultiWrapValue
	}
	tests := []struct {
		name      string
		idsToSet  []string
		idsToTest []string
		want      bool
	}{
		{
			name:      "a single ID",
			idsToSet:  []string{"Nexus 6"},
			idsToTest: []string{"Nexus 6"},
			want:      true,
		},
		{
			name:      "two sets of equal IDs",
			idsToSet:  []string{"A", "B"},
			idsToTest: []string{"A", "B"},
			want:      true,
		},
		{
			name:      "two sets of equal IDs, with duplicated ids set",
			idsToSet:  []string{"A", "B", "A"},
			idsToTest: []string{"A", "B"},
			want:      true,
		},
		{
			name:      "two sets of equal IDs, with duplicated ids tested",
			idsToSet:  []string{"A", "B"},
			idsToTest: []string{"A", "B", "A"},
			want:      true,
		},
		{
			name:      "two sets of equal IDs in different order",
			idsToSet:  []string{"A", "B"},
			idsToTest: []string{"B", "A"},
			want:      true,
		},
		{
			name:      "two sets of different IDs",
			idsToSet:  []string{"A", "B"},
			idsToTest: []string{"B", "C"},
			want:      false,
		},
	}
	for _, tt := range tests {
		useSetIds := func(s *keyIdSet) {
			s.setIds(tt.idsToSet)
		}
		useSet := func(s *keyIdSet) {
			mwv := &MultiWrapValue{Generation: 6}
			for _, id := range tt.idsToSet {
				mwv.Slots = append(mwv.Slots, &wrapping.BlobInfo{
					KeyInfo: &wrapping.KeyInfo{
						KeyId: id,
					},
				})
			}
			s.set(mwv)
		}

		runTest := func(name string, setter func(*keyIdSet)) {
			t.Run(name, func(t *testing.T) {
				s := &keyIdSet{}
				setter(s)

				mwv := &MultiWrapValue{Generation: 6}
				for _, id := range tt.idsToTest {
					mwv.Slots = append(mwv.Slots, &wrapping.BlobInfo{
						KeyInfo: &wrapping.KeyInfo{
							KeyId: id,
						},
					})
				}
				if got := s.equal(mwv); got != tt.want {
					t.Errorf("equal() = %v, want %v, IDs set: %v, IDs tested: %v",
						got, tt.want, tt.idsToSet, tt.idsToTest)
				}
			})
		}
		runTest(tt.name+".set()", useSet)
		runTest(tt.name+".setIDs", useSetIds)
	}
}

// Test_Encrypt_duplicate_keyIds verifies that if two seal wrappers produce the same Key ID, an error
// will be returned for both.
func Test_Encrypt_duplicate_keyIds(t *testing.T) {
	ctx := context.Background()

	setId := func(w *SealWrapper, keyId string) {
		testWrapper := w.Wrapper.(*ToggleableWrapper).Wrapper.(*wrapping.TestWrapper)
		testWrapper.SetKeyId(keyId)
	}

	getId := func(w *SealWrapper) string {
		id, err := w.Wrapper.KeyId(ctx)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}

	access, _ := NewTestSeal(&TestSealOpts{WrapperCount: 3})

	// Set up - make the key IDs the same for the last two wrappers
	wrappers := access.GetAllSealWrappersByPriority()
	setId(wrappers[1], "this-key-is-duplicated")
	setId(wrappers[2], "this-key-is-duplicated")

	// Some sanity checks
	require.NotEqual(t, wrappers[0].Name, wrappers[1].Name)
	require.NotEqual(t, wrappers[1].Name, wrappers[2].Name)
	require.NotEqual(t, getId(wrappers[0]), getId(wrappers[1]))
	require.Equal(t, getId(wrappers[1]), getId(wrappers[2]))

	// Encrypt a value
	mwv, errorMap := access.Encrypt(ctx, []byte("Rinconete y Cortadillo"))

	// Assertions
	require.NotNilf(t, mwv, "seal 0 should have succeeded")

	requireDuplicateErr := func(w *SealWrapper) {
		require.ErrorContains(t, errorMap[w.Name], fmt.Sprintf("seal %v has returned duplicate key ID", w.Name))
	}
	requireDuplicateErr(wrappers[1])
	requireDuplicateErr(wrappers[2])
}
