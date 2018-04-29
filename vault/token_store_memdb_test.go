package vault

import (
	"reflect"
	"testing"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/token"
)

func TestTokenStore_MemDBIndexes(t *testing.T) {
	var err error
	_, ts, _, _ := TestCoreWithTokenStore(t)

	tm1 := &token.TokenMapping{
		ID:          "testid",
		TokenID:     "testtokenid",
		Accessor:    "testaccessor",
		ParentID:    "testparentid",
		CubbyholeID: "testcubbyholeid",
	}
	err = ts.UpsertTokenMapping(tm1)
	if err != nil {
		t.Fatal(err)
	}

	tm2 := &token.TokenMapping{
		ID:       "testid2",
		TokenID:  "testtokenid2",
		Accessor: "testaccessor2",
		// Use the same parent for both the mappings
		ParentID:    "testparentid",
		CubbyholeID: "testcubbyholeid2",
	}
	err = ts.UpsertTokenMapping(tm2)
	if err != nil {
		t.Fatal(err)
	}

	tmFetched, err := ts.MemDBTokenMappingByTokenID("testtokenid")
	if err != nil {
		t.Fatal(err)
	}

	if tmFetched != tm1 {
		t.Fatalf("bad: same reference expected")
	}

	if !reflect.DeepEqual(tm1, tmFetched) {
		t.Fatalf("bad: token mapping; expected: %#v\n actual: %#v", tm1, tmFetched)
	}

	tmFetched, err = ts.MemDBTokenMappingByAccessor("testaccessor")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tm1, tmFetched) {
		t.Fatalf("bad: token mapping; expected: %#v\n actual: %#v", tm1, tmFetched)
	}

	tmsFetched, err := ts.MemDBTokenMappingsByParentID("testparentid")
	if err != nil {
		t.Fatal(err)
	}
	if len(tmsFetched) != 2 {
		t.Fatalf("bad: length of mappings; expected: 2, actual: %d", len(tmsFetched))
	}
	tm1Found := false
	tm2Found := false
	for _, tm := range tmsFetched {
		if tm.ID == tm1.ID {
			tm1Found = true
		}
		if tm.ID == tm2.ID {
			tm2Found = true
		}
	}
	if !tm1Found || !tm2Found {
		t.Fatalf("expected both token mappings to be returned")
	}
}

func TestTokenStore_MemDBDeleteTokenMappingByTokenID(t *testing.T) {
	var err error
	_, ts, _, _ := TestCoreWithTokenStore(t)

	tm1 := &token.TokenMapping{
		ID:          "testid",
		TokenID:     "testtokenid",
		Accessor:    "testaccessor",
		ParentID:    "testparentid",
		CubbyholeID: "testcubbyholeid",
	}
	err = ts.UpsertTokenMapping(tm1)
	if err != nil {
		t.Fatal(err)
	}

	err = ts.DeleteTokenMappingByTokenID("testtokenid")
	if err != nil {
		t.Fatal(err)
	}

	tmFetched, err := ts.MemDBTokenMappingByTokenID("testtokenid")
	if err != nil {
		t.Fatal(err)
	}

	if tmFetched != nil {
		t.Fatalf("expected a nil token mapping")
	}
}

func TestTokenStore_MemDBTokenMappings(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	for i := 0; i < 100; i++ {
		random, err := uuid.GenerateUUID()
		if err != nil {
			t.Fatal(err)
		}
		tm := &token.TokenMapping{
			ID:          random + "id",
			TokenID:     random + "tokenid",
			Accessor:    random + "accessor",
			ParentID:    random + "parentid",
			CubbyholeID: random + "cubbyholeid",
		}

		err = ts.UpsertTokenMapping(tm)
		if err != nil {
			t.Fatal(err)
		}
	}

	tokenMappings, err := ts.MemDBTokenMappings()
	if err != nil {
		t.Fatal(err)
	}

	// 100 from the loop above and 1 for the root token
	if len(tokenMappings) != 101 {
		t.Fatalf("bad: len(tokenMappings); expected %d, got %d", 101, len(tokenMappings))
	}
}

func TestTokenStore_MemDBTokenMappingsByParentID(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "tokenid", "", []string{"root"})
	testMakeToken(t, ts, "tokenid", "secondtokenid", "", []string{"foo"})
	testMakeToken(t, ts, "tokenid", "thirdtokenid", "", []string{"bar"})

	tokenMapping, err := ts.MemDBTokenMappingByTokenID(root)
	if err != nil {
		t.Fatal(err)
	}

	if tokenMapping.TokenID != root {
		t.Fatalf("bad: token ID in mapping; expected %q, got %q", root, tokenMapping.TokenID)
	}

	tokenMapping, err = ts.MemDBTokenMappingByTokenID("tokenid")
	if err != nil {
		t.Fatal(err)
	}

	if tokenMapping.TokenID != "tokenid" {
		t.Fatalf("bad: token ID in mapping; expected %q, got %q", "tokenid", tokenMapping.TokenID)
	}

	tokenMapping, err = ts.MemDBTokenMappingByTokenID("secondtokenid")
	if err != nil {
		t.Fatal(err)
	}

	if tokenMapping.TokenID != "secondtokenid" {
		t.Fatalf("bad: token ID in mapping; expected %q, got %q", "secondtokenid", tokenMapping.TokenID)
	}

	tokenMappings, err := ts.MemDBTokenMappingsByParentID(root)
	if err != nil {
		t.Fatal(err)
	}

	if len(tokenMappings) != 1 || tokenMappings[0].TokenID != "tokenid" {
		t.Fatalf("failed to fetch token mapping by parent id")
	}

	tokenMappings, err = ts.MemDBTokenMappingsByParentID("tokenid")
	if err != nil {
		t.Fatal(err)
	}

	tm1Found := false
	tm2Found := false

	for _, tm := range tokenMappings {
		switch tm.TokenID {
		case "secondtokenid":
			tm1Found = true
		case "thirdtokenid":
			tm2Found = true
		}
	}

	if !tm1Found || !tm2Found {
		t.Fatalf("failed to fetch token mappings by parent id")
	}
}
