// +build all unit

package gocql

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestUUIDNil(t *testing.T) {
	var uuid UUID
	want, got := "00000000-0000-0000-0000-000000000000", uuid.String()
	if want != got {
		t.Fatalf("TestNil: expected %q got %q", want, got)
	}
}

var testsUUID = []struct {
	input   string
	variant int
	version int
}{
	{"b4f00409-cef8-4822-802c-deb20704c365", VariantIETF, 4},
	{"B4F00409-CEF8-4822-802C-DEB20704C365", VariantIETF, 4}, //Use capital letters
	{"f81d4fae-7dec-11d0-a765-00a0c91e6bf6", VariantIETF, 1},
	{"00000000-7dec-11d0-a765-00a0c91e6bf6", VariantIETF, 1},
	{"3051a8d7-aea7-1801-e0bf-bc539dd60cf3", VariantFuture, 1},
	{"3051a8d7-aea7-2801-e0bf-bc539dd60cf3", VariantFuture, 2},
	{"3051a8d7-aea7-3801-e0bf-bc539dd60cf3", VariantFuture, 3},
	{"3051a8d7-aea7-4801-e0bf-bc539dd60cf3", VariantFuture, 4},
	{"3051a8d7-aea7-3801-e0bf-bc539dd60cf3", VariantFuture, 5},
	{"d0e817e1-e4b1-1801-3fe6-b4b60ccecf9d", VariantNCSCompat, 0},
	{"d0e817e1-e4b1-1801-bfe6-b4b60ccecf9d", VariantIETF, 1},
	{"d0e817e1-e4b1-1801-dfe6-b4b60ccecf9d", VariantMicrosoft, 0},
	{"d0e817e1-e4b1-1801-ffe6-b4b60ccecf9d", VariantFuture, 0},
}

func TestPredefinedUUID(t *testing.T) {
	for i := range testsUUID {
		uuid, err := ParseUUID(testsUUID[i].input)
		if err != nil {
			t.Errorf("ParseUUID #%d: %v", i, err)
			continue
		}

		if str := uuid.String(); str != strings.ToLower(testsUUID[i].input) {
			t.Errorf("String #%d: expected %q got %q", i, testsUUID[i].input, str)
			continue
		}

		if variant := uuid.Variant(); variant != testsUUID[i].variant {
			t.Errorf("Variant #%d: expected %d got %d", i, testsUUID[i].variant, variant)
		}

		if testsUUID[i].variant == VariantIETF {
			if version := uuid.Version(); version != testsUUID[i].version {
				t.Errorf("Version #%d: expected %d got %d", i, testsUUID[i].version, version)
			}
		}

		json, err := uuid.MarshalJSON()
		if err != nil {
			t.Errorf("MarshalJSON #%d: %v", i, err)
		}
		expectedJson := `"` + strings.ToLower(testsUUID[i].input) + `"`
		if string(json) != expectedJson {
			t.Errorf("MarshalJSON #%d: expected %v got %v", i, expectedJson, string(json))
		}

		var unmarshaled UUID
		err = unmarshaled.UnmarshalJSON(json)
		if err != nil {
			t.Errorf("UnmarshalJSON #%d: %v", i, err)
		}
		if unmarshaled != uuid {
			t.Errorf("UnmarshalJSON #%d: expected %v got %v", i, uuid, unmarshaled)
		}
	}
}

func TestInvalidUUIDCharacter(t *testing.T) {
	_, err := ParseUUID("z4f00409-cef8-4822-802c-deb20704c365")
	if err == nil || !strings.Contains(err.Error(), "invalid UUID") {
		t.Fatalf("expected invalid UUID error, got '%v' ", err)
	}
}

func TestInvalidUUIDLength(t *testing.T) {
	_, err := ParseUUID("4f00")
	if err == nil || !strings.Contains(err.Error(), "invalid UUID") {
		t.Fatalf("expected invalid UUID error, got '%v' ", err)
	}

	_, err = UUIDFromBytes(TimeUUID().Bytes()[:15])
	if err == nil || err.Error() != "UUIDs must be exactly 16 bytes long" {
		t.Fatalf("expected error '%v', got '%v'", "UUIDs must be exactly 16 bytes long", err)
	}
}

func TestRandomUUID(t *testing.T) {
	for i := 0; i < 20; i++ {
		uuid, err := RandomUUID()
		if err != nil {
			t.Errorf("RandomUUID: %v", err)
		}
		if variant := uuid.Variant(); variant != VariantIETF {
			t.Errorf("wrong variant. expected %d got %d", VariantIETF, variant)
		}
		if version := uuid.Version(); version != 4 {
			t.Errorf("wrong version. expected %d got %d", 4, version)
		}
	}
}

func TestRandomUUIDInvalidAPICalls(t *testing.T) {
	uuid, err := RandomUUID()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if node := uuid.Node(); node != nil {
		t.Fatalf("expected nil, got %v", node)
	}

	if stamp := uuid.Timestamp(); stamp != 0 {
		t.Fatalf("expceted 0, got %v", stamp)
	}
	zeroT := time.Time{}
	if to := uuid.Time(); to != zeroT {
		t.Fatalf("expected %v, got %v", zeroT, to)
	}
}

func TestUUIDFromTime(t *testing.T) {
	date := time.Date(1982, 5, 5, 12, 34, 56, 400, time.UTC)
	uuid := UUIDFromTime(date)

	if uuid.Time() != date {
		t.Errorf("embedded time incorrect. Expected %v got %v", date, uuid.Time())
	}
}

func TestParseUUID(t *testing.T) {
	uuid, _ := ParseUUID("486f3a88-775b-11e3-ae07-d231feb1dc81")
	if uuid.Time() != time.Date(2014, 1, 7, 5, 19, 29, 222516000, time.UTC) {
		t.Errorf("Expected date of 1/7/2014 at 5:19:29.222516, got %v", uuid.Time())
	}
}

func TestTimeUUID(t *testing.T) {
	var node []byte
	timestamp := int64(0)
	for i := 0; i < 20; i++ {
		uuid := TimeUUID()

		if variant := uuid.Variant(); variant != VariantIETF {
			t.Errorf("wrong variant. expected %d got %d", VariantIETF, variant)
		}
		if version := uuid.Version(); version != 1 {
			t.Errorf("wrong version. expected %d got %d", 1, version)
		}

		if n := uuid.Node(); !bytes.Equal(n, node) && i > 0 {
			t.Errorf("wrong node. expected %x, got %x", node, n)
		} else if i == 0 {
			node = n
		}

		ts := uuid.Timestamp()
		if ts < timestamp {
			t.Errorf("timestamps must grow")
		}
		timestamp = ts
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var withHyphens, withoutHypens, tooLong UUID

	withHyphens.UnmarshalJSON([]byte(`"486f3a88-775b-11e3-ae07-d231feb1dc81"`))
	if withHyphens.Time().Truncate(time.Second) != time.Date(2014, 1, 7, 5, 19, 29, 0, time.UTC) {
		t.Errorf("Expected date of 1/7/2014 at 5:19:29, got %v", withHyphens.Time())
	}

	withoutHypens.UnmarshalJSON([]byte(`"486f3a88775b11e3ae07d231feb1dc81"`))
	if withoutHypens.Time().Truncate(time.Second) != time.Date(2014, 1, 7, 5, 19, 29, 0, time.UTC) {
		t.Errorf("Expected date of 1/7/2014 at 5:19:29, got %v", withoutHypens.Time())
	}

	err := tooLong.UnmarshalJSON([]byte(`"486f3a88-775b-11e3-ae07-d231feb1dc81486f3a88"`))
	if err == nil {
		t.Errorf("no error for invalid JSON UUID")
	}

}
