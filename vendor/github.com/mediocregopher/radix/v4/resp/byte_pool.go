package resp

import "io"

// bytePool is a non-thread-safe pool of []byte instances which can help absorb
// allocations.
//
// bytePool uses an internal threshold and a counter to free []bytes on a
// regular basis. Everytime a []byte is put back in the pool the counter is
// incremented by len([]byte). If the counter then exceeds the threshold that
// []byte is freed rather than being put back.
type bytePool struct {
	used, threshold int
	pool            []*[]byte
}

func newBytePool(threshold int) *bytePool {
	return &bytePool{threshold: threshold}
}

func (bp *bytePool) get() *[]byte {
	if len(bp.pool) == 0 {
		b := make([]byte, 0, 32)
		return &b
	}
	b := bp.pool[len(bp.pool)-1]
	bp.pool = bp.pool[:len(bp.pool)-1]
	return b
}

func (bp *bytePool) put(b *[]byte) {
	if bp.used += cap(*b); bp.used > bp.threshold {
		bp.used = 0
		return
	}

	*b = (*b)[:0]
	bp.pool = append(bp.pool, b)
}

type byteReader struct {
	b    []byte
	pool *byteReaderPool
}

func (br *byteReader) Read(b []byte) (int, error) {
	if len(br.b) == 0 {
		return 0, io.EOF
	} else if len(b) == 0 {
		return 0, nil
	}

	n := copy(b, br.b)
	br.b = br.b[n:]
	if len(br.b) == 0 {
		br.b = nil
		br.pool.put(br)
		return n, io.EOF
	}
	return n, nil
}

type byteReaderPool struct {
	pool []*byteReader
}

func newByteReaderPool() *byteReaderPool {
	return new(byteReaderPool)
}

func (brp *byteReaderPool) get(b []byte) io.Reader {
	if len(brp.pool) == 0 {
		return &byteReader{b: b, pool: brp}
	}
	br := brp.pool[len(brp.pool)-1]
	brp.pool = brp.pool[:len(brp.pool)-1]
	br.b = b
	return br
}

func (brp *byteReaderPool) put(br *byteReader) {
	brp.pool = append(brp.pool, br)
}
