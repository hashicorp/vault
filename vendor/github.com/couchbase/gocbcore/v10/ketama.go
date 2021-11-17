package gocbcore

import (
	"crypto/md5" // nolint: gosec
	"fmt"
	"sort"
)

// "Point" in the ring hash entry. See lcbvb_CONTINUUM
type routeKetamaContinuum struct {
	index uint32
	point uint32
}

type ketamaSorter struct {
	elems []routeKetamaContinuum
}

func (c ketamaSorter) Len() int           { return len(c.elems) }
func (c ketamaSorter) Swap(i, j int)      { c.elems[i], c.elems[j] = c.elems[j], c.elems[i] }
func (c ketamaSorter) Less(i, j int) bool { return c.elems[i].point < c.elems[j].point }

type ketamaContinuum struct {
	entries []routeKetamaContinuum
}

func ketamaHash(key []byte) uint32 {
	digest := md5.Sum(key) // nolint: gosec

	return ((uint32(digest[3])&0xFF)<<24 |
		(uint32(digest[2])&0xFF)<<16 |
		(uint32(digest[1])&0xFF)<<8 |
		(uint32(digest[0]) & 0xFF)) & 0xffffffff
}

func newKetamaContinuum(serverList []string) *ketamaContinuum {
	continuum := ketamaContinuum{}

	// Libcouchbase presorts this. Might not strictly be required..
	sort.Strings(serverList)

	for ss, authority := range serverList {
		// 160 points per server
		for hh := 0; hh < 40; hh++ {
			hostkey := []byte(fmt.Sprintf("%s-%d", authority, hh))
			digest := md5.Sum(hostkey) // nolint: gosec

			for nn := 0; nn < 4; nn++ {

				var d1 = uint32(digest[3+nn*4]&0xff) << 24
				var d2 = uint32(digest[2+nn*4]&0xff) << 16
				var d3 = uint32(digest[1+nn*4]&0xff) << 8
				var d4 = uint32(digest[0+nn*4] & 0xff)
				var point = d1 | d2 | d3 | d4

				continuum.entries = append(continuum.entries, routeKetamaContinuum{
					point: point,
					index: uint32(ss),
				})
			}
		}
	}

	sort.Sort(ketamaSorter{continuum.entries})

	return &continuum
}

func (continuum ketamaContinuum) IsValid() bool {
	return len(continuum.entries) > 0
}

func (continuum ketamaContinuum) nodeByHash(hash uint32) (int, error) {
	var lowp = uint32(0)
	var highp = uint32(len(continuum.entries))
	var maxp = highp

	if len(continuum.entries) <= 0 {
		logErrorf("0-length ketama map!  Mapping to node 0.")
		return 0, errCliInternalError
	}

	// Copied from libcouchbase vbucket.c (map_ketama)
	for {
		midp := lowp + (highp-lowp)/2
		if midp == maxp {
			// Roll over to first entry
			return int(continuum.entries[0].index), nil
		}

		mid := continuum.entries[midp].point
		var prev uint32
		if midp == 0 {
			prev = 0
		} else {
			prev = continuum.entries[midp-1].point
		}

		if hash <= mid && hash > prev {
			return int(continuum.entries[midp].index), nil
		}

		if mid < hash {
			lowp = midp + 1
		} else {
			highp = midp - 1
		}

		if lowp > highp {
			return int(continuum.entries[0].index), nil
		}
	}
}

func (continuum ketamaContinuum) NodeByKey(key []byte) (int, error) {
	return continuum.nodeByHash(ketamaHash(key))
}
