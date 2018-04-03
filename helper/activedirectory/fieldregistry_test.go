package activedirectory

import (
	"testing"
)

func TestFieldRegistryListsFields(t *testing.T) {
	fields := FieldRegistry.List()
	if len(fields) != 36 {
		t.FailNow()
	}
}

func TestFieldRegistryEqualityComparisonsWork(t *testing.T) {

	fields := FieldRegistry.List()

	foundGivenName := false
	foundSurname := false

	for _, field := range fields {
		if field == FieldRegistry.GivenName {
			foundGivenName = true
		}
		if field == FieldRegistry.Surname {
			foundSurname = true
		}
	}

	if !foundGivenName || !foundSurname {
		t.FailNow()
	}
}

func TestFieldRegistryParsesFieldsByString(t *testing.T) {

	field, err := FieldRegistry.Parse("sn")
	if err != nil {
		t.FailNow()
	}

	if field != FieldRegistry.Surname {
		t.FailNow()
	}
}
