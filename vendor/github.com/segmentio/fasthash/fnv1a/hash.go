package fnv1a

const (
	// FNV-1a
	offset64 = uint64(14695981039346656037)
	prime64  = uint64(1099511628211)

	// Init64 is what 64 bits hash values should be initialized with.
	Init64 = offset64
)

// HashString64 returns the hash of s.
func HashString64(s string) uint64 {
	return AddString64(Init64, s)
}

// HashBytes64 returns the hash of u.
func HashBytes64(b []byte) uint64 {
	return AddBytes64(Init64, b)
}

// HashUint64 returns the hash of u.
func HashUint64(u uint64) uint64 {
	return AddUint64(Init64, u)
}

// AddString64 adds the hash of s to the precomputed hash value h.
func AddString64(h uint64, s string) uint64 {
	/*
		This is an unrolled version of this algorithm:

		for _, c := range s {
			h = (h ^ uint64(c)) * prime64
		}

		It seems to be ~1.5x faster than the simple loop in BenchmarkHash64:

		- BenchmarkHash64/hash_function-4   30000000   56.1 ns/op   642.15 MB/s   0 B/op   0 allocs/op
		- BenchmarkHash64/hash_function-4   50000000   38.6 ns/op   932.35 MB/s   0 B/op   0 allocs/op

	*/
	for len(s) >= 8 {
		h = (h ^ uint64(s[0])) * prime64
		h = (h ^ uint64(s[1])) * prime64
		h = (h ^ uint64(s[2])) * prime64
		h = (h ^ uint64(s[3])) * prime64
		h = (h ^ uint64(s[4])) * prime64
		h = (h ^ uint64(s[5])) * prime64
		h = (h ^ uint64(s[6])) * prime64
		h = (h ^ uint64(s[7])) * prime64
		s = s[8:]
	}

	if len(s) >= 4 {
		h = (h ^ uint64(s[0])) * prime64
		h = (h ^ uint64(s[1])) * prime64
		h = (h ^ uint64(s[2])) * prime64
		h = (h ^ uint64(s[3])) * prime64
		s = s[4:]
	}

	if len(s) >= 2 {
		h = (h ^ uint64(s[0])) * prime64
		h = (h ^ uint64(s[1])) * prime64
		s = s[2:]
	}

	if len(s) > 0 {
		h = (h ^ uint64(s[0])) * prime64
	}

	return h
}

// AddBytes64 adds the hash of b to the precomputed hash value h.
func AddBytes64(h uint64, b []byte) uint64 {
	for len(b) >= 8 {
		h = (h ^ uint64(b[0])) * prime64
		h = (h ^ uint64(b[1])) * prime64
		h = (h ^ uint64(b[2])) * prime64
		h = (h ^ uint64(b[3])) * prime64
		h = (h ^ uint64(b[4])) * prime64
		h = (h ^ uint64(b[5])) * prime64
		h = (h ^ uint64(b[6])) * prime64
		h = (h ^ uint64(b[7])) * prime64
		b = b[8:]
	}

	if len(b) >= 4 {
		h = (h ^ uint64(b[0])) * prime64
		h = (h ^ uint64(b[1])) * prime64
		h = (h ^ uint64(b[2])) * prime64
		h = (h ^ uint64(b[3])) * prime64
		b = b[4:]
	}

	if len(b) >= 2 {
		h = (h ^ uint64(b[0])) * prime64
		h = (h ^ uint64(b[1])) * prime64
		b = b[2:]
	}

	if len(b) > 0 {
		h = (h ^ uint64(b[0])) * prime64
	}

	return h
}

// AddUint64 adds the hash value of the 8 bytes of u to h.
func AddUint64(h uint64, u uint64) uint64 {
	h = (h ^ ((u >> 56) & 0xFF)) * prime64
	h = (h ^ ((u >> 48) & 0xFF)) * prime64
	h = (h ^ ((u >> 40) & 0xFF)) * prime64
	h = (h ^ ((u >> 32) & 0xFF)) * prime64
	h = (h ^ ((u >> 24) & 0xFF)) * prime64
	h = (h ^ ((u >> 16) & 0xFF)) * prime64
	h = (h ^ ((u >> 8) & 0xFF)) * prime64
	h = (h ^ ((u >> 0) & 0xFF)) * prime64
	return h
}
