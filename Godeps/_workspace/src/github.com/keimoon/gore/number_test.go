package gore

import (
	"testing"
)

func TestFixInt(t *testing.T) {
	x := int64(-1234567812345678)
	b := FixInt(x).Bytes()
	y, _ := ToFixInt(b)
	if x != y {
		t.Fatal(x, y)
	}
}

func TestVarInt(t *testing.T) {
	x := int64(-123456789123456789)
	b := VarInt(x).Bytes()
	y, _ := ToVarInt(b)
	if x != y {
		t.Fatal(x, y)
	}
}

func BenchmarkVarInt(b *testing.B) {
	x := VarInt(-1234567812345678)
	for i := 0; i < b.N; i++ {
		ToVarInt(x.Bytes())
	}
}
