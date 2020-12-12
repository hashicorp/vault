// +build appengine s390x

package murmur

import "encoding/binary"

func getBlock(data []byte, n int) (int64, int64) {
	k1 := int64(binary.LittleEndian.Uint64(data[n*16:]))
	k2 := int64(binary.LittleEndian.Uint64(data[(n*16)+8:]))
	return k1, k2
}
