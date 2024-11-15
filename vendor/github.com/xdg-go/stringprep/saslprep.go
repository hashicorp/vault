package stringprep

var mapNonASCIISpaceToASCIISpace = Mapping{
	0x00A0: []rune{0x0020},
	0x1680: []rune{0x0020},
	0x2000: []rune{0x0020},
	0x2001: []rune{0x0020},
	0x2002: []rune{0x0020},
	0x2003: []rune{0x0020},
	0x2004: []rune{0x0020},
	0x2005: []rune{0x0020},
	0x2006: []rune{0x0020},
	0x2007: []rune{0x0020},
	0x2008: []rune{0x0020},
	0x2009: []rune{0x0020},
	0x200A: []rune{0x0020},
	0x200B: []rune{0x0020},
	0x202F: []rune{0x0020},
	0x205F: []rune{0x0020},
	0x3000: []rune{0x0020},
}

// SASLprep is a pre-defined stringprep profile for user names and passwords
// as described in RFC-4013.
//
// Because the stringprep distinction between query and stored strings was
// intended for compatibility across profile versions, but SASLprep was never
// updated and is now deprecated, this profile only operates in stored
// strings mode, prohibiting unassigned code points.
var SASLprep Profile = saslprep

var saslprep = Profile{
	Mappings: []Mapping{
		TableB1,
		mapNonASCIISpaceToASCIISpace,
	},
	Normalize: true,
	Prohibits: []Set{
		TableA1,
		TableC1_2,
		TableC2_1,
		TableC2_2,
		TableC3,
		TableC4,
		TableC5,
		TableC6,
		TableC7,
		TableC8,
		TableC9,
	},
	CheckBiDi: true,
}
