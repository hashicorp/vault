package ldap

import (
	"reflect"
	"testing"

	"github.com/vanackere/asn1-ber"
)

type compileTest struct {
	filterStr  string
	filterType ber.Tag
}

var testFilters = []compileTest{
	compileTest{filterStr: "(&(sn=Miller)(givenName=Bob))", filterType: FilterAnd},
	compileTest{filterStr: "(|(sn=Miller)(givenName=Bob))", filterType: FilterOr},
	compileTest{filterStr: "(!(sn=Miller))", filterType: FilterNot},
	compileTest{filterStr: "(sn=Miller)", filterType: FilterEqualityMatch},
	compileTest{filterStr: "(sn=Mill*)", filterType: FilterSubstrings},
	compileTest{filterStr: "(sn=*Mill)", filterType: FilterSubstrings},
	compileTest{filterStr: "(sn=*Mill*)", filterType: FilterSubstrings},
	compileTest{filterStr: "(sn>=Miller)", filterType: FilterGreaterOrEqual},
	compileTest{filterStr: "(sn<=Miller)", filterType: FilterLessOrEqual},
	compileTest{filterStr: "(sn=*)", filterType: FilterPresent},
	compileTest{filterStr: "(sn~=Miller)", filterType: FilterApproxMatch},
	// compileTest{ filterStr: "()", filterType: FilterExtensibleMatch },
}

func TestFilter(t *testing.T) {
	// Test Compiler and Decompiler
	for _, i := range testFilters {
		filter, err := CompileFilter(i.filterStr)
		if err != nil {
			t.Errorf("Problem compiling %s - %s", i.filterStr, err.Error())
		} else if filter.Tag != i.filterType {
			t.Errorf("%q Expected %q got %q", i.filterStr, filterMap[i.filterType], filterMap[filter.Tag])
		} else {
			o, err := DecompileFilter(filter)
			if err != nil {
				t.Errorf("Problem compiling %s - %s", i.filterStr, err.Error())
			} else if i.filterStr != o {
				t.Errorf("%q expected, got %q", i.filterStr, o)
			}
		}
	}
}

type binTestFilter struct {
	bin []byte
	str string
}

var binTestFilters = []binTestFilter{
	{bin: []byte{0x87, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72}, str: "(member=*)"},
}

func TestFiltersDecode(t *testing.T) {
	for i, test := range binTestFilters {
		p := ber.DecodePacket(test.bin)
		if filter, err := DecompileFilter(p); err != nil {
			t.Errorf("binTestFilters[%d], DecompileFilter returned : %s", i, err)
		} else if filter != test.str {
			t.Errorf("binTestFilters[%d], %q expected, got %q", i, test.str, filter)
		}
	}
}

func TestFiltersEncode(t *testing.T) {
	for i, test := range binTestFilters {
		p, err := CompileFilter(test.str)
		if err != nil {
			t.Errorf("binTestFilters[%d], CompileFilter returned : %s", i, err)
			continue
		}
		b := p.Bytes()
		if !reflect.DeepEqual(b, test.bin) {
			t.Errorf("binTestFilters[%d], %q expected for CompileFilter(%q), got %q", i, test.bin, test.str, b)
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
		CompileFilter(filters[i%maxIdx])
	}
}

func BenchmarkFilterDecompile(b *testing.B) {
	b.StopTimer()
	filters := make([]*ber.Packet, len(testFilters))

	// Test Compiler and Decompiler
	for idx, i := range testFilters {
		filters[idx], _ = CompileFilter(i.filterStr)
	}

	maxIdx := len(filters)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DecompileFilter(filters[i%maxIdx])
	}
}
