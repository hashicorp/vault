package ldap_test

import (
	"strings"
	"testing"

	"gopkg.in/asn1-ber.v1"
	"gopkg.in/ldap.v2"
)

type compileTest struct {
	filterStr string

	expectedFilter string
	expectedType   int
	expectedErr    string
}

var testFilters = []compileTest{
	compileTest{
		filterStr:      "(&(sn=Miller)(givenName=Bob))",
		expectedFilter: "(&(sn=Miller)(givenName=Bob))",
		expectedType:   ldap.FilterAnd,
	},
	compileTest{
		filterStr:      "(|(sn=Miller)(givenName=Bob))",
		expectedFilter: "(|(sn=Miller)(givenName=Bob))",
		expectedType:   ldap.FilterOr,
	},
	compileTest{
		filterStr:      "(!(sn=Miller))",
		expectedFilter: "(!(sn=Miller))",
		expectedType:   ldap.FilterNot,
	},
	compileTest{
		filterStr:      "(sn=Miller)",
		expectedFilter: "(sn=Miller)",
		expectedType:   ldap.FilterEqualityMatch,
	},
	compileTest{
		filterStr:      "(sn=Mill*)",
		expectedFilter: "(sn=Mill*)",
		expectedType:   ldap.FilterSubstrings,
	},
	compileTest{
		filterStr:      "(sn=*Mill)",
		expectedFilter: "(sn=*Mill)",
		expectedType:   ldap.FilterSubstrings,
	},
	compileTest{
		filterStr:      "(sn=*Mill*)",
		expectedFilter: "(sn=*Mill*)",
		expectedType:   ldap.FilterSubstrings,
	},
	compileTest{
		filterStr:      "(sn=*i*le*)",
		expectedFilter: "(sn=*i*le*)",
		expectedType:   ldap.FilterSubstrings,
	},
	compileTest{
		filterStr:      "(sn=Mi*l*r)",
		expectedFilter: "(sn=Mi*l*r)",
		expectedType:   ldap.FilterSubstrings,
	},
	compileTest{
		filterStr:      "(sn=Mi*le*)",
		expectedFilter: "(sn=Mi*le*)",
		expectedType:   ldap.FilterSubstrings,
	},
	compileTest{
		filterStr:      "(sn=*i*ler)",
		expectedFilter: "(sn=*i*ler)",
		expectedType:   ldap.FilterSubstrings,
	},
	compileTest{
		filterStr:      "(sn>=Miller)",
		expectedFilter: "(sn>=Miller)",
		expectedType:   ldap.FilterGreaterOrEqual,
	},
	compileTest{
		filterStr:      "(sn<=Miller)",
		expectedFilter: "(sn<=Miller)",
		expectedType:   ldap.FilterLessOrEqual,
	},
	compileTest{
		filterStr:      "(sn=*)",
		expectedFilter: "(sn=*)",
		expectedType:   ldap.FilterPresent,
	},
	compileTest{
		filterStr:      "(sn~=Miller)",
		expectedFilter: "(sn~=Miller)",
		expectedType:   ldap.FilterApproxMatch,
	},
	compileTest{
		filterStr:      `(objectGUID='\fc\fe\a3\ab\f9\90N\aaGm\d5I~\d12)`,
		expectedFilter: `(objectGUID='\fc\fe\a3\ab\f9\90N\aaGm\d5I~\d12)`,
		expectedType:   ldap.FilterEqualityMatch,
	},
	compileTest{
		filterStr:      `(objectGUID=абвгдеёжзийклмнопрстуфхцчшщъыьэюя)`,
		expectedFilter: `(objectGUID=\c3\90\c2\b0\c3\90\c2\b1\c3\90\c2\b2\c3\90\c2\b3\c3\90\c2\b4\c3\90\c2\b5\c3\91\c2\91\c3\90\c2\b6\c3\90\c2\b7\c3\90\c2\b8\c3\90\c2\b9\c3\90\c2\ba\c3\90\c2\bb\c3\90\c2\bc\c3\90\c2\bd\c3\90\c2\be\c3\90\c2\bf\c3\91\c2\80\c3\91\c2\81\c3\91\c2\82\c3\91\c2\83\c3\91\c2\84\c3\91\c2\85\c3\91\c2\86\c3\91\c2\87\c3\91\c2\88\c3\91\c2\89\c3\91\c2\8a\c3\91\c2\8b\c3\91\c2\8c\c3\91\c2\8d\c3\91\c2\8e\c3\91\c2\8f)`,
		expectedType:   ldap.FilterEqualityMatch,
	},
	compileTest{
		filterStr:      `(objectGUID=함수목록)`,
		expectedFilter: `(objectGUID=\c3\ad\c2\95\c2\a8\c3\ac\c2\88\c2\98\c3\ab\c2\aa\c2\a9\c3\ab\c2\a1\c2\9d)`,
		expectedType:   ldap.FilterEqualityMatch,
	},
	compileTest{
		filterStr:      `(objectGUID=`,
		expectedFilter: ``,
		expectedType:   0,
		expectedErr:    "unexpected end of filter",
	},
	compileTest{
		filterStr:      `(objectGUID=함수목록`,
		expectedFilter: ``,
		expectedType:   0,
		expectedErr:    "unexpected end of filter",
	},
	// compileTest{ filterStr: "()", filterType: FilterExtensibleMatch },
}

var testInvalidFilters = []string{
	`(objectGUID=\zz)`,
	`(objectGUID=\a)`,
}

func TestFilter(t *testing.T) {
	// Test Compiler and Decompiler
	for _, i := range testFilters {
		filter, err := ldap.CompileFilter(i.filterStr)
		if err != nil {
			if i.expectedErr == "" || !strings.Contains(err.Error(), i.expectedErr) {
				t.Errorf("Problem compiling '%s' - '%v' (expected error to contain '%v')", i.filterStr, err, i.expectedErr)
			}
		} else if filter.Tag != ber.Tag(i.expectedType) {
			t.Errorf("%q Expected %q got %q", i.filterStr, ldap.FilterMap[uint64(i.expectedType)], ldap.FilterMap[uint64(filter.Tag)])
		} else {
			o, err := ldap.DecompileFilter(filter)
			if err != nil {
				t.Errorf("Problem compiling %s - %s", i.filterStr, err.Error())
			} else if i.expectedFilter != o {
				t.Errorf("%q expected, got %q", i.expectedFilter, o)
			}
		}
	}
}

func TestInvalidFilter(t *testing.T) {
	for _, filterStr := range testInvalidFilters {
		if _, err := ldap.CompileFilter(filterStr); err == nil {
			t.Errorf("Problem compiling %s - expected err", filterStr)
		}
	}
}

func BenchmarkFilterCompile(b *testing.B) {
	b.StopTimer()
	filters := make([]string, len(testFilters))

	// Test Compiler and Decompiler
	for idx, i := range testFilters {
		filters[idx] = i.filterStr
	}

	maxIdx := len(filters)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ldap.CompileFilter(filters[i%maxIdx])
	}
}

func BenchmarkFilterDecompile(b *testing.B) {
	b.StopTimer()
	filters := make([]*ber.Packet, len(testFilters))

	// Test Compiler and Decompiler
	for idx, i := range testFilters {
		filters[idx], _ = ldap.CompileFilter(i.filterStr)
	}

	maxIdx := len(filters)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ldap.DecompileFilter(filters[i%maxIdx])
	}
}
