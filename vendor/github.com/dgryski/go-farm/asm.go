// +build ignore

package main

import (
	"flag"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

const k0 uint64 = 0xc3a5c85c97cb3127
const k1 uint64 = 0xb492b66fbe98f273
const k2 uint64 = 0x9ae16a3b2f90404f

const c1 uint32 = 0xcc9e2d51
const c2 uint32 = 0x1b873593

func shiftMix(val GPVirtual) GPVirtual {
	r := GP64()
	MOVQ(val, r)
	SHRQ(Imm(47), r)
	XORQ(val, r)
	return r
}

func shiftMix64(val uint64) uint64 {
	return val ^ (val >> 47)
}

func hashLen16MulLine(a, b, c, d, k, mul GPVirtual) GPVirtual {
	tmpa := GP64()
	MOVQ(a, tmpa)

	ADDQ(b, tmpa)
	RORQ(Imm(43), tmpa)
	ADDQ(d, tmpa)
	tmpc := GP64()
	MOVQ(c, tmpc)
	RORQ(Imm(30), tmpc)
	ADDQ(tmpc, tmpa)

	ADDQ(c, a)
	ADDQ(k, b)
	RORQ(Imm(18), b)
	ADDQ(b, a)

	r := hashLen16Mul(tmpa, a, mul)
	return r
}

func hashLen16Mul(u, v, mul GPVirtual) GPVirtual {
	XORQ(v, u)
	IMULQ(mul, u)
	a := shiftMix(u)

	XORQ(a, v)
	IMULQ(mul, v)
	b := shiftMix(v)

	IMULQ(mul, b)

	return b
}

func hashLen0to16(sbase, slen GPVirtual) {
	CMPQ(slen, Imm(8))
	JL(LabelRef("check4"))
	{
		a := GP64()
		MOVQ(Mem{Base: sbase}, a)

		b := GP64()
		t := GP64()
		MOVQ(slen, t)
		SUBQ(Imm(8), t)
		ADDQ(sbase, t)
		MOVQ(Mem{Base: t}, b)

		rk2 := GP64()
		MOVQ(Imm(k2), rk2)

		ADDQ(rk2, a)

		mul := slen
		SHLQ(Imm(1), mul)
		ADDQ(rk2, mul)

		c := GP64()
		MOVQ(b, c)
		RORQ(Imm(37), c)
		IMULQ(mul, c)
		ADDQ(a, c)

		d := GP64()
		MOVQ(a, d)
		RORQ(Imm(25), d)
		ADDQ(b, d)
		IMULQ(mul, d)

		r := hashLen16Mul(c, d, mul)
		Store(r, ReturnIndex(0))
		RET()
	}

	Label("check4")

	CMPQ(slen, Imm(4))
	JL(LabelRef("check0"))
	{
		rk2 := GP64()
		MOVQ(Imm(k2), rk2)

		mul := GP64()
		MOVQ(slen, mul)
		SHLQ(Imm(1), mul)
		ADDQ(rk2, mul)

		a := GP64()
		MOVL(Mem{Base: sbase}, a.As32())

		SHLQ(Imm(3), a)
		ADDQ(slen, a)

		b := GP64()
		SUBQ(Imm(4), slen)
		ADDQ(slen, sbase)
		MOVL(Mem{Base: sbase}, b.As32())
		r := hashLen16Mul(a, b, mul)

		Store(r, ReturnIndex(0))
		RET()
	}

	Label("check0")
	TESTQ(slen, slen)
	JZ(LabelRef("empty"))
	{

		a := GP64()
		MOVBQZX(Mem{Base: sbase}, a)

		base := GP64()
		MOVQ(slen, base)
		SHRQ(Imm(1), base)

		b := GP64()
		ADDQ(sbase, base)
		MOVBQZX(Mem{Base: base}, b)

		MOVQ(slen, base)
		SUBQ(Imm(1), base)
		c := GP64()
		ADDQ(sbase, base)
		MOVBQZX(Mem{Base: base}, c)

		SHLQ(Imm(8), b)
		ADDQ(b, a)
		y := a

		SHLQ(Imm(2), c)
		ADDQ(c, slen)
		z := slen

		rk0 := GP64()
		MOVQ(Imm(k0), rk0)
		IMULQ(rk0, z)

		rk2 := GP64()
		MOVQ(Imm(k2), rk2)

		IMULQ(rk2, y)
		XORQ(y, z)

		r := shiftMix(z)
		IMULQ(rk2, r)

		Store(r, ReturnIndex(0))
		RET()
	}

	Label("empty")

	ret := GP64()
	MOVQ(Imm(k2), ret)
	Store(ret, ReturnIndex(0))
	RET()
}

func hashLen17to32(sbase, slen GPVirtual) {
	mul := GP64()
	MOVQ(slen, mul)
	SHLQ(Imm(1), mul)

	rk2 := GP64()
	MOVQ(Imm(k2), rk2)
	ADDQ(rk2, mul)

	a := GP64()
	MOVQ(Mem{Base: sbase}, a)

	rk1 := GP64()
	MOVQ(Imm(k1), rk1)
	IMULQ(rk1, a)

	b := GP64()
	MOVQ(Mem{Base: sbase, Disp: 8}, b)

	base := GP64()
	MOVQ(slen, base)
	SUBQ(Imm(16), base)
	ADDQ(sbase, base)

	c := GP64()
	MOVQ(Mem{Base: base, Disp: 8}, c)
	IMULQ(mul, c)

	d := GP64()
	MOVQ(Mem{Base: base}, d)
	IMULQ(rk2, d)

	r := hashLen16MulLine(a, b, c, d, rk2, mul)
	Store(r, ReturnIndex(0))
	RET()
}

// Return an 8-byte hash for 33 to 64 bytes.
func hashLen33to64(sbase, slen GPVirtual) {
	mul := GP64()
	MOVQ(slen, mul)
	SHLQ(Imm(1), mul)

	rk2 := GP64()
	MOVQ(Imm(k2), rk2)
	ADDQ(rk2, mul)

	a := GP64()
	MOVQ(Mem{Base: sbase}, a)
	IMULQ(rk2, a)

	b := GP64()
	MOVQ(Mem{Base: sbase, Disp: 8}, b)

	base := GP64()
	MOVQ(slen, base)
	SUBQ(Imm(16), base)
	ADDQ(sbase, base)

	c := GP64()
	MOVQ(Mem{Base: base, Disp: 8}, c)
	IMULQ(mul, c)

	d := GP64()
	MOVQ(Mem{Base: base}, d)
	IMULQ(rk2, d)

	y := GP64()
	MOVQ(a, y)

	ADDQ(b, y)
	RORQ(Imm(43), y)
	ADDQ(d, y)
	tmpc := GP64()
	MOVQ(c, tmpc)
	RORQ(Imm(30), tmpc)
	ADDQ(tmpc, y)

	ADDQ(a, c)
	ADDQ(rk2, b)
	RORQ(Imm(18), b)
	ADDQ(b, c)

	tmpy := GP64()
	MOVQ(y, tmpy)
	z := hashLen16Mul(tmpy, c, mul)

	e := GP64()
	MOVQ(Mem{Base: sbase, Disp: 16}, e)
	IMULQ(mul, e)

	f := GP64()
	MOVQ(Mem{Base: sbase, Disp: 24}, f)

	base = GP64()
	MOVQ(slen, base)
	SUBQ(Imm(32), base)
	ADDQ(sbase, base)
	g := GP64()
	MOVQ(Mem{Base: base}, g)
	ADDQ(y, g)
	IMULQ(mul, g)

	h := GP64()
	MOVQ(Mem{Base: base, Disp: 8}, h)
	ADDQ(z, h)
	IMULQ(mul, h)

	r := hashLen16MulLine(e, f, g, h, a, mul)
	Store(r, ReturnIndex(0))
	RET()
}

// Return a 16-byte hash for s[0] ... s[31], a, and b.  Quick and dirty.
func weakHashLen32WithSeeds(sbase GPVirtual, disp int, a, b GPVirtual) {

	w := Mem{Base: sbase, Disp: disp + 0}
	x := Mem{Base: sbase, Disp: disp + 8}
	y := Mem{Base: sbase, Disp: disp + 16}
	z := Mem{Base: sbase, Disp: disp + 24}

	// a += w
	ADDQ(w, a)

	// b = bits.RotateLeft64(b+a+z, -21)
	ADDQ(a, b)
	ADDQ(z, b)
	RORQ(Imm(21), b)

	// c := a
	c := GP64()
	MOVQ(a, c)

	// a += x
	// a += y
	ADDQ(x, a)
	ADDQ(y, a)

	// b += bits.RotateLeft64(a, -44)
	atmp := GP64()
	MOVQ(a, atmp)
	RORQ(Imm(44), atmp)
	ADDQ(atmp, b)

	// a += z
	// b += c
	ADDQ(z, a)
	ADDQ(c, b)

	XCHGQ(a, b)
}

func hashLoopBody(x, y, z, vlo, vhi, wlo, whi, sbase GPVirtual, mul1 GPVirtual, mul2 uint64) {
	ADDQ(y, x)
	ADDQ(vlo, x)
	ADDQ(Mem{Base: sbase, Disp: 8}, x)
	RORQ(Imm(37), x)

	IMULQ(mul1, x)

	ADDQ(vhi, y)
	ADDQ(Mem{Base: sbase, Disp: 48}, y)
	RORQ(Imm(42), y)
	IMULQ(mul1, y)

	if mul2 != 1 {
		t := GP64()
		MOVQ(U32(mul2), t)
		IMULQ(whi, t)
		XORQ(t, x)
	} else {
		XORQ(whi, x)
	}

	if mul2 != 1 {
		t := GP64()
		MOVQ(U32(mul2), t)
		IMULQ(vlo, t)
		ADDQ(t, y)
	} else {
		ADDQ(vlo, y)
	}

	ADDQ(Mem{Base: sbase, Disp: 40}, y)

	ADDQ(wlo, z)
	RORQ(Imm(33), z)
	IMULQ(mul1, z)

	{
		IMULQ(mul1, vhi)
		MOVQ(x, vlo)
		ADDQ(wlo, vlo)
		weakHashLen32WithSeeds(sbase, 0, vhi, vlo)
	}

	{
		ADDQ(z, whi)
		MOVQ(y, wlo)
		ADDQ(Mem{Base: sbase, Disp: 16}, wlo)
		weakHashLen32WithSeeds(sbase, 32, whi, wlo)
	}

	XCHGQ(z, x)
}

func fp64() {

	TEXT("Fingerprint64", NOSPLIT, "func(s []byte) uint64")

	slen := GP64()
	sbase := GP64()

	Load(Param("s").Base(), sbase)
	Load(Param("s").Len(), slen)

	CMPQ(slen, Imm(16))
	JG(LabelRef("check32"))
	hashLen0to16(sbase, slen)

	Label("check32")
	CMPQ(slen, Imm(32))
	JG(LabelRef("check64"))
	hashLen17to32(sbase, slen)

	Label("check64")
	CMPQ(slen, Imm(64))
	JG(LabelRef("long"))
	hashLen33to64(sbase, slen)

	Label("long")

	seed := uint64(81)

	vlo, vhi, wlo, whi := GP64(), GP64(), GP64(), GP64()
	XORQ(vlo, vlo)
	XORQ(vhi, vhi)
	XORQ(wlo, wlo)
	XORQ(whi, whi)

	x := GP64()

	eightOne := uint64(81)

	MOVQ(Imm(eightOne*k2), x)
	ADDQ(Mem{Base: sbase}, x)

	y := GP64()
	y64 := uint64(seed*k1) + 113
	MOVQ(Imm(y64), y)

	z := GP64()
	MOVQ(Imm(shiftMix64(y64*k2+113)*k2), z)

	endIdx := GP64()
	MOVQ(slen, endIdx)
	tmp := GP64()
	SUBQ(Imm(1), endIdx)
	MOVQ(U64(^uint64(63)), tmp)
	ANDQ(tmp, endIdx)
	last64Idx := GP64()
	MOVQ(slen, last64Idx)
	SUBQ(Imm(1), last64Idx)
	ANDQ(Imm(63), last64Idx)
	SUBQ(Imm(63), last64Idx)
	ADDQ(endIdx, last64Idx)

	last64 := GP64()
	MOVQ(last64Idx, last64)
	ADDQ(sbase, last64)

	end := GP64()
	MOVQ(slen, end)

	Label("loop")

	rk1 := GP64()
	MOVQ(Imm(k1), rk1)

	hashLoopBody(x, y, z, vlo, vhi, wlo, whi, sbase, rk1, 1)

	ADDQ(Imm(64), sbase)
	SUBQ(Imm(64), end)
	CMPQ(end, Imm(64))
	JG(LabelRef("loop"))

	MOVQ(last64, sbase)

	mul := GP64()
	MOVQ(z, mul)
	ANDQ(Imm(0xff), mul)
	SHLQ(Imm(1), mul)
	ADDQ(rk1, mul)

	MOVQ(last64, sbase)

	SUBQ(Imm(1), slen)
	ANDQ(Imm(63), slen)
	ADDQ(slen, wlo)

	ADDQ(wlo, vlo)
	ADDQ(vlo, wlo)

	hashLoopBody(x, y, z, vlo, vhi, wlo, whi, sbase, mul, 9)

	{
		a := hashLen16Mul(vlo, wlo, mul)
		ADDQ(z, a)
		b := shiftMix(y)
		rk0 := GP64()
		MOVQ(Imm(k0), rk0)
		IMULQ(rk0, b)
		ADDQ(b, a)

		c := hashLen16Mul(vhi, whi, mul)
		ADDQ(x, c)

		r := hashLen16Mul(a, c, mul)
		Store(r, ReturnIndex(0))
	}

	RET()
}

func fmix(h GPVirtual) GPVirtual {
	h2 := GP32()
	MOVL(h, h2)
	SHRL(Imm(16), h2)
	XORL(h2, h)

	MOVL(Imm(0x85ebca6b), h2)
	IMULL(h2, h)

	MOVL(h, h2)
	SHRL(Imm(13), h2)
	XORL(h2, h)

	MOVL(Imm(0xc2b2ae35), h2)
	IMULL(h2, h)

	MOVL(h, h2)
	SHRL(Imm(16), h2)
	XORL(h2, h)
	return h
}

func mur(a, h GPVirtual) GPVirtual {
	imul3l(c1, a, a)
	RORL(Imm(17), a)
	imul3l(c2, a, a)
	XORL(a, h)
	RORL(Imm(19), h)

	LEAL(Mem{Base: h, Index: h, Scale: 4}, a)
	LEAL(Mem{Base: a, Disp: 0xe6546b64}, h)

	return h
}

func hash32Len5to12(sbase, slen GPVirtual) {

	a := GP32()
	MOVL(slen.As32(), a)
	b := GP32()
	MOVL(a, b)
	SHLL(Imm(2), b)
	ADDL(a, b)

	c := GP32()
	MOVL(U32(9), c)

	d := GP32()
	MOVL(b, d)

	ADDL(Mem{Base: sbase, Disp: 0}, a)

	t := GP64()
	MOVQ(slen, t)
	SUBQ(Imm(4), t)
	ADDQ(sbase, t)
	ADDL(Mem{Base: t}, b)

	MOVQ(slen, t)
	SHRQ(Imm(1), t)
	ANDQ(Imm(4), t)
	ADDQ(sbase, t)
	ADDL(Mem{Base: t}, c)

	t = mur(a, d)
	t = mur(b, t)
	t = mur(c, t)
	t = fmix(t)

	Store(t, ReturnIndex(0))
	RET()
}

func hash32Len13to24Seed(sbase, slen GPVirtual) {
	slen2 := GP64()
	MOVQ(slen, slen2)
	SHRQ(Imm(1), slen2)
	ADDQ(sbase, slen2)

	a := GP32()
	MOVL(Mem{Base: slen2, Disp: -4}, a)

	b := GP32()
	MOVL(Mem{Base: sbase, Disp: 4}, b)

	send := GP64()
	MOVQ(slen, send)
	ADDQ(sbase, send)

	c := GP32()
	MOVL(Mem{Base: send, Disp: -8}, c)

	d := GP32()
	MOVL(Mem{Base: slen2}, d)

	e := GP32()
	MOVL(Mem{Base: sbase}, e)

	f := GP32()
	MOVL(Mem{Base: send, Disp: -4}, f)

	h := GP32()
	MOVL(U32(c1), h)
	IMULL(d, h)
	ADDL(slen.As32(), h)

	RORL(Imm(12), a)
	ADDL(f, a)

	ctmp := GP32()
	MOVL(c, ctmp)
	h = mur(ctmp, h)
	ADDL(a, h)

	RORL(Imm(3), a)
	ADDL(c, a)

	h = mur(e, h)
	ADDL(a, h)

	ADDL(f, a)
	RORL(Imm(12), a)
	ADDL(d, a)

	h = mur(b, h)
	ADDL(a, h)

	h = fmix(h)

	Store(h, ReturnIndex(0))
	RET()
}

func hash32Len0to4(sbase, slen GPVirtual) {
	b := GP32()
	c := GP32()

	XORL(b, b)
	MOVL(U32(9), c)

	TESTQ(slen, slen)
	JZ(LabelRef("done"))

	l := GP64()
	v := GP32()
	MOVQ(slen, l)

	c1reg := GP32()
	MOVL(U32(c1), c1reg)

	for i := 0; i < 4; i++ {
		IMULL(c1reg, b)
		MOVBLSX(Mem{Base: sbase, Disp: i}, v)
		ADDL(v, b)
		XORL(b, c)
		SUBQ(Imm(1), l)
		TESTQ(l, l)
		JZ(LabelRef("done"))
	}

	Label("done")

	s32 := GP32()
	MOVL(slen.As32(), s32)
	r := mur(s32, c)
	r = mur(b, r)
	r = fmix(r)

	Store(r, ReturnIndex(0))
	RET()
}

func fp32() {

	TEXT("Fingerprint32", NOSPLIT, "func(s []byte) uint32")

	sbase := GP64()
	slen := GP64()

	Load(Param("s").Base(), sbase)
	Load(Param("s").Len(), slen)

	CMPQ(slen, Imm(24))
	JG(LabelRef("long"))

	CMPQ(slen, Imm(12))
	JG(LabelRef("hash_13_24"))

	CMPQ(slen, Imm(4))
	JG(LabelRef("hash_5_12"))
	hash32Len0to4(sbase, slen)

	Label("hash_5_12")
	hash32Len5to12(sbase, slen)

	Label("hash_13_24")
	hash32Len13to24Seed(sbase, slen)

	Label("long")

	h := GP32()
	MOVL(slen.As32(), h)

	g := GP32()
	MOVL(U32(c1), g)
	IMULL(h, g)

	f := GP32()
	MOVL(g, f)

	// len > 24

	send := GP64()
	MOVQ(slen, send)
	ADDQ(sbase, send)
	c1reg := GP32()
	MOVL(U32(c1), c1reg)
	c2reg := GP32()
	MOVL(U32(c2), c2reg)

	shuf := func(r GPVirtual, disp int) {
		a := GP32()
		MOVL(Mem{Base: send, Disp: disp}, a)
		IMULL(c1reg, a)
		RORL(Imm(17), a)
		IMULL(c2reg, a)
		XORL(a, r)
		RORL(Imm(19), r)
		MOVL(r, a)
		SHLL(Imm(2), a)
		ADDL(a, r)
		ADDL(Imm(0xe6546b64), r)
	}

	shuf(h, -4)
	shuf(g, -8)
	shuf(h, -16)
	shuf(g, -12)

	PREFETCHT0(Mem{Base: sbase})
	{
		a := GP32()
		MOVL(Mem{Base: send, Disp: -20}, a)
		IMULL(c1reg, a)
		RORL(Imm(17), a)
		IMULL(c2reg, a)

		ADDL(a, f)
		RORL(Imm(19), f)
		ADDL(Imm(113), f)

	}

	loop32Body := func(f, g, h, sbase, slen GPVirtual, disp int) {
		a, b, c, d, e := GP32(), GP32(), GP32(), GP32(), GP32()

		MOVL(Mem{Base: sbase, Disp: disp + 0}, a)
		ADDL(a, h)

		MOVL(Mem{Base: sbase, Disp: disp + 4}, b)
		ADDL(b, g)

		MOVL(Mem{Base: sbase, Disp: disp + 8}, c)
		ADDL(c, f)

		MOVL(Mem{Base: sbase, Disp: disp + 12}, d)
		t := GP32()
		MOVL(d, t)
		h = mur(t, h)

		MOVL(Mem{Base: sbase, Disp: disp + 16}, e)
		ADDL(e, h)

		MOVL(c, t)
		g = mur(t, g)
		ADDL(a, g)

		imul3l(c1, e, t)
		ADDL(b, t)
		f = mur(t, f)
		ADDL(d, f)

		ADDL(g, f)
		ADDL(f, g)
	}

	Label("loop80")
	CMPQ(slen, Imm(80+20))
	JL(LabelRef("loop20"))
	{
		PREFETCHT0(Mem{Base: sbase, Disp: 20})
		loop32Body(f, g, h, sbase, slen, 0)
		PREFETCHT0(Mem{Base: sbase, Disp: 40})
		loop32Body(f, g, h, sbase, slen, 20)
		PREFETCHT0(Mem{Base: sbase, Disp: 60})
		loop32Body(f, g, h, sbase, slen, 40)
		PREFETCHT0(Mem{Base: sbase, Disp: 80})
		loop32Body(f, g, h, sbase, slen, 60)

		ADDQ(Imm(80), sbase)
		SUBQ(Imm(80), slen)
		JMP(LabelRef("loop80"))
	}

	Label("loop20")
	CMPQ(slen, Imm(20))
	JLE(LabelRef("after"))
	{
		loop32Body(f, g, h, sbase, slen, 0)

		ADDQ(Imm(20), sbase)
		SUBQ(Imm(20), slen)
		JMP(LabelRef("loop20"))
	}

	Label("after")

	c1reg = GP32()
	MOVL(U32(c1), c1reg)

	RORL(Imm(11), g)
	IMULL(c1reg, g)

	RORL(Imm(17), g)
	IMULL(c1reg, g)

	RORL(Imm(11), f)
	IMULL(c1reg, f)

	RORL(Imm(17), f)
	IMULL(c1reg, f)

	ADDL(g, h)
	RORL(Imm(19), h)

	t := GP32()
	MOVL(h, t)
	SHLL(Imm(2), t)
	ADDL(t, h)
	ADDL(Imm(0xe6546b64), h)

	RORL(Imm(17), h)
	IMULL(c1reg, h)

	ADDL(f, h)
	RORL(Imm(19), h)

	t = GP32()
	MOVL(h, t)
	SHLL(Imm(2), t)
	ADDL(t, h)
	ADDL(Imm(0xe6546b64), h)

	RORL(Imm(17), h)
	IMULL(c1reg, h)

	Store(h, ReturnIndex(0))
	RET()
}

var go111 = flag.Bool("go111", true, "use assembly instructions present in go1.11 and later")

func imul3l(m uint32, x, y Register) {
	if *go111 {
		IMUL3L(U32(m), x, y)
	} else {
		t := GP32()
		MOVL(U32(m), t)
		IMULL(t, x)
		MOVL(x, y)
	}
}

func main() {

	flag.Parse()

	ConstraintExpr("amd64,!purego")

	fp64()
	fp32()

	Generate()
}
