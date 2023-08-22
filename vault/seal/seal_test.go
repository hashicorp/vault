package seal

import (
	"testing"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
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
