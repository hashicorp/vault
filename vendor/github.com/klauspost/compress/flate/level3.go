package flate

import "fmt"

// fastEncL3
type fastEncL3 struct {
	fastGen
	table [1 << 16]tableEntryPrev
}

// Encode uses a similar algorithm to level 2, will check up to two candidates.
func (e *fastEncL3) Encode(dst *tokens, src []byte) {
	const (
		inputMargin            = 12 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
		tableBits              = 16
		tableSize              = 1 << tableBits
		hashBytes              = 5
	)

	if debugDeflate && e.cur < 0 {
		panic(fmt.Sprint("e.cur < 0: ", e.cur))
	}

	// Protect against e.cur wraparound.
	for e.cur >= bufferReset {
		if len(e.hist) == 0 {
			for i := range e.table[:] {
				e.table[i] = tableEntryPrev{}
			}
			e.cur = maxMatchOffset
			break
		}
		// Shift down everything in the table that isn't already too far away.
		minOff := e.cur + int32(len(e.hist)) - maxMatchOffset
		for i := range e.table[:] {
			v := e.table[i]
			if v.Cur.offset <= minOff {
				v.Cur.offset = 0
			} else {
				v.Cur.offset = v.Cur.offset - e.cur + maxMatchOffset
			}
			if v.Prev.offset <= minOff {
				v.Prev.offset = 0
			} else {
				v.Prev.offset = v.Prev.offset - e.cur + maxMatchOffset
			}
			e.table[i] = v
		}
		e.cur = maxMatchOffset
	}

	s := e.addBlock(src)

	// Skip if too small.
	if len(src) < minNonLiteralBlockSize {
		// We do not fill the token table.
		// This will be picked up by caller.
		dst.n = uint16(len(src))
		return
	}

	// Override src
	src = e.hist
	nextEmit := s

	// sLimit is when to stop looking for offset/length copies. The inputMargin
	// lets us use a fast path for emitLiteral in the main loop, while we are
	// looking for copies.
	sLimit := int32(len(src) - inputMargin)

	// nextEmit is where in src the next emitLiteral should start from.
	cv := load6432(src, s)
	for {
		const skipLog = 7
		nextS := s
		var candidate tableEntry
		for {
			nextHash := hashLen(cv, tableBits, hashBytes)
			s = nextS
			nextS = s + 1 + (s-nextEmit)>>skipLog
			if nextS > sLimit {
				goto emitRemainder
			}
			candidates := e.table[nextHash]
			now := load6432(src, nextS)

			// Safe offset distance until s + 4...
			minOffset := e.cur + s - (maxMatchOffset - 4)
			e.table[nextHash] = tableEntryPrev{Prev: candidates.Cur, Cur: tableEntry{offset: s + e.cur}}

			// Check both candidates
			candidate = candidates.Cur
			if candidate.offset < minOffset {
				cv = now
				// Previous will also be invalid, we have nothing.
				continue
			}

			if uint32(cv) == load3232(src, candidate.offset-e.cur) {
				if candidates.Prev.offset < minOffset || uint32(cv) != load3232(src, candidates.Prev.offset-e.cur) {
					break
				}
				// Both match and are valid, pick longest.
				offset := s - (candidate.offset - e.cur)
				o2 := s - (candidates.Prev.offset - e.cur)
				l1, l2 := matchLen(src[s+4:], src[s-offset+4:]), matchLen(src[s+4:], src[s-o2+4:])
				if l2 > l1 {
					candidate = candidates.Prev
				}
				break
			} else {
				// We only check if value mismatches.
				// Offset will always be invalid in other cases.
				candidate = candidates.Prev
				if candidate.offset > minOffset && uint32(cv) == load3232(src, candidate.offset-e.cur) {
					break
				}
			}
			cv = now
		}

		// Call emitCopy, and then see if another emitCopy could be our next
		// move. Repeat until we find no match for the input immediately after
		// what was consumed by the last emitCopy call.
		//
		// If we exit this loop normally then we need to call emitLiteral next,
		// though we don't yet know how big the literal will be. We handle that
		// by proceeding to the next iteration of the main loop. We also can
		// exit this loop via goto if we get close to exhausting the input.
		for {
			// Invariant: we have a 4-byte match at s, and no need to emit any
			// literal bytes prior to s.

			// Extend the 4-byte match as long as possible.
			//
			t := candidate.offset - e.cur
			l := e.matchlenLong(s+4, t+4, src) + 4

			// Extend backwards
			for t > 0 && s > nextEmit && src[t-1] == src[s-1] {
				s--
				t--
				l++
			}
			if nextEmit < s {
				if false {
					emitLiteral(dst, src[nextEmit:s])
				} else {
					for _, v := range src[nextEmit:s] {
						dst.tokens[dst.n] = token(v)
						dst.litHist[v]++
						dst.n++
					}
				}
			}

			dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
			s += l
			nextEmit = s
			if nextS >= s {
				s = nextS + 1
			}

			if s >= sLimit {
				t += l
				// Index first pair after match end.
				if int(t+8) < len(src) && t > 0 {
					cv = load6432(src, t)
					nextHash := hashLen(cv, tableBits, hashBytes)
					e.table[nextHash] = tableEntryPrev{
						Prev: e.table[nextHash].Cur,
						Cur:  tableEntry{offset: e.cur + t},
					}
				}
				goto emitRemainder
			}

			// Store every 5th hash in-between.
			for i := s - l + 2; i < s-5; i += 6 {
				nextHash := hashLen(load6432(src, i), tableBits, hashBytes)
				e.table[nextHash] = tableEntryPrev{
					Prev: e.table[nextHash].Cur,
					Cur:  tableEntry{offset: e.cur + i}}
			}
			// We could immediately start working at s now, but to improve
			// compression we first update the hash table at s-2 to s.
			x := load6432(src, s-2)
			prevHash := hashLen(x, tableBits, hashBytes)

			e.table[prevHash] = tableEntryPrev{
				Prev: e.table[prevHash].Cur,
				Cur:  tableEntry{offset: e.cur + s - 2},
			}
			x >>= 8
			prevHash = hashLen(x, tableBits, hashBytes)

			e.table[prevHash] = tableEntryPrev{
				Prev: e.table[prevHash].Cur,
				Cur:  tableEntry{offset: e.cur + s - 1},
			}
			x >>= 8
			currHash := hashLen(x, tableBits, hashBytes)
			candidates := e.table[currHash]
			cv = x
			e.table[currHash] = tableEntryPrev{
				Prev: candidates.Cur,
				Cur:  tableEntry{offset: s + e.cur},
			}

			// Check both candidates
			candidate = candidates.Cur
			minOffset := e.cur + s - (maxMatchOffset - 4)

			if candidate.offset > minOffset {
				if uint32(cv) == load3232(src, candidate.offset-e.cur) {
					// Found a match...
					continue
				}
				candidate = candidates.Prev
				if candidate.offset > minOffset && uint32(cv) == load3232(src, candidate.offset-e.cur) {
					// Match at prev...
					continue
				}
			}
			cv = x >> 8
			s++
			break
		}
	}

emitRemainder:
	if int(nextEmit) < len(src) {
		// If nothing was added, don't encode literals.
		if dst.n == 0 {
			return
		}

		emitLiteral(dst, src[nextEmit:])
	}
}
