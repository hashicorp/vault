//line parse.y:4
package hcl

import __yyfmt__ "fmt"

//line parse.y:4
import (
	"fmt"
	"strconv"
)

//line parse.y:13
type hclSymType struct {
	yys     int
	b       bool
	f       float64
	num     int
	str     string
	obj     *Object
	objlist []*Object
}

const BOOL = 57346
const FLOAT = 57347
const NUMBER = 57348
const COMMA = 57349
const COMMAEND = 57350
const IDENTIFIER = 57351
const EQUAL = 57352
const NEWLINE = 57353
const STRING = 57354
const MINUS = 57355
const LEFTBRACE = 57356
const RIGHTBRACE = 57357
const LEFTBRACKET = 57358
const RIGHTBRACKET = 57359
const PERIOD = 57360
const EPLUS = 57361
const EMINUS = 57362

var hclToknames = []string{
	"BOOL",
	"FLOAT",
	"NUMBER",
	"COMMA",
	"COMMAEND",
	"IDENTIFIER",
	"EQUAL",
	"NEWLINE",
	"STRING",
	"MINUS",
	"LEFTBRACE",
	"RIGHTBRACE",
	"LEFTBRACKET",
	"RIGHTBRACKET",
	"PERIOD",
	"EPLUS",
	"EMINUS",
}
var hclStatenames = []string{}

const hclEofCode = 1
const hclErrCode = 2
const hclMaxDepth = 200

//line parse.y:259

//line yacctab:1
var hclExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 6,
	10, 7,
	-2, 17,
	-1, 7,
	10, 8,
	-2, 18,
}

const hclNprod = 36
const hclPrivate = 57344

var hclTokenNames []string
var hclStates []string

const hclLast = 62

var hclAct = []int{

	35, 3, 21, 22, 9, 30, 31, 29, 17, 26,
	25, 26, 25, 10, 26, 25, 18, 24, 13, 24,
	23, 37, 24, 44, 45, 42, 34, 38, 39, 9,
	32, 6, 6, 43, 7, 7, 2, 40, 28, 26,
	25, 6, 41, 11, 7, 46, 37, 24, 14, 36,
	27, 15, 5, 13, 19, 1, 4, 8, 33, 20,
	16, 12,
}
var hclPact = []int{

	32, -1000, 32, -1000, 3, -1000, -1000, -1000, 39, -1000,
	4, -1000, -1000, 23, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -14, -14, 9, 6, -1000, -1000, 22, -1000, -1000,
	36, 19, -1000, 16, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 34, -1000, -1000,
}
var hclPgo = []int{

	0, 3, 2, 59, 58, 36, 52, 49, 43, 1,
	0, 57, 7, 56, 55,
}
var hclR1 = []int{

	0, 14, 14, 5, 5, 8, 8, 13, 13, 9,
	9, 9, 9, 9, 9, 6, 6, 11, 11, 3,
	3, 4, 4, 4, 10, 10, 7, 7, 7, 7,
	2, 2, 1, 1, 12, 12,
}
var hclR2 = []int{

	0, 0, 1, 1, 2, 3, 2, 1, 1, 3,
	3, 3, 3, 3, 1, 2, 2, 1, 1, 3,
	2, 1, 3, 2, 1, 1, 1, 1, 2, 2,
	2, 1, 2, 1, 2, 2,
}
var hclChk = []int{

	-1000, -14, -5, -9, -13, -6, 9, 12, -11, -9,
	10, -8, -6, 14, 9, 12, -7, 4, 12, -8,
	-3, -2, -1, 16, 13, 6, 5, -5, 15, -12,
	19, 20, -12, -4, 17, -10, -7, 12, -2, -1,
	15, 6, 6, 17, 7, 8, -10,
}
var hclDef = []int{

	1, -2, 2, 3, 0, 14, -2, -2, 0, 4,
	0, 15, 16, 0, 17, 18, 9, 10, 11, 12,
	13, 26, 27, 0, 0, 31, 33, 0, 6, 28,
	0, 0, 29, 0, 20, 21, 24, 25, 30, 32,
	5, 34, 35, 19, 0, 23, 22,
}
var hclTok1 = []int{

	1,
}
var hclTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20,
}
var hclTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var hclDebug = 0

type hclLexer interface {
	Lex(lval *hclSymType) int
	Error(s string)
}

const hclFlag = -1000

func hclTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(hclToknames) {
		if hclToknames[c-4] != "" {
			return hclToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func hclStatname(s int) string {
	if s >= 0 && s < len(hclStatenames) {
		if hclStatenames[s] != "" {
			return hclStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func hcllex1(lex hclLexer, lval *hclSymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = hclTok1[0]
		goto out
	}
	if char < len(hclTok1) {
		c = hclTok1[char]
		goto out
	}
	if char >= hclPrivate {
		if char < hclPrivate+len(hclTok2) {
			c = hclTok2[char-hclPrivate]
			goto out
		}
	}
	for i := 0; i < len(hclTok3); i += 2 {
		c = hclTok3[i+0]
		if c == char {
			c = hclTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = hclTok2[1] /* unknown char */
	}
	if hclDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", hclTokname(c), uint(char))
	}
	return c
}

func hclParse(hcllex hclLexer) int {
	var hcln int
	var hcllval hclSymType
	var hclVAL hclSymType
	hclS := make([]hclSymType, hclMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	hclstate := 0
	hclchar := -1
	hclp := -1
	goto hclstack

ret0:
	return 0

ret1:
	return 1

hclstack:
	/* put a state and value onto the stack */
	if hclDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", hclTokname(hclchar), hclStatname(hclstate))
	}

	hclp++
	if hclp >= len(hclS) {
		nyys := make([]hclSymType, len(hclS)*2)
		copy(nyys, hclS)
		hclS = nyys
	}
	hclS[hclp] = hclVAL
	hclS[hclp].yys = hclstate

hclnewstate:
	hcln = hclPact[hclstate]
	if hcln <= hclFlag {
		goto hcldefault /* simple state */
	}
	if hclchar < 0 {
		hclchar = hcllex1(hcllex, &hcllval)
	}
	hcln += hclchar
	if hcln < 0 || hcln >= hclLast {
		goto hcldefault
	}
	hcln = hclAct[hcln]
	if hclChk[hcln] == hclchar { /* valid shift */
		hclchar = -1
		hclVAL = hcllval
		hclstate = hcln
		if Errflag > 0 {
			Errflag--
		}
		goto hclstack
	}

hcldefault:
	/* default state action */
	hcln = hclDef[hclstate]
	if hcln == -2 {
		if hclchar < 0 {
			hclchar = hcllex1(hcllex, &hcllval)
		}

		/* look through exception table */
		xi := 0
		for {
			if hclExca[xi+0] == -1 && hclExca[xi+1] == hclstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			hcln = hclExca[xi+0]
			if hcln < 0 || hcln == hclchar {
				break
			}
		}
		hcln = hclExca[xi+1]
		if hcln < 0 {
			goto ret0
		}
	}
	if hcln == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			hcllex.Error("syntax error")
			Nerrs++
			if hclDebug >= 1 {
				__yyfmt__.Printf("%s", hclStatname(hclstate))
				__yyfmt__.Printf(" saw %s\n", hclTokname(hclchar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for hclp >= 0 {
				hcln = hclPact[hclS[hclp].yys] + hclErrCode
				if hcln >= 0 && hcln < hclLast {
					hclstate = hclAct[hcln] /* simulate a shift of "error" */
					if hclChk[hclstate] == hclErrCode {
						goto hclstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if hclDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", hclS[hclp].yys)
				}
				hclp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if hclDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", hclTokname(hclchar))
			}
			if hclchar == hclEofCode {
				goto ret1
			}
			hclchar = -1
			goto hclnewstate /* try again in the same state */
		}
	}

	/* reduction by production hcln */
	if hclDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", hcln, hclStatname(hclstate))
	}

	hclnt := hcln
	hclpt := hclp
	_ = hclpt // guard against "declared and not used"

	hclp -= hclR2[hcln]
	hclVAL = hclS[hclp+1]

	/* consult goto table to find next state */
	hcln = hclR1[hcln]
	hclg := hclPgo[hcln]
	hclj := hclg + hclS[hclp].yys + 1

	if hclj >= hclLast {
		hclstate = hclAct[hclg]
	} else {
		hclstate = hclAct[hclj]
		if hclChk[hclstate] != -hcln {
			hclstate = hclAct[hclg]
		}
	}
	// dummy call; replaced with literal code
	switch hclnt {

	case 1:
		//line parse.y:39
		{
			hclResult = &Object{Type: ValueTypeObject}
		}
	case 2:
		//line parse.y:43
		{
			hclResult = &Object{
				Type:  ValueTypeObject,
				Value: ObjectList(hclS[hclpt-0].objlist).Flat(),
			}
		}
	case 3:
		//line parse.y:52
		{
			hclVAL.objlist = []*Object{hclS[hclpt-0].obj}
		}
	case 4:
		//line parse.y:56
		{
			hclVAL.objlist = append(hclS[hclpt-1].objlist, hclS[hclpt-0].obj)
		}
	case 5:
		//line parse.y:62
		{
			hclVAL.obj = &Object{
				Type:  ValueTypeObject,
				Value: ObjectList(hclS[hclpt-1].objlist).Flat(),
			}
		}
	case 6:
		//line parse.y:69
		{
			hclVAL.obj = &Object{
				Type: ValueTypeObject,
			}
		}
	case 7:
		//line parse.y:77
		{
			hclVAL.str = hclS[hclpt-0].str
		}
	case 8:
		//line parse.y:81
		{
			hclVAL.str = hclS[hclpt-0].str
		}
	case 9:
		//line parse.y:87
		{
			hclVAL.obj = hclS[hclpt-0].obj
			hclVAL.obj.Key = hclS[hclpt-2].str
		}
	case 10:
		//line parse.y:92
		{
			hclVAL.obj = &Object{
				Key:   hclS[hclpt-2].str,
				Type:  ValueTypeBool,
				Value: hclS[hclpt-0].b,
			}
		}
	case 11:
		//line parse.y:100
		{
			hclVAL.obj = &Object{
				Key:   hclS[hclpt-2].str,
				Type:  ValueTypeString,
				Value: hclS[hclpt-0].str,
			}
		}
	case 12:
		//line parse.y:108
		{
			hclS[hclpt-0].obj.Key = hclS[hclpt-2].str
			hclVAL.obj = hclS[hclpt-0].obj
		}
	case 13:
		//line parse.y:113
		{
			hclVAL.obj = &Object{
				Key:   hclS[hclpt-2].str,
				Type:  ValueTypeList,
				Value: hclS[hclpt-0].objlist,
			}
		}
	case 14:
		//line parse.y:121
		{
			hclVAL.obj = hclS[hclpt-0].obj
		}
	case 15:
		//line parse.y:127
		{
			hclS[hclpt-0].obj.Key = hclS[hclpt-1].str
			hclVAL.obj = hclS[hclpt-0].obj
		}
	case 16:
		//line parse.y:132
		{
			hclVAL.obj = &Object{
				Key:   hclS[hclpt-1].str,
				Type:  ValueTypeObject,
				Value: []*Object{hclS[hclpt-0].obj},
			}
		}
	case 17:
		//line parse.y:142
		{
			hclVAL.str = hclS[hclpt-0].str
		}
	case 18:
		//line parse.y:146
		{
			hclVAL.str = hclS[hclpt-0].str
		}
	case 19:
		//line parse.y:152
		{
			hclVAL.objlist = hclS[hclpt-1].objlist
		}
	case 20:
		//line parse.y:156
		{
			hclVAL.objlist = nil
		}
	case 21:
		//line parse.y:162
		{
			hclVAL.objlist = []*Object{hclS[hclpt-0].obj}
		}
	case 22:
		//line parse.y:166
		{
			hclVAL.objlist = append(hclS[hclpt-2].objlist, hclS[hclpt-0].obj)
		}
	case 23:
		//line parse.y:170
		{
			hclVAL.objlist = hclS[hclpt-1].objlist
		}
	case 24:
		//line parse.y:176
		{
			hclVAL.obj = hclS[hclpt-0].obj
		}
	case 25:
		//line parse.y:180
		{
			hclVAL.obj = &Object{
				Type:  ValueTypeString,
				Value: hclS[hclpt-0].str,
			}
		}
	case 26:
		//line parse.y:189
		{
			hclVAL.obj = &Object{
				Type:  ValueTypeInt,
				Value: hclS[hclpt-0].num,
			}
		}
	case 27:
		//line parse.y:196
		{
			hclVAL.obj = &Object{
				Type:  ValueTypeFloat,
				Value: hclS[hclpt-0].f,
			}
		}
	case 28:
		//line parse.y:203
		{
			fs := fmt.Sprintf("%d%s", hclS[hclpt-1].num, hclS[hclpt-0].str)
			f, err := strconv.ParseFloat(fs, 64)
			if err != nil {
				panic(err)
			}

			hclVAL.obj = &Object{
				Type:  ValueTypeFloat,
				Value: f,
			}
		}
	case 29:
		//line parse.y:216
		{
			fs := fmt.Sprintf("%f%s", hclS[hclpt-1].f, hclS[hclpt-0].str)
			f, err := strconv.ParseFloat(fs, 64)
			if err != nil {
				panic(err)
			}

			hclVAL.obj = &Object{
				Type:  ValueTypeFloat,
				Value: f,
			}
		}
	case 30:
		//line parse.y:231
		{
			hclVAL.num = hclS[hclpt-0].num * -1
		}
	case 31:
		//line parse.y:235
		{
			hclVAL.num = hclS[hclpt-0].num
		}
	case 32:
		//line parse.y:241
		{
			hclVAL.f = hclS[hclpt-0].f * -1
		}
	case 33:
		//line parse.y:245
		{
			hclVAL.f = hclS[hclpt-0].f
		}
	case 34:
		//line parse.y:251
		{
			hclVAL.str = "e" + strconv.FormatInt(int64(hclS[hclpt-0].num), 10)
		}
	case 35:
		//line parse.y:255
		{
			hclVAL.str = "e-" + strconv.FormatInt(int64(hclS[hclpt-0].num), 10)
		}
	}
	goto hclstack /* stack new state and value */
}
