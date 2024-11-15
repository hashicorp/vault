package protocol

import (
	"sort"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

const noFieldName uint32 = 0xFFFFFFFF

type ofsName struct {
	ofs  uint32
	name string
}

type fieldNames struct { // use struct here to get a stable pointer
	items []ofsName
}

func (fn *fieldNames) search(ofs uint32) int {
	// binary search
	return sort.Search(len(fn.items), func(i int) bool { return fn.items[i].ofs >= ofs })
}

func (fn *fieldNames) insert(ofs uint32) {
	if ofs == noFieldName {
		return
	}
	i := fn.search(ofs)
	switch {
	case i >= len(fn.items): // not found -> append
		fn.items = append(fn.items, ofsName{ofs: ofs})
	case fn.items[i].ofs == ofs: // duplicate
	default: // insert
		fn.items = append(fn.items, ofsName{})
		copy(fn.items[i+1:], fn.items[i:])
		fn.items[i] = ofsName{ofs: ofs}
	}
}

func (fn *fieldNames) name(ofs uint32) string {
	i := fn.search(ofs)
	if i < len(fn.items) {
		return fn.items[i].name
	}
	return ""
}

func (fn *fieldNames) decode(dec *encoding.Decoder) (err error) {
	// TODO sniffer - python client texts are returned differently?
	// - double check offset calc (CESU8 issue?)
	pos := uint32(0)
	for i, on := range fn.items {
		diff := int(on.ofs - pos)
		if diff > 0 {
			dec.Skip(diff)
		}
		var n int
		var s string
		n, s, err = dec.CESU8LIString()
		fn.items[i].name = s
		// len byte + size + diff
		pos += uint32(n + diff) //nolint: gosec
	}
	return err
}
