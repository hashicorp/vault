package radix

import (
	"fmt"
)

// this type isn't _strictly_ necessary, but it helps to make the tests
// deterministic.

// poolConnColl implements a fixed-capacity collection of poolConns which
// supports random access and push/pop operations.
type poolConnColl struct {
	s        []*poolConn
	first, l int
}

func newPoolConnColl(capacity int) *poolConnColl {
	if capacity == 0 {
		panic("pool conn collection capacity of 0 is not supported")
	}
	return &poolConnColl{
		s: make([]*poolConn, capacity),
	}
}

func (p *poolConnColl) len() int {
	return p.l
}

func (p *poolConnColl) get(i int) *poolConn {
	if i >= p.l {
		panic(fmt.Sprintf("cannot get index %d, out of range. first:%d len:%d cap:%d", i, p.first, p.l, len(p.s)))
	}
	i += p.first
	if l := len(p.s); i >= l {
		i -= l
	}
	return p.s[i]
}

func (p *poolConnColl) remove(conn *poolConn) {
	if p.len() == 0 {
		return
	}

	var found bool
	i, c := p.first, 0
	for {
		if c == p.l {
			break
		} else if found = p.s[i] == conn; found {
			break
		}

		i++
		c++
		if i == len(p.s) {
			i = 0
		}
	}
	if !found {
		return
	}

	if i < p.first {
		copy(p.s[i:], p.s[i+1:p.first])
		p.s[p.first-1] = nil
	} else {
		copy(p.s[p.first+1:], p.s[p.first:i])
		p.s[p.first] = nil
		p.first++
		if p.first == len(p.s) {
			p.first = 0
		}
	}
	p.l--
}

func (p *poolConnColl) pushFront(c *poolConn) {
	if len(p.s) == p.l {
		panic("cannot push onto front of poolConnColl, it is already full")
	}

	i := p.first - 1
	if i == -1 {
		i = len(p.s) - 1
	}
	if p.s[i] != nil {
		panic(fmt.Sprintf("pushing onto front at %d but it already has an element", i))
	}
	p.s[i] = c
	p.first = i
	p.l++
}

func (p *poolConnColl) popBack() *poolConn {
	if p.len() == 0 {
		return nil
	}

	i := p.first + p.l - 1
	if i >= len(p.s) {
		i -= len(p.s)
	}

	c := p.s[i]
	p.s[i] = nil
	p.l--
	return c
}
